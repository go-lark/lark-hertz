// Package larkhertz is Hertz middleware for go-lark
package larkhertz

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/go-lark/lark"
)

// DefaultLarkMessageKey still public for compatibility
const DefaultLarkMessageKey = "go-lark-message"

// LarkMiddleware .
type LarkMiddleware struct {
	messageKey string

	enableTokenVerification bool
	verificationToken       string

	enableEncryption bool
	encryptKey       []byte

	enableURLBinding bool
	urlPrefix        string
}

// NewLarkMiddleware .
func NewLarkMiddleware() *LarkMiddleware {
	return &LarkMiddleware{
		messageKey: DefaultLarkMessageKey,
	}
}

// WithTokenVerification .
func (opt *LarkMiddleware) WithTokenVerification(token string) *LarkMiddleware {
	opt.enableTokenVerification = true
	opt.verificationToken = token

	return opt
}

// WithEncryption .
func (opt *LarkMiddleware) WithEncryption(key string) *LarkMiddleware {
	opt.enableEncryption = true
	opt.encryptKey = lark.EncryptKey(key)

	return opt
}

// BindURLPrefix .
func (opt *LarkMiddleware) BindURLPrefix(prefix string) *LarkMiddleware {
	opt.enableURLBinding = true
	opt.urlPrefix = prefix

	return opt
}

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

// LarkChallengeHandler Lark challenge handler
func (opt *LarkMiddleware) LarkChallengeHandler() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer ctx.Next(c)

		if !opt.checkURL(ctx) {
			return
		}

		inputBody, ok := opt.getBody(ctx)
		if !ok {
			handleGetBodyError(c)
			return
		}

		var challenge lark.EventChallengeReq
		err := json.Unmarshal(inputBody, &challenge)
		if err != nil || challenge.Challenge == "" {
			return
		}

		if !opt.checkToken(ctx, challenge.Token) {
			handleCheckTokenError(c)
			return
		}

		if challenge.Type == "url_verification" {
			hlog.CtxDebugf(c, "Handling challenge: %v", challenge.Challenge)
			ctx.AbortWithStatusJSON(http.StatusOK, utils.H{
				"challenge": challenge.Challenge,
			})
		}
	}
}

func (opt *LarkMiddleware) checkURL(ctx *app.RequestContext) bool {
	if opt.enableURLBinding && string(ctx.Request.RequestURI()) != opt.urlPrefix {
		// url not match just pass
		return false
	}
	return true
}

func (opt *LarkMiddleware) getBody(ctx *app.RequestContext) ([]byte, bool) {
	body := ctx.Request.Body()
	inputBody := body
	if opt.enableEncryption {
		decryptedData, err := opt.decodeEncryptedJSON(body)
		if err != nil {
			return nil, false
		}
		inputBody = decryptedData
	}
	return inputBody, true
}

func (opt *LarkMiddleware) checkToken(ctx *app.RequestContext, token string) bool {
	if opt.enableTokenVerification && token != opt.verificationToken {
		return false
	}
	return true
}
