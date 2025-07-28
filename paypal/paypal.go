package paypal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

const customIDPrefix = "STEAMID-"

// Config represents paypal config.
type Config struct {
	Live      bool
	ClientID  string
	Secret    string
	WebhookID string
}

// Client represents paypal client.
type Client struct {
	pc *paypalClient

	webhookID string
}

func New(conf Config) (*Client, error) {
	base := apiBaseSandbox
	if conf.Live {
		base = apiBaseLive
	}
	c, err := NewClient(conf.ClientID, conf.Secret, base)
	if err != nil {
		return nil, err
	}

	return &Client{c, conf.WebhookID}, nil
}

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

type subscriptionPayload struct {
	Resource struct {
		ID       string `json:"id"`
		CustomID string `json:"custom_id"`
		Status   string `json:"status"`
	} `json:"resource"`
}

func (c *Client) IsCancelled(ctx context.Context, req *http.Request) (steamID string, cancelled bool, err error) {
	res, err := c.pc.VerifyWebhookSignature(ctx, req, c.webhookID)
	if err != nil {
		return "", false, fmt.Errorf("invalid signature: %s", err)
	}
	if strings.ToUpper(res.VerificationStatus) != "SUCCESS" {
		return "", false, fmt.Errorf("verification failed: %s", res.VerificationStatus)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", false, err
	}
	defer req.Body.Close()

	var sub subscriptionPayload
	if err = json.Unmarshal(body, &sub); err != nil {
		return "", false, err
	}

	cancelled = slices.Contains([]string{"CANCELLED", "SUSPENDED"}, sub.Resource.Status)
	return sub.Resource.CustomID, cancelled, nil
}

func (c *Client) planName(ctx context.Context, planID string) (name string, err error) {
	p, err := c.pc.GetSubscriptionPlan(ctx, planID)
	if err != nil {
		return
	}
	return p.Name, nil
}
