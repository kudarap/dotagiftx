// Package paypal is copied from https://github.com/plutov/paypal because I don't want to import the library
// and lazy to write a client, app just needs a few things.
package paypal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

const (
	// apiBaseSandbox points to the sandbox (for testing) version of the API
	apiBaseSandbox = "https://api-m.sandbox.paypal.com"

	// apiBaseLive points to the live version of the API
	apiBaseLive = "https://api-m.paypal.com"

	// requestNewTokenBeforeExpiresIn is used by SendWithAuth and try to get new Token when it's about to expire
	requestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second
)

type (
	// paypalClient represents a Paypal REST API paypalClient
	paypalClient struct {
		// sync.Mutex
		mu                   sync.Mutex
		Client               *http.Client
		ClientID             string
		Secret               string
		APIBase              string
		Log                  io.Writer // If user set log file name all requests will be logged there
		Token                *TokenResponse
		tokenExpiresAt       time.Time
		returnRepresentation bool
	}

	// TokenResponse is for API response for the /oauth2/token endpoint
	TokenResponse struct {
		RefreshToken string         `json:"refresh_token"`
		Token        string         `json:"access_token"`
		Type         string         `json:"token_type"`
		ExpiresIn    expirationTime `json:"expires_in"`
	}

	// VerifyWebhookResponse struct
	VerifyWebhookResponse struct {
		VerificationStatus string `json:"verification_status,omitempty"`
	}

	WebhookEventTypesResponse struct {
		EventTypes []WebhookEventType `json:"event_types"`
	}

	// ErrorResponseDetail struct
	ErrorResponseDetail struct {
		Field       string `json:"field"`
		Issue       string `json:"issue"`
		Name        string `json:"name"`
		Message     string `json:"message"`
		Description string `json:"description"`
		Links       []Link `json:"link"`
	}

	// ErrorResponse https://developer.paypal.com/docs/api/errors/
	ErrorResponse struct {
		Response        *http.Response        `json:"-"`
		Name            string                `json:"name"`
		DebugID         string                `json:"debug_id"`
		Message         string                `json:"message"`
		InformationLink string                `json:"information_link"`
		Details         []ErrorResponseDetail `json:"details"`
	}

	// Link struct
	Link struct {
		Href        string `json:"href"`
		Rel         string `json:"rel,omitempty"`
		Method      string `json:"method,omitempty"`
		Description string `json:"description,omitempty"`
		Enctype     string `json:"enctype,omitempty"`
	}

	// WebhookEventType struct
	WebhookEventType struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status,omitempty"`
	}

	expirationTime int64
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %s, %+v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*e = expirationTime(i)
	return nil
}

// ToDuration convert ExpirationTime to time.Duration
func (e *expirationTime) ToDuration() time.Duration {
	seconds := int64(*e)
	return time.Duration(seconds) * time.Second
}

// NewClient returns new paypalClient struct
// APIBase is a base API URL, for testing you can use paypal.apiBaseSandbox
func NewClient(clientID string, secret string, apiBase string) (*paypalClient, error) {
	if clientID == "" || secret == "" || apiBase == "" {
		return nil, errors.New("ClientID, Secret and apiBase are required to create a paypalClient")
	}

	return &paypalClient{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  apiBase,
	}, nil
}

// GetAccessToken returns struct of TokenResponse
// No need to call SetAccessToken to apply new access token for current paypalClient
// Endpoint: POST /v1/oauth2/token
func (c *paypalClient) GetAccessToken(ctx context.Context) (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	response := &TokenResponse{}
	err = c.SendWithBasicAuth(req, response)

	// Set Token fur current paypalClient
	if response.Token != "" {
		c.Token = response
		c.tokenExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	}

	return response, err
}

// SetHTTPClient sets *http.Client to current client
func (c *paypalClient) SetHTTPClient(client *http.Client) {
	c.Client = client
}

// SetAccessToken sets saved token to current client
func (c *paypalClient) SetAccessToken(token string) {
	c.Token = &TokenResponse{
		Token: token,
	}
	c.tokenExpiresAt = time.Time{}
}

// SetLog will set/change the output destination.
// If log file is set paypal will log all requests and responses to this Writer
func (c *paypalClient) SetLog(log io.Writer) {
	c.Log = log
}

// SetReturnRepresentation enables verbose response
// Verbose response: https://developer.paypal.com/docs/api/orders/v2/#orders-authorize-header-parameters
func (c *paypalClient) SetReturnRepresentation() {
	c.returnRepresentation = true
}

// Send makes a request to the API, the response body will be
// unmarshalled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *paypalClient) Send(req *http.Request, v interface{}) (retErr error) {
	var (
		err  error
		resp *http.Response
		data []byte
	)

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}
	if c.returnRepresentation {
		req.Header.Set("Prefer", "return=representation")
	}
	if c.Log != nil {
		if reqDump, err := httputil.DumpRequestOut(req, true); err == nil {
			logMsg := fmt.Sprintf("Request: %s\n", string(reqDump))
			if _, logErr := c.Log.Write([]byte(logMsg)); logErr != nil {
				return logErr
			}
		}
	}

	resp, err = c.Client.Do(req)
	if err != nil {
		return err
	}

	if c.Log != nil {
		if respDump, err := httputil.DumpResponse(resp, true); err == nil {
			logMsg := fmt.Sprintf("Response from %s: %s\n", req.URL.RequestURI(), string(respDump))
			if _, logErr := c.Log.Write([]byte(logMsg)); logErr != nil {
				return logErr
			}
		}
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil && retErr == nil {
			retErr = err
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &ErrorResponse{Response: resp}

		var bodyBuffer bytes.Buffer
		teeReader := io.TeeReader(resp.Body, &bodyBuffer)
		data, err = io.ReadAll(teeReader)
		resp.Body = io.NopCloser(&bodyBuffer)

		if err == nil && len(data) > 0 {
			err := json.Unmarshal(data, errResp)
			if err != nil {
				return err
			}
		}

		return errResp
	}
	if v == nil {
		return nil
	}

	if w, ok := v.(io.Writer); ok {
		_, err := io.Copy(w, resp.Body)
		return err
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// SendWithAuth makes a request to the API and apply OAuth2 header automatically.
// If the access token soon to be expired or already expired, it will try to get a new one before
// making the main request
// client.Token will be updated when changed
func (c *paypalClient) SendWithAuth(req *http.Request, v interface{}) error {
	// c.Lock()
	c.mu.Lock()
	// Note: Here we do not want to `defer c.Unlock()` because we need `c.Send(...)`
	// to happen outside of the locked section.

	if c.Token == nil || (!c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < requestNewTokenBeforeExpiresIn) {
		// c.Token will be updated in GetAccessToken call
		if _, err := c.GetAccessToken(req.Context()); err != nil {
			// c.Unlock()
			c.mu.Unlock()
			return err
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	// Unlock the client mutex before sending the request, this allows multiple requests
	// to be in progress at the same time.
	// c.Unlock()
	c.mu.Unlock()
	return c.Send(req, v)
}

// SendWithBasicAuth makes a request to the API using clientID:secret basic auth
func (c *paypalClient) SendWithBasicAuth(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.ClientID, c.Secret)

	return c.Send(req, v)
}

// NewRequest constructs a request
// Convert payload to a JSON
func (c *paypalClient) NewRequest(ctx context.Context, method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	return http.NewRequestWithContext(ctx, method, url, buf)
}

type (
	SubscriptionPlan struct {
		ID        string `json:"id,omitempty"`
		ProductId string `json:"product_id"`
		Name      string `json:"name"`
	}

	Subscription struct {
		SubscriptionDetailResp
	}

	// SubscriptionDetailResp struct
	SubscriptionDetailResp struct {
		PlanID   string `json:"plan_id"`
		CustomID string `json:"custom_id,omitempty"`
	}
)

// GetSubscriptionDetails shows details for a subscription, by ID.
// Endpoint: GET /v1/billing/subscriptions/
func (c *paypalClient) GetSubscriptionDetails(ctx context.Context, subscriptionID string) (*SubscriptionDetailResp, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/billing/subscriptions/%s", c.APIBase, subscriptionID), nil)
	response := &SubscriptionDetailResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// GetSubscriptionPlan get subscription plan
// Endpoint: GET /v1/billing/plans/:plan_id
func (c *paypalClient) GetSubscriptionPlan(ctx context.Context, planId string) (*SubscriptionPlan, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/plans/", planId), nil)
	response := &SubscriptionPlan{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// VerifyWebhookSignature - Use this to verify the signature of a webhook recieved from paypal.
// Endpoint: POST /v1/notifications/verify-webhook-signature
func (c *paypalClient) VerifyWebhookSignature(ctx context.Context, httpReq *http.Request, webhookID string) (*VerifyWebhookResponse, error) {
	type verifyWebhookSignatureRequest struct {
		AuthAlgo         string          `json:"auth_algo,omitempty"`
		CertURL          string          `json:"cert_url,omitempty"`
		TransmissionID   string          `json:"transmission_id,omitempty"`
		TransmissionSig  string          `json:"transmission_sig,omitempty"`
		TransmissionTime string          `json:"transmission_time,omitempty"`
		WebhookID        string          `json:"webhook_id,omitempty"`
		Event            json.RawMessage `json:"webhook_event,omitempty"`
	}

	// Read the content
	var bodyBytes []byte
	if httpReq.Body != nil {
		bodyBytes, _ = io.ReadAll(httpReq.Body)
	} else {
		return nil, errors.New("cannot verify webhook for HTTP Request with empty body")
	}
	// Restore the io.ReadCloser to its original state
	httpReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	verifyRequest := verifyWebhookSignatureRequest{
		AuthAlgo:         httpReq.Header.Get("PAYPAL-AUTH-ALGO"),
		CertURL:          httpReq.Header.Get("PAYPAL-CERT-URL"),
		TransmissionID:   httpReq.Header.Get("PAYPAL-TRANSMISSION-ID"),
		TransmissionSig:  httpReq.Header.Get("PAYPAL-TRANSMISSION-SIG"),
		TransmissionTime: httpReq.Header.Get("PAYPAL-TRANSMISSION-TIME"),
		WebhookID:        webhookID,
		Event:            json.RawMessage(bodyBytes),
	}

	response := &VerifyWebhookResponse{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/verify-webhook-signature"), verifyRequest)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetWebhookEventTypes - Lists all webhook event types.
// Endpoint: GET /v1/notifications/webhooks-event-types
func (c *paypalClient) GetWebhookEventTypes(ctx context.Context) (*WebhookEventTypesResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks-event-types"), nil)
	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()
	resp := &WebhookEventTypesResponse{}
	if err != nil {
		return nil, err
	}

	err = c.SendWithAuth(req, resp)
	return resp, err
}
