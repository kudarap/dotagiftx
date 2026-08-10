<div align="center">

# DotagiftX

**A marketplace for giftable Dota 2 items**

Buy and sell Dota 2 items as Steam gifts. DotagiftX pairs a Go API backend with
a Next.js web client, backed by RethinkDB, Redis, and optional ClickHouse
analytics, with automated gift-delivery verification powered by Steam.

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8)](https://go.dev)
[![Next.js](https://img.shields.io/badge/Next.js-16-black)](https://nextjs.org)

</div>

---

## Features

- **Marketplace** for giftable Dota 2 items with item catalogs, listing,
  selling, and market tracking
- **Steam integration** — login via Steam OpenID, player lookup, inventory
  scanning, and gift-delivery verification
- **Phantasm crawler** — background crawler that verifies inventories and item
  deliveries
- **PayPal checkout** — sandbox and live payment support with webhook handling
- **Gift wrapping & gifting** workflows with delivery status tracking
- **Image management** — upload, resize, and whitelisted image sources
- **Admin tooling** — hammer (user moderation) and reporting services
- **Optional analytics** — ClickHouse capture of track and market statistics
  with RethinkDB change feeds
- **Web client** — responsive Next.js UI with SSR, MUI, and Paypal buttons

## Architecture

The repository contains two Go binaries and a Next.js web client.

```
cmd/dxserver   HTTP API server (main entry point)
cmd/dxworker   Background worker for delivery/inventory verification jobs
web/           Next.js web client
```

Backend packages:

| Package       | Purpose                                        |
| ------------- | ---------------------------------------------- |
| `dotagiftx`   | Core domain: users, auth, items, markets, ...  |
| `http`        | HTTP router, handlers, middlewares, JWT auth   |
| `steam`       | Steam API client (auth, players, inventories)  |
| `phantasm`    | Inventory crawler service                      |
| `steaminvorg` | Steam inventory lookup / verification helpers  |
| `verify`      | Gift delivery verification logic              |
| `rethink`     | RethinkDB data stores and change feeds         |
| `redis`       | Redis/Valkey cache                             |
| `clickhouse`  | Optional stats capture and migrations          |
| `paypal`      | PayPal REST client and webhooks                |
| `file`        | File/image upload management                   |
| `logging`     | Structured logging and file output             |
| `discord`     | Discord webhook notifications                  |
| `tracing`     | OpenTracing spans (optional)                   |
| `worker`      | Job definitions for the background worker      |
| `config`      | Environment-based configuration                |
| `dota2`       | Dota 2 item/hero/treasure static data          |

### Data stores

- **RethinkDB** — primary datastore
- **Redis / Valkey** — caching and queues
- **ClickHouse** — optional analytics (disabled by default)

## API reference

The Postman collection documents the public API:

- [Postman Collection](postman.json)

## Requirements

- Go 1.26 https://go.dev/dl
- Docker 29.x https://docs.docker.com/get-docker
- Node 24.x https://nodejs.org/en/download
- Bun 1.x https://bun.com/get

## Credentials

- [Steam](https://steamcommunity.com/dev)
  - `DG_STEAM_KEY`
- [PayPal](https://developer.paypal.com)
  - `DG_PAYPAL_CLIENTID`
  - `DG_PAYPAL_SECRET`
  - `NEXT_PUBLIC_PAYPAL_CLIENT_ID` (web)

## Local Setup

### 1. Backend

Create your env config and adjust values as needed. Change `DG_PAYPAL_*` to your
own [PayPal sandbox](https://developer.paypal.com/dashboard/applications/sandbox)
credentials.

```shell
cp .env.sample .env
```

Open a new terminal to start the databases:

```shell
make local
```

Open another terminal to run the backend server:

```shell
make run
```

The API is served at `http://localhost:8000`. To run the background worker
instead, use `make run-worker`.

### 2. Web client

```shell
cd web
cp .env.sample .env
bun install
bun dev
```

The web app is served at `http://localhost:3000`. See [web/README.md](web/README.md).

## Configuration

All configuration is read from environment variables prefixed with `DG_` (the
web client uses `NEXT_PUBLIC_` prefixed variables). See [.env.sample](.env.sample)
for the full list of variables and defaults.

Key variables:

| Variable                        | Description                                  |
| ------------------------------- | -------------------------------------------- |
| `DG_SIGKEY`                     | JWT signing key (set a random secret)        |
| `DG_DIVINEKEY`                  | API request signing key (set a random secret)|
| `DG_PROD`                       | Production mode toggle                       |
| `DG_ADDR`                       | API listen address                           |
| `DG_RETHINK_ADDR`               | RethinkDB host                               |
| `DG_REDIS_ADDR`                 | Redis/Valkey host                            |
| `DG_STEAM_KEY`                  | Steam Web API key                            |
| `DG_PAYPAL_CLIENTID`/`_SECRET`  | PayPal REST credentials                      |
| `DG_PAYPAL_LIVE`                | Use PayPal live vs sandbox                   |
| `DG_STATS_CAPTURE_ENABLED`      | Enable ClickHouse analytics                  |
| `DG_PHANTASM_ADDRS`             | Phantasm crawler addresses (comma separated) |
| `DG_DISCORD_WEBHOOK_URL`        | Discord notifications webhook                |

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full contributing guide.

Common commands:

```shell
make test        # lint + vulncheck + Go tests
make lint        # golangci-lint
make vuln        # govulncheck
make build       # build dxserver binary
make build-worker# build dxworker binary
make docker-build# build the API Docker image
```

## Deployment

### Docker

Build and run the API image:

```shell
make docker-build
make docker-run
```

The image is built with `--platform=linux/amd64` support and can be pushed to
any registry. The `VERSION` file controls the image tag.

### Web client

The web client is a standard Next.js app that can be deployed to Vercel, a
Node/Bun server, or any static hosting provider. See [web/vercel.json](web/vercel.json).

## License

DotagiftX is released under the [MIT License](LICENSE).

Dota 2 is a trademark of Valve Corporation. This project is an independent
community project and is not affiliated with or endorsed by Valve Corporation.
