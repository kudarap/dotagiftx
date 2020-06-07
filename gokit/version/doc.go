package version

/*

Implementation:
	var tag, commit, built string

	func initVer(cfg Config) *version.Version {
		v := version.New(cfg.Prod, tag, commit, built)
		return v
	}

Usage:
	go build -ldflags=" \
			-X main.tag=`git describe --tag --abbrev=0` \
			-X main.commit=`git rev-parse HEAD` \
			-X main.built=`date -u +%s`" \
			./cmd/app.go

*/
