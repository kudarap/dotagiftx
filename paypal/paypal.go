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

func (c *Client) Subscription(ctx context.Context, id string) (plan, steamID string, err error) {
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
	return plan, strings.TrimPrefix(sub.CustomID, customIDPrefix), nil
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

type subscriptionPayload struct {
	Resource struct {
		ID          string `json:"id"`
		CustomID    string `json:"custom_id"`
		Status      string `json:"status"`
		BillingInfo struct {
			OutstandingBalance struct {
				CurrencyCode string `json:"currency_code"`
				Value        string `json:"value"`
			} `json:"outstanding_balance"`
			CycleExecutions []struct {
				TenureType                  string `json:"tenure_type"`
				Sequence                    int    `json:"sequence"`
				CyclesCompleted             int    `json:"cycles_completed"`
				CyclesRemaining             int    `json:"cycles_remaining"`
				CurrentPricingSchemeVersion int    `json:"current_pricing_scheme_version"`
				TotalCycles                 int    `json:"total_cycles"`
			} `json:"cycle_executions"`
			LastPayment struct {
				Amount struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"amount"`
				Time time.Time `json:"time"`
			} `json:"last_payment"`
			FailedPaymentsCount int `json:"failed_payments_count"`
		} `json:"billing_info"`
	} `json:"resource"`
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
	defer req.Body.Close()

	var sub subscriptionPayload
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
