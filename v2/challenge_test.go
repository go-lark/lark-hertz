package larkhertz

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/go-lark/lark/v2"
	"github.com/stretchr/testify/assert"
)

func TestChallengePassed(t *testing.T) {
	var (
		r          = server.Default()
		middleware = NewLarkMiddleware()
	)
	r.Use(middleware.LarkChallengeHandler())
	r.POST("/", func(_ context.Context, _ *app.RequestContext) {
		// do nothing
	})

	message := lark.EventChallenge{
		Challenge: "test",
		Type:      "url_verification",
	}
	resp := performRequest(r, "POST", "/", message)
	var respData lark.EventChallenge
	if assert.NotNil(t, resp.Body) {
		json.NewDecoder(resp.Body).Decode(&respData)
		assert.Equal(t, "test", respData.Challenge)
	}
}

func TestChallengeMismatch(t *testing.T) {
	r := server.Default()
	middleware := NewLarkMiddleware().BindURLPrefix("/abc")
	r.Use(middleware.LarkChallengeHandler())
	r.POST("/", func(_ context.Context, _ *app.RequestContext) {
		// do nothing
	})

	message := lark.EventChallenge{
		Challenge: "test",
		Type:      "url_verification",
	}
	resp := performRequest(r, "POST", "/", message)
	var respData lark.EventChallenge
	if assert.NotNil(t, resp.Body) {
		err := json.NewDecoder(resp.Body).Decode(&respData)
		assert.Error(t, err)
	}
}
