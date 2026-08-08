package paypal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
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
	pc *BaseClient

	webhookID string
}

func New(conf Config) (*Client, error) {
	c, err := NewBaseClient(conf.ClientID, conf.Secret, conf.Live)
	if err != nil {
		return nil, err
	}

	return &Client{c, conf.WebhookID}, nil
}

func (c *Client) Subscription(
	ctx context.Context,
	id string,
) (plan, steamID, subscriptionID string, lastPayment time.Time, err error) {
	if c.pc.Token == nil {
		_, err = c.pc.GetAccessToken(context.Background())
		if err != nil {
			return
		}
	}

	sub, err := c.pc.GetSubscriptionDetails(ctx, id)
	if err != nil {
		return
	}
	plan, err = c.planName(ctx, sub.PlanID)
	if err != nil {
		return
	}
	return plan, strings.TrimPrefix(sub.CustomID, customIDPrefix), sub.ID, sub.BillingInfo.LastPayment.Time, nil
}

// CreateSubscription creates a new subscription for the given plan ID
func (c *Client) CreateSubscription(ctx context.Context, planID, customID string) (id string, err error) {
	if planID == "" || customID == "" {
		return "", errors.New("missing plan ID or custom ID")
	}

	if c.pc.Token == nil {
		_, err = c.pc.GetAccessToken(context.Background())
		if err != nil {
			return
		}
	}

	sub, err := c.pc.CreateSubscription(ctx, planID, customIDPrefix+customID)
	if err != nil {
		return
	}

	return sub.ID, nil
}

type SubscriptionEventPayload struct {
	Resource Subscription `json:"resource"`
}

func (c *Client) IsCancelled(
	ctx context.Context,
	req *http.Request,
) (steamID string, cancelled bool, lastPayment time.Time, err error) {
	res, err := c.pc.VerifyWebhookSignature(ctx, req, c.webhookID)
	if err != nil {
		err = fmt.Errorf("invalid signature: %s", err)
		return
	}
	if strings.ToUpper(res.VerificationStatus) != "SUCCESS" {
		err = fmt.Errorf("verification failed: %s", res.VerificationStatus)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	defer func() { _ = req.Body.Close() }()

	var sub SubscriptionEventPayload
	if err = json.Unmarshal(body, &sub); err != nil {
		return
	}

	cancelled = slices.Contains([]string{"CANCELLED", "SUSPENDED"}, sub.Resource.Status)
	lastPayment = sub.Resource.BillingInfo.LastPayment.Time
	return strings.TrimPrefix(sub.Resource.CustomID, customIDPrefix), cancelled, lastPayment, nil
}

func (c *Client) planName(ctx context.Context, planID string) (name string, err error) {
	p, err := c.pc.GetSubscriptionPlan(ctx, planID)
	if err != nil {
		return
	}
	return p.Name, nil
}
