package paypal

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GrantNewAccessTokenFromAuthCode - Use this call to grant a new access token, using the previously obtained authorization code.
// Endpoint: POST /v1/identity/openidconnect/tokenservice
func (c *Client) GrantNewAccessTokenFromAuthCode(ctx context.Context, code, redirectURI string) (*TokenResponse, error) {
	token := &TokenResponse{}

	q := url.Values{}
	q.Set("grant_type", "authorization_code")
	q.Set("code", code)
	q.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/identity/openidconnect/tokenservice"), strings.NewReader(q.Encode()))
	if err != nil {
		return token, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err = c.SendWithBasicAuth(req, token); err != nil {
		return token, err
	}

	return token, nil
}

// GrantNewAccessTokenFromRefreshToken - Use this call to grant a new access token, using a refresh token.
// Endpoint: POST /v1/identity/openidconnect/tokenservice
func (c *Client) GrantNewAccessTokenFromRefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	type request struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}

	token := &TokenResponse{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/identity/openidconnect/tokenservice"), request{GrantType: "refresh_token", RefreshToken: refreshToken})
	if err != nil {
		return token, err
	}

	if err = c.SendWithAuth(req, token); err != nil {
		return token, err
	}

	return token, nil
}

// GetUserInfo - Use this call to retrieve user profile attributes.
// Endpoint: GET /v1/identity/openidconnect/userinfo/?schema=<Schema>
// Pass the schema that is used to return as per openidconnect protocol. The only supported schema value is openid.
func (c *Client) GetUserInfo(ctx context.Context, schema string) (*UserInfo, error) {
	u := &UserInfo{}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v1/identity/openidconnect/userinfo/?schema=", schema), nil)
	if err != nil {
		return u, err
	}

	if err = c.SendWithAuth(req, u); err != nil {
		return u, err
	}

	return u, nil
}

func (c *Client) GetSellerAccessToken(ctx context.Context,shareId,authCode, sellerNonce string) (*SellerAccessToken,error) {

	token := &SellerAccessToken{}

	buf := bytes.NewBuffer([]byte(fmt.Sprintf("grant_type=authorization_code&code=%s&code_verifier=%s", authCode,sellerNonce)))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return token, err
	}

	basic := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:",shareId)))
	req.Header.Set("Authorization","Basic " + basic)
	err = c.Send(req, token);
	if  err != nil {
		return token, err
	}

	return token, nil
}

func (c *Client) GetSellerCredentials(ctx context.Context,sellerAccessToken,partnerMerchantId string) (*SellerCredentials,error) {

	credentials := &SellerCredentials{}
	uri := fmt.Sprintf("/v1/customer/partners/%s/merchant-integrations/credentials/", partnerMerchantId)
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, uri), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s",sellerAccessToken))

	if err = c.SendWithAuth(req, credentials); err != nil {
		return credentials, err
	}

	return credentials, nil
}
