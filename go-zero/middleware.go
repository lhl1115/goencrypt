package go_zero

import (
	"go-api-encrypt/utils"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wumansgy/goEncrypt/aes"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type EncryptMiddleware struct {
	AesKey     string
	AesIv      string
	AppID      string
	AppSecret  string
	ApiTimeout int
}

func NewEncryptMiddleware(aesKey, aesIv, appID, appSecret string, apiTimeout int) *EncryptMiddleware {
	return &EncryptMiddleware{
		AesKey:     aesKey,
		AesIv:      aesIv,
		ApiTimeout: apiTimeout,
		AppID:      appID,
		AppSecret:  appSecret,
	}
}

func (m *EncryptMiddleware) getAppSecret(appid string) string {
	switch appid {
	case m.AppID:
		return m.AppSecret
	}
	return ""
}

func (m *EncryptMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appid := r.Header.Get("app_id")
		timestamp := r.Header.Get("timestamp")
		signature := r.Header.Get("signature")
		if appid == "" || timestamp == "" || signature == "" {
			httpx.WriteJson(w, http.StatusBadRequest, "请求头部错误")
			return
		}

		timeInt64, _ := strconv.ParseInt(timestamp, 10, 64)
		timeRequest := time.Unix(timeInt64, 0)
		// 判断时间戳是否超时
		if int(math.Abs(time.Since(timeRequest).Seconds())) > m.ApiTimeout {
			httpx.WriteJson(w, http.StatusBadRequest, "系统时间错误")
			return
		}

		// 获取参数
		paramMap := make(map[string]interface{})
		switch r.Method {
		case "POST":
			var err error
			paramMap, err = utils.GetPostJsonParams(r)
			if err != nil {
				httpx.WriteJson(w, http.StatusBadRequest, "请求参数错误")
				return
			}
		case "GET":
			query := r.URL.Query()
			for k := range query {
				paramMap[k] = r.FormValue(k)
			}
		case "OPTIONS":
			next(w, r)
			return
		default:
			next(w, r)
			return
		}

		var signText strings.Builder
		// 拼接app secret
		signText.WriteString(m.getAppSecret(appid))
		// 拼接参数 按key排序 然后以key+value拼接
		array := utils.MapKeys(paramMap)
		sort.Strings(array)
		for _, v := range array {
			if paramMap[v] == nil {
				continue
			}
			signText.WriteString(v)
			signText.WriteString(utils.AnyToString(paramMap[v]))
		}
		// 拼接时间戳
		signText.WriteString(timestamp)

		// 加密后与sign对比
		aesEncrypt, err := aes.AesCbcEncryptHex([]byte(signText.String()), []byte(m.AesKey), []byte(m.AesIv))
		if err != nil {
			httpx.WriteJson(w, http.StatusBadRequest, "签名错误")
			return
		}
		if aesEncrypt != signature {
			httpx.WriteJson(w, http.StatusBadRequest, "签名错误")
			return
		}
		next(w, r)
	}
}
