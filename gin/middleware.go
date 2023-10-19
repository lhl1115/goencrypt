package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/lhl1115/go-api-encrypt/utils"
	"github.com/wumansgy/goEncrypt/aes"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var AesKey string
var AesIv string
var AppID string
var AppSecret string
var ApiTimeout int

func getAppSecret(key string) string {
	switch key {
	case AppID:
		return AppSecret
	}
	return ""
}

func Encrypt(aesKey, aesIv, appID, appSecret string, apiTimeout int) gin.HandlerFunc {
	AesKey = aesKey
	AesIv = aesIv
	AppID = appID
	AppSecret = appSecret
	ApiTimeout = apiTimeout
	return func(c *gin.Context) {

		req := c.Request
		appid := req.Header.Get("app_id")
		timestamp := req.Header.Get("timestamp")
		signature := req.Header.Get("signature")
		if appid == "" || timestamp == "" || signature == "" {
			c.JSON(http.StatusBadRequest, "请求头部错误")
			c.Abort()
			return
		}

		timeInt64, _ := strconv.ParseInt(timestamp, 10, 64)
		timeRequest := time.Unix(timeInt64, 0)
		// 判断时间戳是否超时
		if int(math.Abs(time.Since(timeRequest).Seconds())) > ApiTimeout {
			c.JSON(http.StatusBadRequest, "系统时间错误")
			c.Abort()
			return
		}

		// 获取参数
		paramMap := make(map[string]interface{})
		switch req.Method {
		case "POST":
			var err error
			paramMap, err = utils.GetPostJsonParams(req)
			if err != nil {
				c.JSON(http.StatusBadRequest, "请求参数错误")
				c.Abort()
				return
			}
		case "GET":
			query := req.URL.Query()
			for k := range query {
				paramMap[k] = req.FormValue(k)
			}
		case "OPTIONS":
			c.Next()
			return
		default:
			c.Next()
			return
		}

		var signText strings.Builder
		// 拼接app secret
		signText.WriteString(getAppSecret(appid))
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
		aesEncrypt, err := aes.AesCbcEncryptHex([]byte(signText.String()), []byte(AesKey), []byte(AesIv))
		if err != nil {
			c.JSON(http.StatusBadRequest, "签名错误")
			c.Abort()
			return
		}
		if aesEncrypt != signature {
			c.JSON(http.StatusBadRequest, "签名错误")
			c.Abort()
			return
		}
		c.Next()
	}
}
