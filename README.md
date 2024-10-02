# DotagiftX

Market place for giftable Dota 2 items

### Tech Stack

- Go 1.23
- RethinkDB 2.4
- Redis 6.2
- Docker 20

### Architecture

- Standard Package Layout
- Dependency Injections
- Containerized

### Entities

- auth
- user
- item
- market
- catalog(market index)
- report

### API endpoints

- public

  - [x] `GET /auth/steam` -- user login/register
  - [x] `GET /auth/renew` -- renews access token
  - [x] `GET /auth/revoke` -- revokes access token
  - [x] `GET /items` -- item search
  - [x] `GET /items/{item-id}` -- item details
  - [x] `GET /catalogs` -- indexed market search
  - [x] `GET /catalogs/{item-id}` -- indexed market search
  - [x] `GET /markets` -- market search
  - [x] `GET /markets/{market-id}` -- item market details
  - [x] `GET /users/{steam-id}` -- user details
  - [x] `GET /stats/top_origins` -- top origins stats
  - [x] `GET /stats/top_heroes` -- top heroes stats
  - [x] `GET /stats/market_summary` -- market status count
  - [x] `GET /catalogs_trend` -- trending items
  - [x] `GET /reports` -- report list
  - [x] `GET /reports/{report-id}` -- report details
  - [x] `GET /` -- api info

- private
  - [x] `GET /my/profile` -- user profile details
  - [x] `GET /my/markets` -- user market list
  - [x] `GET /my/markets/{market-id}` -- user market listing details
  - [x] `POST /my/markets` -- create user market
  - [x] `PATCH /my/markets` -- update user market
  - [x] `POST /items` -- create item
  - [x] `POST /items_import` -- yaml items import
  - [x] `POST /reports` -- create user report
