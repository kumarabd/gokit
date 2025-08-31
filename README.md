# Bucky Web Toolkit

<img align="right" width="159px" src="https://github.com/kumarabd/gokit/blob/master/assets/bucky.png">

[![Build Status](https://github.com/kumarabd/gokit/actions/workflows/go.yml/badge.svg?branch=v0.0.1)
[![codecov](https://codecov.io/gh/realnighthawk/bucky/branch/master/graph/badge.svg?token=8JUPZAA8ZW)](https://codecov.io/gh/realnighthawk/bucky)
[![Go Report Card](https://goreportcard.com/badge/github.com/kumarabd/gokit)](https://goreportcard.com/report/github.com/kumarabd/gokit)
<!--- [![GoDoc](https://pkg.go.dev/github.com/kumarabd/gokit?status.svg)](https://pkg.go.dev/github.com/kumarabd/gokit?tab=doc) 
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) 
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/gin-gonic/gin)](https://www.tickgit.com/browse?repo=github.com/gin-gonic/gin)
[![Sourcegraph](https://sourcegraph.com/github.com/gin-gonic/gin/-/badge.svg)](https://sourcegraph.com/github.com/gin-gonic/gin?badge)--->
[![Open Source Helpers](https://www.codetriage.com/realnighthawk/bucky/badges/users.svg)](https://www.codetriage.com/realnighthawk/bucky)
[![Release](https://img.shields.io/github/issues/realnighthawk/bucky)](https://github.com/kumarabd/gokit/releases)


Golang based library/toolkit for building Microservices written in Go (Golang). It is as durable as the [vibranium](https://marvel.fandom.com/wiki/Vibranium), which is also up to 40 times faster. If you are in need of full packed features and good productivity, bucky is your guy.


## Contents

- [Bucky Web Toolkit](#bucky-web-toolkit)
  - [Contents](#contents)
  - [Quick Start](#quick-start)
  - [CLI Tool](#cli-tool)
  - [Features](#features)
  - [Usage](#usage)
  - [Benchmarks](#benchmarks)
  - [Users](#users)

## Quick Start

### Using the CLI Tool

The easiest way to get started with GoKit is using our CLI tool:

```bash
# Install the CLI
go install github.com/kumarabd/gokit/cli@latest

# Create a new HTTP service
gokit new service --name user-service --template http

# Add monitoring to your service
gokit add monitoring --service user-service

# Add tracing to your service
gokit add tracing --service user-service
```

### Using the Library Directly

```go
package main

import (
    "github.com/kumarabd/gokit/config"
    "github.com/kumarabd/gokit/logger"
    "github.com/kumarabd/gokit/server"
)

type Config struct {
    Server server.HostPort `yaml:"server"`
    Log    logger.Options `yaml:"log"`
}

func main() {
    // Load configuration
    var cfg Config
    if err := config.New(&cfg); err != nil {
        log.Fatal(err)
    }

    // Initialize logger
    logger := logger.New(cfg.Log)
    defer logger.Close()

    logger.Info("Service started")
}
```

## CLI Tool

GoKit includes a powerful CLI tool for scaffolding and managing microservices:

### Installation

```bash
# Build from source
cd cli && go build -o gokit main.go

# Or install globally
go install github.com/kumarabd/gokit/cli@latest
```

### Commands

```bash
# Create new services
gokit new service --name user-service --template http
gokit new service --name payment-service --template grpc
gokit new service --name worker-service --template worker

# Add features to existing services
gokit add monitoring --service ./user-service
gokit add tracing --service ./payment-service
gokit add caching --service ./notification-service

# Show version
gokit version
```

### Supported Templates

- **HTTP Service**: RESTful API with health checks and graceful shutdown
- **gRPC Service**: gRPC server with protocol buffer support
- **Event Service**: Event-driven service for message processing
- **Worker Service**: Background job processing service

### Supported Features

- **Monitoring**: Prometheus metrics and monitoring
- **Tracing**: OpenTelemetry distributed tracing
- **Caching**: In-memory caching with TTL
- **Client**: HTTP client utilities
- **Middleware**: Common HTTP middleware (CORS, logging, recovery)

## Features

GoKit provides a comprehensive set of utilities for building microservices:

### Core Components

- **Configuration Management**: YAML, environment variables, and command-line flags
- **Structured Logging**: JSON logging with zerolog integration
- **Error Handling**: Standardized error types with severity levels
- **Caching**: Interface-based caching with in-memory implementation
- **HTTP Client**: Simple HTTP client for making requests
- **Server Abstractions**: HTTP and gRPC server interfaces

### Observability

- **Monitoring**: Prometheus metrics integration
- **Tracing**: OpenTelemetry support for distributed tracing
- **Health Checks**: Built-in health check endpoints
- **Graceful Shutdown**: Proper service shutdown handling

### Development Tools

- **CLI Scaffolding**: Generate new services with best practices
- **Feature Addition**: Add capabilities to existing services
- **Project Templates**: Standardized project structure
- **Documentation**: Comprehensive documentation and examples

## Usage

## Benchmarks

Bucky uses a custom version of [HttpRouter](https://github.com/julienschmidt/httprouter)

[See all benchmarks](/BENCHMARKS.md)

| Benchmark name                 |       (1) |             (2) |          (3) |             (4) |
| ------------------------------ | ---------:| ---------------:| ------------:| ---------------:|
| BenchmarkGin_GithubAll         | **43550** | **27364 ns/op** |   **0 B/op** | **0 allocs/op** |
| BenchmarkAce_GithubAll         |     40543 |     29670 ns/op |       0 B/op |     0 allocs/op |
| BenchmarkAero_GithubAll        |     57632 |     20648 ns/op |       0 B/op |     0 allocs/op |
| BenchmarkBear_GithubAll        |      9234 |    216179 ns/op |   86448 B/op |   943 allocs/op |
| BenchmarkBeego_GithubAll       |      7407 |    243496 ns/op |   71456 B/op |   609 allocs/op |
| BenchmarkBone_GithubAll        |       420 |   2922835 ns/op |  720160 B/op |  8620 allocs/op |
| BenchmarkChi_GithubAll         |      7620 |    238331 ns/op |   87696 B/op |   609 allocs/op |
| BenchmarkDenco_GithubAll       |     18355 |     64494 ns/op |   20224 B/op |   167 allocs/op |
| BenchmarkEcho_GithubAll        |     31251 |     38479 ns/op |       0 B/op |     0 allocs/op |
| BenchmarkGocraftWeb_GithubAll  |      4117 |    300062 ns/op |  131656 B/op |  1686 allocs/op |
| BenchmarkGoji_GithubAll        |      3274 |    416158 ns/op |   56112 B/op |   334 allocs/op |
| BenchmarkGojiv2_GithubAll      |      1402 |    870518 ns/op |  352720 B/op |  4321 allocs/op |
| BenchmarkGoJsonRest_GithubAll  |      2976 |    401507 ns/op |  134371 B/op |  2737 allocs/op |
| BenchmarkGoRestful_GithubAll   |       410 |   2913158 ns/op |  910144 B/op |  2938 allocs/op |
| BenchmarkGorillaMux_GithubAll  |       346 |   3384987 ns/op |  251650 B/op |  1994 allocs/op |
| BenchmarkGowwwRouter_GithubAll |     10000 |    143025 ns/op |   72144 B/op |   501 allocs/op |
| BenchmarkHttpRouter_GithubAll  |     55938 |     21360 ns/op |       0 B/op |     0 allocs/op |
| BenchmarkHttpTreeMux_GithubAll |     10000 |    153944 ns/op |   65856 B/op |   671 allocs/op |
| BenchmarkKocha_GithubAll       |     10000 |    106315 ns/op |   23304 B/op |   843 allocs/op |
| BenchmarkLARS_GithubAll        |     47779 |     25084 ns/op |       0 B/op |     0 allocs/op |
| BenchmarkMacaron_GithubAll     |      3266 |    371907 ns/op |  149409 B/op |  1624 allocs/op |
| BenchmarkMartini_GithubAll     |       331 |   3444706 ns/op |  226551 B/op |  2325 allocs/op |
| BenchmarkPat_GithubAll         |       273 |   4381818 ns/op | 1483152 B/op | 26963 allocs/op |
| BenchmarkPossum_GithubAll      |     10000 |    164367 ns/op |   84448 B/op |   609 allocs/op |
| BenchmarkR2router_GithubAll    |     10000 |    160220 ns/op |   77328 B/op |   979 allocs/op |
| BenchmarkRivet_GithubAll       |     14625 |     82453 ns/op |   16272 B/op |   167 allocs/op |
| BenchmarkTango_GithubAll       |      6255 |    279611 ns/op |   63826 B/op |  1618 allocs/op |
| BenchmarkTigerTonic_GithubAll  |      2008 |    687874 ns/op |  193856 B/op |  4474 allocs/op |
| BenchmarkTraffic_GithubAll     |       355 |   3478508 ns/op |  820744 B/op | 14114 allocs/op |
| BenchmarkVulcan_GithubAll      |      6885 |    193333 ns/op |   19894 B/op |   609 allocs/op |

- (1): Total Repetitions achieved in constant time, higher means more confident result
- (2): Single Repetition Duration (ns/op), lower is better
- (3): Heap Memory (B/op), lower is better
- (4): Average Allocations per Repetition (allocs/op), lower is better

## Users

Awesome project lists using [Bucky](https://github.com/kumarabd/bucky) web framework.

* [krypton](https://github.com/kumarabd/krypton): A foundation framework written in Go.

## The gopher

The gopher used here was created using [Gopherize.me](https://gopherize.me/). WebGo stays out of developers' way, so sitback and enjoy a cup of coffee like this gopher.

