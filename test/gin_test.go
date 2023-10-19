package test

import (
	"github.com/lhl1115/go-api-encrypt/utils"
	"github.com/wumansgy/goEncrypt/aes"
	"strconv"
	"testing"
	"time"
)

func TestGinPost(t *testing.T) {
	// post请求
	url := "http://localhost:8082/pong"
	headers := map[string]string{
		"app_id": utils.AppID,
	}
	headers["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)

	signText := utils.AppSecret + headers["timestamp"]

	aesEncrypt, err := aes.AesCbcEncryptHex([]byte(signText), []byte(utils.AesKey), []byte(utils.AesIv))
	headers["signature"] = aesEncrypt

	code, json, err := utils.HTTPPostJson(url, nil, headers)
	if err != nil {
		t.Error(err)
		return
	}
	if code != 200 {
		t.Error("code == ", code)
		return
	}
	t.Log(json)
}

func TestGinGet(t *testing.T) {

}
