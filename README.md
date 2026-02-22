# DotagiftX

Marketplace for giftable Dota 2 items

### API reference

- [Postman Collection](/postman.json)

### Requirements

- Go 1.26
- Node 24.x
- Yarn 1.22
- Docker 29.x

### Local Setup

- Create a new env config and change accordingly. Change `DG_PAYPAL_*` values with your own sandbox account credentials.

```shell
$ cp .env.sample .env
```

- Open a new terminal to setup databases.

```shell
$ make local
```

- Open a new terminal to run backend server.

```shell
$ make run
```

- Open a new terminal to setup and run web client.

```shell
$ cd web
$ cp .env.sample .env
$ yarn
$ yarn dev
```
