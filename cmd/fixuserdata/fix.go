package main

import (
	"flag"
	"log/slog"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/rethink"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func fixUserData(
	badSteamIDs []string,
	slogger *slog.Logger,
	userStg dotagiftx.UserStorage,
	marketSvc dotagiftx.MarketService,
	rethinkClient *rethink.Client,
) {
	var finalRun bool
	flag.BoolVar(&finalRun, "final", false, "final run")
	flag.Parse()

	slogger.Info("starting")
	allStart := time.Now()

	var completed, errorSet []string
	for _, steamID := range badSteamIDs {
		start := time.Now()
		user, err := userStg.Get(steamID)
		if err != nil {
			slogger.Error("get user. skipping", "sid", steamID, "err", err)
			continue
		}

		var joinedDate time.Time
		// trash old auth data
		authQuery := r.Table("auth").GetAllByIndex("username", steamID).
			OrderBy("created_at")
		err = authQuery.Limit(1).Field("created_at").Default(time.Time{}).
			ReadOne(&joinedDate, rethinkClient.Session())
		if err != nil {
			slogger.Warn("get join date", "sid", steamID, "err", err)
		}
		if finalRun && !joinedDate.IsZero() {
			err = r.Table("auth_trash").Insert(authQuery).Exec(rethinkClient.Session())
			if err != nil {
				slogger.Error("mark old auth data for trashing. skipping", "sid", steamID, "err", err)
				continue
			}
			err = authQuery.Delete().Exec(rethinkClient.Session())
			if err != nil {
				slogger.Error("mark old auth data for deletion. skipping", "sid", steamID, "err", err)
				continue
			}
		}

		// clean user data
		userQuery := r.Table("user").GetAllByIndex("steam_id", steamID).
			OrderBy("created_at").Slice(1)
		var userIDs []string
		err = userQuery.Field("id").ReadAll(&userIDs, rethinkClient.Session())
		if err != nil {
			slogger.Error("get user ids. skipping", "sid", steamID, "err", err)
			continue
		}

		if finalRun {
			err = r.Table("user_trash").Insert(userQuery).Exec(rethinkClient.Session())
			if err != nil {
				slogger.Error("mark old user data for trashing. skipping", "sid", steamID, "err", err)
				continue
			}
			err = userQuery.Delete().Exec(rethinkClient.Session())
			if err != nil {
				slogger.Error("mark old user data for deletion. skipping", "sid", steamID, "err", err)
				continue
			}
			// assign oldest join date to user
			if !joinedDate.IsZero() && joinedDate.Before(*user.CreatedAt) {
				user.CreatedAt = &joinedDate
				err = userStg.Update(user)
				if err != nil {
					slogger.Error("update user. skipping", "sid", steamID, "err", err)
					continue
				}
				slogger.Info("updated user join date", "sid", steamID, "joined_at", joinedDate)
			}

			// merge market data
			for _, uid := range userIDs {
				slogger.Info("merging market data...",
					"old_user_id", uid,
					"user_id", user.ID,
				)
				err = r.Table("market").GetAllByIndex("user_id", uid).Update(map[string]any{
					"old_user_id":      uid,
					"user_id":          user.ID,
					"manual_fix_notes": "fix account merge",
				}).Exec(rethinkClient.Session())
				if err != nil {
					slogger.Error("merge market data. skipping", "sid", steamID, "err", err)
					continue
				}
			}
		}

		if finalRun {
			// generate user stats
			if err = marketSvc.UpdateUserRankScore(user.ID); err != nil {
				slogger.Error("generate user stats. skipping", "sid", steamID, "err", err)
				continue
			}
		}

		slogger.Info("done", "sid", steamID, "elapsed", time.Since(start))
		completed = append(completed, steamID)
	}

	slogger.Info("all DONE!",
		"completed", len(completed),
		"total", len(badSteamIDs),
		"errors", len(errorSet),
		"errors_set", errorSet,
		"elapsed", time.Since(allStart),
	)
}
