package larkhertz

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-lark/lark"
)

// SetMessageKey .
func (opt *LarkMiddleware) SetMessageKey(key string) *LarkMiddleware {
	opt.messageKey = key

	return opt
}

// GetMessage from hertz context
func (opt *LarkMiddleware) GetMessage(c *app.RequestContext) (msg *lark.EventMessage, ok bool) {
	if message, ok := c.Get(opt.messageKey); ok {
		msg, ok := message.(lark.EventMessage)
		return &msg, ok
	}

	return nil, false
}

// LarkMessageHandler Lark message handler
func (opt *LarkMiddleware) LarkMessageHandler() app.HandlerFunc {
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

		var message lark.EventMessage
		err := json.Unmarshal(inputBody, &message)
		if err != nil {
			handleJSONUnmarshalError(c, err)
			return
		}

		if !opt.checkToken(ctx, message.Token) {
			handleCheckTokenError(c)
			return
		}
		hlog.CtxDebugf(c, "Handling message: %v", message.EventType)
		ctx.Set(opt.messageKey, message)
	}
}
