# Sentinel Tick

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/filecoin-project/sentinel-tick)
[![docker build status](https://img.shields.io/docker/cloud/build/filecoin/sentinel-tick?style=flat-square)](https://hub.docker.com/repository/docker/filecoin/sentinel-tick)
![Go](https://github.com/filecoin-project/sentinel-tick/workflows/Go/badge.svg)

A component of [**Sentinel**](https://github.com/filecoin-project/sentinel), a collection of services which monitor the health and function of the Filecoin network. 

**Sentinel-Tick** collects price information for exchanges where Filecoin is traded and writes it to a PostgreSQL (normally the same one as other Sentinel components).

## Usage

`sentinel-tick --help`

## Code of Conduct

Sentinel Tick follows the [Filecoin Project Code of Conduct](https://github.com/filecoin-project/community/blob/master/CODE_OF_CONDUCT.md). Before contributing, please acquaint yourself with our social courtesies and expectations.


## Contributing

Welcoming [new issues](https://github.com/filecoin-project/sentinel-tick/issues/new) and [pull requests](https://github.com/filecoin-project/sentinel-tick/pulls).


## License

Sentinel Tick is dual-licensed under Apache 2.0 and MIT terms:

- Apache License, Version 2.0, ([LICENSE-APACHE](https://github.com/filecoin-project/sentinel-tick/blob/master/LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
- MIT license ([LICENSE-MIT](https://github.com/filecoin-project/sentinel-tick/blob/master/LICENSE-MIT) or http://opensource.org/licenses/MIT)
