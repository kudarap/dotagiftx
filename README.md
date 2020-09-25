# DotagiftX

Market place for giftable Dota 2 items

### Tech Stack

- Docker 19
- RethinkDB 2.4
- Go 1.15

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

### API endpoints

- public

  - [x] `GET /auth/steam` -- user login/register
  - [x] `GET /auth/renew` -- renews access token
  - [x] `GET /auth/revoke` -- revokes access token
  - [x] `GET /items` -- item search
  - [x] `GET /items/{item-id}` -- item details
  - [x] `GET /catalog` -- indexed market search
  - [x] `GET /catalog/{item-id}` -- indexed market search
  - [x] `GET /markets` -- market search
  - [x] `GET /markets/{market-id}` -- item market details
  - [x] `GET /users/{steam-id}` -- user details
  - [x] `GET /stats/top-origins` -- top origins stats
  - [x] `GET /stats/top-heroes` -- top heroes stats
  - [x] `GET /` -- api info

- private
  - [x] `GET /my/profile` -- user profile details
  - [x] `GET /my/markets` -- user market list
  - [x] `GET /my/markets/{market-id}` -- user market listing details
  - [x] `POST /my/markets` -- create user market
  - [x] `PATCH /my/markets` -- update user market
  - [x] `POST /items` -- create item
