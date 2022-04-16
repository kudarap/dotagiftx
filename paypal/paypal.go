package paypal

import (
	"context"
	"strings"

	"github.com/plutov/paypal/v4"
)

type Client struct {
	paypalClient *paypal.Client
}

func New(clientID, secret string, live bool) (*Client, error) {
	base := paypal.APIBaseSandBox
	if live {
		base = paypal.APIBaseLive
	}
	c, err := paypal.NewClient(clientID, secret, base)
	if err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

const customIDPrefix = "STEAMID-"

func (c *Client) Subscription(id string) (plan, steamID string, err error) {
	if c.paypalClient.Token == nil {
		_, err = c.paypalClient.GetAccessToken(context.Background())
		if err != nil {
			return
		}
	}

	ctx := context.Background()
	sub, err := c.paypalClient.GetSubscriptionDetails(ctx, id)
	if err != nil {
		return
	}
	plan, err = c.planName(ctx, sub.PlanID)
	if err != nil {
		return
	}
	return plan, strings.TrimPrefix(sub.CustomID, customIDPrefix), nil
}

func (c *Client) planName(ctx context.Context, planID string) (name string, err error) {
	p, err := c.paypalClient.GetSubscriptionPlan(ctx, planID)
	if err != nil {
		return
	}
	return strings.ToUpper(p.Name), nil
}
