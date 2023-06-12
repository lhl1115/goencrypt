package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetPostJsonParams(r *http.Request) (map[string]interface{}, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if len(bodyBytes) == 0 {
		return m, nil
	}
	if err = r.Body.Close(); err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// MapKeys get map keys to slice
func MapKeys(m map[string]interface{}) []string {

	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

func AnyToString(value interface{}) string {
	// interface{} 转 string
	var str string
	if value == nil {
		return str
	}

	switch value := value.(type) {
	case float64:
		str = strconv.FormatFloat(value, 'f', -1, 64)
	case float32:
		str = strconv.FormatFloat(float64(value), 'f', -1, 64)
	case int:
		str = strconv.Itoa(value)
	case uint:
		str = strconv.Itoa(int(value))
	case int8:
		str = strconv.Itoa(int(value))
	case uint8:
		str = strconv.Itoa(int(value))
	case int16:
		str = strconv.Itoa(int(value))
	case uint16:
		str = strconv.Itoa(int(value))
	case int32:
		str = strconv.Itoa(int(value))
	case uint32:
		str = strconv.Itoa(int(value))
	case int64:
		str = strconv.FormatInt(value, 10)
	case uint64:
		str = strconv.FormatUint(value, 10)
	case string:
		str = value
	case time.Time:
		str = value.Format("2006-01-02 15:04:05")
	case []byte:
		str = string(value)
	default:
		bytesData, _ := json.Marshal(value)
		str = string(bytesData)
	}

	return str
}

func HTTPPostJson(apiURL string, request []byte, headers map[string]string) (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(request)))
	if err != nil {
		return "", err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
