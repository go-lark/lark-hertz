package larkhertz

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/go-lark/lark"
	"github.com/stretchr/testify/assert"
)

func performRequest(r *server.Hertz, method, path string, body interface{}) *ut.ResponseRecorder {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	return ut.PerformRequest(r.Engine, method, path, &ut.Body{Body: buf, Len: buf.Len()})
}

func TestLarkMiddleware(t *testing.T) {
	lm := NewLarkMiddleware()
	assert.False(t, lm.enableEncryption)
	assert.Empty(t, lm.encryptKey)
	assert.False(t, lm.enableTokenVerification)
	assert.Empty(t, lm.verificationToken)
	assert.False(t, lm.enableURLBinding)
	assert.Empty(t, lm.urlPrefix)
	assert.Equal(t, DefaultLarkMessageKey, lm.messageKey)

	lm.SetMessageKey("aaa")
	assert.Equal(t, "aaa", lm.messageKey)
	lm.WithEncryption("bbb")
	assert.True(t, lm.enableEncryption)
	assert.Equal(t, lark.EncryptKey("bbb"), lm.encryptKey)
	lm.WithTokenVerification("ccc")
	assert.True(t, lm.enableTokenVerification)
	assert.Equal(t, "ccc", lm.verificationToken)
	lm.BindURLPrefix("/ddd")
	assert.True(t, lm.enableURLBinding)
	assert.Equal(t, "/ddd", lm.urlPrefix)
}
