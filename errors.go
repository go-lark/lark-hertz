package larkhertz

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func handleURLBindingError(c context.Context) {
	hlog.CtxErrorf(c, "url binding not match")
}

func handleGetBodyError(c context.Context) {
	hlog.CtxErrorf(c, "get body error")
}

func handleJSONUnmarshalError(c context.Context, err error) {
	hlog.CtxWarnf(c, "json unmarshal failed: %v", err)
}

func handleInvalidSchemaError(c context.Context) {
	hlog.CtxErrorf(c, "invalid event schema")
}

func handleCheckTokenError(c context.Context) {
	hlog.CtxErrorf(c, "token verification failed")
}
