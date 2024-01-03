# Lark Hertz

[![build](https://github.com/go-lark/lark-hertz/actions/workflows/ci.yml/badge.svg)](https://github.com/go-lark/lark-hertz/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/go-lark/lark-hertz/branch/main/graph/badge.svg?token=MQL8MFPF2Q)](https://codecov.io/gh/go-lark/lark-hertz)

Hertz middleware for go-lark.

## Middlewares

- `LarkChallengeHandler`: URL challenge for general events and card callback
- `LarkEventHandler`: Event v2 (schema 2.0)
- `LarkCardHandler`: Card callback
- `LarkMessageHandler`: (Legacy) Incoming message event (schema 1.0)

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
	larkhertz "github.com/go-lark/lark-hertz"
)

func main() {
	r := server.Default()
	middleware := larkhertz.NewLarkMiddleware()

    // lark server challenge
	r.Use(middleware.LarkChallengeHandler())

    // all supported events
    eventGroup := r.Group("/event")
    {
        eventGroup.Use(middleware.LarkEventHandler())
        eventGroup.POST("/", func(c context.Context, ctx *app.RequestContext) {
            if event, ok := middleware.GetEvent(e); ok { // => returns `*lark.EventV2`
            }
        })
    }

    // card callback only
    cardGroup := r.Group("/card")
    {
        cardGroup.Use(middleware.LarkCardHandler())
        cardGroup.POST("/callback", func(c context.Context, ctx *app.RequestContext) {
            if card, ok := middleware.GetCardCallback(c); ok { // => returns `*lark.EventCardCallback`
            }
        })
    }

	r.Spin()
}
```

### Event v2

The default mode is event v1. However, Lark has provided event v2 and it applied automatically to newly created bots.

To enable EventV2, we use `LarkEventHandler` instead of `LarkMessageHandler`:
```go
r.Use(middleware.LarkEventHandler())
```

Get the event (e.g. Message):
```go
r.POST("/", func(c context.Context, ctx *app.RequestContext) {
    event, ok = middleware.GetEvent(ctx)
    if evt, ok := middleware.GetEvent(c); ok { // => GetEvent instead of GetMessage
        if evt.Header.EventType == lark.EventTypeMessageReceived {
            if msg, err := evt.GetMessageReceived(); err == nil {
                fmt.Println(msg.Message.Content)
            }
            // you may have to parse other events
        }
    }
})
```

### Card Callback

We may also setup callback for card actions (e.g. button). The URL challenge part is the same.

We may use `LarkCardHandler` to handle the actions:
```go
r.Use(middleware.LarkCardHandler())
r.POST("/", func(c context.Context, ctx *app.RequestContext) {
    if event, ok = middleware.GetCardCallback(ctx); ok {
    }
})
```

### Token Verification

```go
middleware.WithTokenVerfication("asodjiaoijoi121iuhiaud")
```

### Encryption

> Notice: encryption is not available for card callback, due to restriction from Lark Open Platform.

```go
middleware.WithEncryption("1231asda")
```

### URL Binding

Only bind specific URL for events:

```go
middleware.BindURLPrefix("/abc")
```

## About

Copyright (c) go-lark Developers, 2018-2024.
