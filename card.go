package larkhertz

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-lark/lark"
)

// GetCardCallback from gin context
func (opt LarkMiddleware) GetCardCallback(c *app.RequestContext) (*lark.EventCardCallback, bool) {
	if card, ok := c.Get(opt.cardKey); ok {
		msg, ok := card.(lark.EventCardCallback)
		return &msg, ok
	}

	return nil, false
}

// LarkCardHandler card callback handler
// Encryption is automatically ignored, because it's not supported officially
func (opt LarkMiddleware) LarkCardHandler() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer ctx.Next(c)

		inputBody, ok := opt.getBody(ctx)
		if !ok {
			handleGetBodyError(c)
			return
		}

		var event lark.EventCardCallback
		err := json.Unmarshal(inputBody, &event)
		if err != nil {
			handleJSONUnmarshalError(c, err)
			return
		}
		if opt.enableTokenVerification {
			nonce := ctx.Request.Header.Get("X-Lark-Request-Nonce")
			timestamp := ctx.Request.Header.Get("X-Lark-Request-Timestamp")
			signature := ctx.Request.Header.Get("X-Lark-Signature")
			token := opt.cardSignature(nonce, timestamp, string(inputBody), opt.verificationToken)
			if signature != token {
				return
			}
		}
		ctx.Set(opt.cardKey, event)
	}
}

func (opt LarkMiddleware) cardSignature(nonce string, timestamp string, body string, token string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(token)
	b.WriteString(body)
	bs := []byte(b.String())
	h := sha256.New()
	h.Write(bs)
	bs = h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
