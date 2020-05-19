# Tabia

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

```bash
bin/tabia bitbucket --help
bin/tabia bitbucket projects --help
bin/tabia bitbucket repositories --help
```
