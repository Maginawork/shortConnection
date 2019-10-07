
package shortlib

import (
	"net/http"
)

type Processor interface {
	/*
	*	基础接口
	*	参数：方法，url参数，请求体
	*	返回：需要返回的http response
	 */
	ProcessRequest(method, request_url string, params map[string]string, body []byte, w http.ResponseWriter, r *http.Request) error
}

type BaseProcessor struct {
	RedisCli      *RedisAdaptor
	Configure     *Configure
	HostName      string
	Lru           *LRU
	CountFunction CreateCountFunc
}
