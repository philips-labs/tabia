# Tabia

[![Go CI](https://github.com/philips-labs/tabia/workflows/Go%20CI/badge.svg)](https://github.com/philips-labs/tabia/actions)
[![codecov](https://codecov.io/gh/philips-labs/tabia/branch/develop/graph/badge.svg?token=K2R9WOXNBm)](https://codecov.io/gh/philips-labs/tabia)

Tabia means characteristic in Swahili. Tabia is giving us insights on the characteristics of our code bases.

## Setup

Copy `.env.example` to `.env` and fill out the bitbucket token. This environment variable is read by the CLI and tests. Also vscode will read the variable when running tests or starting debugger.

```bash
cp .env.example .env
source .env
env | grep TABIA
```

## Build

To build the CLI you can make use of the `build` target using `make`.

```bash
make build
```

## Test

To run tests you can make use of the `test` target using `make`.

```bash
make test
```

## Run

### Bitbucket

To interact with Bitbucket `tabia` makes use of the [Bitbucket 1.0 Rest API](https://docs.atlassian.com/bitbucket-server/rest/7.3.0/bitbucket-rest.html).

```bash
bin/tabia bitbucket --help
bin/tabia bitbucket projects --help
bin/tabia bitbucket repositories --help
```
