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
