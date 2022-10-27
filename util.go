package larkhertz

import (
	"encoding/json"

	"github.com/go-lark/lark"
)

func (opt LarkMiddleware) decodeEncryptedJSON(body []byte) ([]byte, error) {
	var encryptedBody lark.EncryptedReq
	err := json.Unmarshal(body, &encryptedBody)
	if err != nil {
		return nil, err
	}
	decryptedData, err := lark.Decrypt(opt.encryptKey, encryptedBody.Encrypt)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}
