package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kudarap/dotagiftx/config"
	"github.com/kudarap/dotagiftx/dotagiftx"
	"github.com/kudarap/dotagiftx/paypal"
	"github.com/kudarap/dotagiftx/rethink"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	configPrefix = "DG"

	cachePaypalSubscriptionFile = ".localdata/paypal-susbcriptions.json"
)

func main() {
	var final, ignoreCache bool
	flag.BoolVar(&final, "final", false, "final mode")
	flag.BoolVar(&ignoreCache, "ignore-cache", false, "use cache")
	flag.Parse()

	config.EnvPrefix = configPrefix
	var conf config.Config
	errCheck(config.Load(&conf))

	paypalClient, err := paypal.NewBaseClient(conf.Paypal.ClientID, conf.Paypal.Secret, conf.Paypal.Live)
	errCheck(err)

	rethinkClient, err := rethink.New(conf.Rethink)
	errCheck(err)
	session := rethinkClient.Session()

	ctx := context.Background()
	paypalSubsIndex := paypalSubscribers(ctx, paypalClient, "ACTIVE", !ignoreCache)

	for _, dgxSubs := range []dotagiftx.UserSubscription{
		dotagiftx.UserSubscriptionPartner,
		dotagiftx.UserSubscriptionTrader,
		dotagiftx.UserSubscriptionSupporter,
	} {
		for _, user := range subscribers(session, dgxSubs) {
			subs, ok := paypalSubsIndex[fmt.Sprintf("STEAMID-%s", user.SteamID)]
			if !ok {
				subs = paypal.Subscription{}
			}

			printPretty(user, subs)
			fixData(final, rethinkClient.Session(), user, subs)
		}
	}

}

func fixData(final bool, db *r.Session, user dotagiftx.User, subs paypal.Subscription) {
	if user.SubscriptionType == "manual" {
		return
	}

	// set paypal subs type on empty
	// set paypal sub id
	payload := map[string]any{
		"paypal": map[string]any{
			"subscription_id":         subs.ID,
			"subscription_last_payed": subs.BillingInfo.LastPayment.Time,
		},
	}
	if user.SubscriptionType != "manual" {
		payload["subscription_type"] = "paypal"
	}

	// set subs expiration on non-active status
	if subs.Status != "ACTIVE" {
		payload["subscription_ends_at"] = subs.BillingInfo.LastPayment.Time.Add(time.Hour * 24 * 30)
	}

	payload["notes_aug_2026"] = "bulk fix subscription"

	if !final {
		fmt.Println("[dry-run] ", user.SteamID, user.Name)
		jsonData, err := json.MarshalIndent(payload, "", "  ")
		errCheck(err)
		fmt.Println(string(jsonData))
		fmt.Println("----")
		return
	}

	err := r.Table("user").Get(user.ID).Update(payload).Exec(db)
	errCheck(err)
}

func paypalSubscribers(
	ctx context.Context,
	paypalClient *paypal.BaseClient,
	status string,
	cache bool,
) map[string]paypal.Subscription {
	var err error
	var paypalSubs []paypal.Subscription

	// check paypal subs cache
	if cache {
		file, err := os.Open(cachePaypalSubscriptionFile)
		if !os.IsNotExist(err) {
			body, err := io.ReadAll(file)
			errCheck(err)

			err = json.Unmarshal(body, &paypalSubs)
			errCheck(err)
		}
		defer func() { _ = file.Close() }()
	}

	if len(paypalSubs) == 0 {
		paypalSubs, err = paypalClient.ListSubscriptions(ctx, status)
		errCheck(err)
	}

	paypalSubsIndex := map[string]paypal.Subscription{}
	for _, subs := range paypalSubs {
		curr, exists := paypalSubsIndex[subs.ID]
		if exists && curr.BillingInfo.LastPayment.Time.Before(subs.BillingInfo.LastPayment.Time) {
			continue
		}

		paypalSubsIndex[subs.CustomID] = subs
	}

	// cache paypal subs response
	if len(paypalSubs) != 0 {
		b, err := json.MarshalIndent(paypalSubs, "", "  ")
		errCheck(err)

		err = os.WriteFile(cachePaypalSubscriptionFile, b, 0600)
		errCheck(err)
	}

	return paypalSubsIndex
}

func subscribers(db *r.Session, subsID dotagiftx.UserSubscription) []dotagiftx.User {
	var res []dotagiftx.User
	q := r.Table("user").
		GetAllByIndex("subscription", subsID).
		OrderBy(r.Desc("updated_at"))
	err := q.ReadAll(&res, db)
	errCheck(err)

	return res
}

func printPretty(user dotagiftx.User, subs paypal.Subscription) {
	data := struct {
		Name         string
		UserID       string
		SteamID      string
		Subscription string
		SubsType     string
		SubsStarts   *time.Time
		SubsEnds     *time.Time
		Boons        []string `json:"-"`

		PaypalSubID     string
		PaypalSubsPlan  string
		PaypalSubStatus string
		PaypalLastPayed time.Time
	}{
		Name:         user.Name,
		UserID:       user.ID,
		SteamID:      user.SteamID,
		Subscription: user.Subscription.String(),
		SubsType:     user.SubscriptionType,
		SubsStarts:   user.SubscribedAt,
		SubsEnds:     user.SubscriptionEndsAt,
		Boons:        user.Boons,

		PaypalSubID:     subs.ID,
		PaypalSubsPlan:  subs.PlanID,
		PaypalSubStatus: subs.Status,
		PaypalLastPayed: subs.BillingInfo.LastPayment.Time,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	errCheck(err)

	fmt.Println(string(jsonData))
}

func errCheck(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}
