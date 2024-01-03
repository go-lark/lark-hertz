// Package larkhertz is Hertz middleware for go-lark
package larkhertz

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-lark/lark"
)

// DefaultLarkMessageKey compat legacy versions
// not use in this repo right now
const DefaultLarkMessageKey = "go-lark-message"

const (
	defaultLarkMessageKey = "go-lark-message"
	defaultLarkCardKey    = "go-lark-card"
)

// LarkMiddleware .
type LarkMiddleware struct {
	messageKey string
	cardKey    string

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
		messageKey: defaultLarkMessageKey,
		cardKey:    defaultLarkCardKey,
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

// SetCardKey .
func (opt *LarkMiddleware) SetCardKey(key string) *LarkMiddleware {
	opt.cardKey = key

	return opt
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
