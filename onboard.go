package paypal

import (
	"context"
	"fmt"
)


type PartnerReferralResponse struct {
	Links []struct {
		Href        string `json:"href"`
		Rel         string `json:"rel"`
		Method      string `json:"method"`
		Description string `json:"description"`
	} `json:"links"`
}

type PartnerReferralMerchantIdResponse struct {
	MerchantId string `json:"merchant_id"`
	TrackingId string `json:"tracking_id"`
	Links      []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}

func (c *Client) GetPartnerReferrals(ctx context.Context,p ReferralRequest) (*PartnerReferralResponse,error) {

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/customer/partner-referrals"), p)
	response := &PartnerReferralResponse{}
	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

func (c *Client) GetMerchantIdByTrackingId(ctx context.Context, partnerId,trackingId string)(*PartnerReferralMerchantIdResponse, error)  {

	uri := fmt.Sprintf("/v1/customer/partners/%s/merchant-integrations?tracking_id=%s",partnerId,trackingId )
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, uri), nil)
	response := &PartnerReferralMerchantIdResponse{}
	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

func (c *Client) GetSellerAccountStatus(ctx context.Context,partnerId,merchantId string) (*PartnerReferralMerchantIdResponse,error) {

	uri := fmt.Sprintf("/v1/customer/partners/%s/merchant-integrations/%s",partnerId,merchantId )
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, uri), nil)
	response := &PartnerReferralMerchantIdResponse{}
	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}