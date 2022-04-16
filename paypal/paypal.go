package paypal

import (
	"context"
	"strings"

	"github.com/plutov/paypal/v4"
)

// Config represents paypal config.
type Config struct {
	ClientID string
	Secret   string
	Live     bool
}

// Client represents paypal client.
type Client struct {
	pc *paypal.Client
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
	if c.pc.Token == nil {
		_, err = c.pc.GetAccessToken(context.Background())
		if err != nil {
			return
		}
	}

	ctx := context.Background()
	sub, err := c.pc.GetSubscriptionDetails(ctx, id)
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
	p, err := c.pc.GetSubscriptionPlan(ctx, planID)
	if err != nil {
		return
	}
	return p.Name, nil
}
