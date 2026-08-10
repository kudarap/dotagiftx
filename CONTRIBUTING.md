# Contributing to DotagiftX

Thank you for your interest in contributing! This guide covers how to set up a
development environment, run checks, and open a good pull request.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting started

### Prerequisites

- Go 1.26 https://go.dev/dl
- Docker 29.x https://docs.docker.com/get-docker
- Node 24.x https://nodejs.org/en/download
- Bun 1.x https://bun.com/get

### Local setup

1. Clone the repository and create your env config.

   ```shell
   git clone https://github.com/kudarap/dotagiftx.git
   cd dotagiftx
   cp .env.sample .env
   ```

2. In one terminal, start the local databases.

   ```shell
   make local
   ```

3. In another terminal, run the backend server.

   ```shell
   make run
   ```

4. In a third terminal, run the web client.

   ```shell
   cd web
   cp .env.sample .env
   bun install
   bun dev
   ```

The API will be served at `http://localhost:8000` and the web app at
`http://localhost:3000`.

## Making changes

- Fork the repository and create a branch from `dev` with a descriptive name,
  e.g. `feature/awesome-marketplace-search` or `fix/duplicate-order`.
- Keep changes focused on a single concern. Small, reviewable PRs are preferred.
- Follow the existing code style. For Go, run `gofmt`; for the web client,
  run `prettier`.

## Running checks

The backend test suite runs lint, vulnerability checks, and unit tests:

```shell
make test
```

To build all binaries:

```shell
make build
make build-worker
```

The web client lints and builds with:

```shell
cd web
bun install
bun run lint
bun run build
```

## Opening a pull request

- Push your branch and open a PR against `dev`.
- Fill out the pull request template and describe what the change does, why,
  and how it was tested.
- Keep commits tidy; the maintainers will squash when merging.
- CI (GitHub Actions) runs Go lint, vulncheck, tests, and web lint on every PR.
  Make sure it is green before requesting review.

## Reporting bugs

Please open an issue using the [bug report template](.github/ISSUE_TEMPLATE/bug_report.yml).
If the issue involves a security vulnerability, follow the instructions in
[SECURITY.md](SECURITY.md) instead of opening a public issue.
