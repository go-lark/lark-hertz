package larkhertz

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-lark/lark"
)

// GetEvent should call GetEvent if you're using EventV2
func (opt *LarkMiddleware) GetEvent(c *app.RequestContext) (*lark.EventV2, bool) {
	if message, ok := c.Get(opt.messageKey); ok {
		event, ok := message.(lark.EventV2)
		if event.Schema != "2.0" {
			return nil, false
		}
		return &event, ok
	}

	return nil, false
}

// LarkEventHandler handle lark event v2
func (opt *LarkMiddleware) LarkEventHandler() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer ctx.Next(c)

		if !opt.checkURL(ctx) {
			handleURLBindingError(c)
			return
		}

		inputBody, ok := opt.getBody(ctx)
		if !ok {
			handleGetBodyError(c)
			return
		}

		var event lark.EventV2
		err := json.Unmarshal(inputBody, &event)
		if err != nil {
			handleJSONUnmarshalError(c, err)
			return
		}
		if event.Schema == "" {
			handleInvalidSchemaError(c)
			return
		}
		if !opt.checkToken(ctx, event.Header.Token) {
			handleCheckTokenError(c)
			return
		}
		hlog.CtxDebugf(c, "Handling event", event.Header.EventType)
		ctx.Set(opt.messageKey, event)
	}
}
