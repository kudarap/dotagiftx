# Dota2 Giftables

Market place for giftable Dota 2 items

### Tech Stack

- Docker 19
- RethinkDB 2.4
- Go 1.14

### Architecture

- Standard Package Layout
- Dependency Injections
- Containerized

### Entities

- auth
- user
- item
- sell

### API endpoints

- public
    - [x] `GET /auth/steam` -- user login/register
    - [x] `GET /auth/renew` -- renews access token
    - [x] `GET /auth/revoke` -- revokes access token
    - [x] `GET /items` -- item search
    - [x] `GET /items/{item-id}` -- item details
    - [x] `GET /market` -- sell search
    - [x] `GET /market/{market-id}` -- item sell details
    - [x] `GET /users/{steam-id}` -- user details

- private
    - [x] `GET /my/profile` -- user profile details
    - [x] `GET /my/sells` -- user sell list
    - [x] `GET /my/sells/{sell-id}` -- user sell listing details
    - [x] `POST /my/sells` -- create user sell
    - [x] `PATCH /my/sells` -- update user sell
    - [x] `POST /items` -- create item
