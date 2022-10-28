# Lark Hertz

[![build](https://github.com/go-lark/lark-hertz/actions/workflows/ci.yml/badge.svg)](https://github.com/go-lark/lark-hertz/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/go-lark/lark-hertz/branch/main/graph/badge.svg?token=MQL8MFPF2Q)](https://codecov.io/gh/go-lark/lark-hertz)

Hertz middleware for go-lark.

## Installation

```shell
go get -u github.com/go-lark/lark-hertz
```

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/go-lark/lark-hertz"
)

func main() {
	r := server.Default()
	middleware := larkhertz.NewLarkMiddleware()
	r.Use(middleware.LarkEventHandler())
	r.Use(middleware.LarkChallengeHandler())

	r.POST("/", func(c context.Context, ctx *app.RequestContext) {
		if evt, err := middleware.GetEvent(ctx); err == nil {
			if evt.Header.EventType == lark.EventTypeMessageReceived {
				if msg, err := evt.GetMessageReceived(); err == nil {
					fmt.Println(msg.Message.Content)
				}
			}
			// you may parse other events
		}
	})
	r.Spin()
}
```

### Token Verification

```go
middleware.WithTokenVerfication("asodjiaoijoi121iuhiaud")
```

### Encryption

```go
middleware.WithEncryption("1231asda")
```

### URL Binding

```go
middleware.BindURLPrefix("/abc")
```

## About

Copyright (c) go-lark Developers, 2018-2022.
