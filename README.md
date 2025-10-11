# Lark Hertz

[![build](https://github.com/go-lark/lark-hertz/actions/workflows/build.yml/badge.svg)](https://github.com/go-lark/lark-hertz/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/go-lark/lark-hertz/branch/main/graph/badge.svg?token=MQL8MFPF2Q)](https://codecov.io/gh/go-lark/lark-hertz)

Hertz Middlewares for go-lark.

## Middlewares

- `LarkChallengeHandler`: URL challenge for general events and card callback
- `LarkEventHandler`: Incoming events (schema 2.0)
- `LarkCardHandler`: Card callback
- ~~`LarkMessageHandler`~~: (Legacy) Incoming message event (schema 1.0)

## Installation

```shell
go get -u github.com/go-lark/lark-hertz/v2
```

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
    "github.com/go-lark/lark/v2"
	larkhertz "github.com/go-lark/lark-hertz/v2"
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
            if evt, ok := middleware.GetEvent(ctx); ok { // => returns `*lark.EventV2`
                if evt.Header.EventType == lark.EventTypeMessageReceived {
                    // message received event
                    // you may also parse other events
                    if msg, err := evt.GetMessageReceived(); err == nil {
                        fmt.Println(msg.Message.Content)
                    }
                }
            }
        })
    }

    // card callback only
    cardGroup := r.Group("/card")
    {
        cardGroup.Use(middleware.LarkCardHandler())
        cardGroup.POST("/callback", func(c context.Context, ctx *app.RequestContext) {
            if card, ok := middleware.GetCardCallback(ctx); ok { // => returns `*lark.EventCardCallback`
            }
        })
    }

	r.Spin()
}
```

Example: [examples/hertz-middleware](https://github.com/go-lark/examples/tree/v2/hertz-middleware)

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
middleware.WithTokenVerfication("yourVerificationToken")
```

### Encryption

> Notice: encryption is not available for card callback, due to restriction from Lark Open Platform.

```go
middleware.WithEncryption("yourEncryptionKey")
```

### URL Binding

Only bind specific URL for events:

```go
middleware.BindURLPrefix("/abc")
```

## About

Copyright (c) go-lark Developers, 2018-2025.
