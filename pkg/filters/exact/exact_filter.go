package exactfilter

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/csyourui/wechat_server/pkg/filters"
	"github.com/csyourui/wechat_server/pkg/utils"
	"time"

	"github.com/go-redis/redis/v7"
)

type ExactFilter struct {
	tag    string
	client *redis.ClusterClient
	lru    *utils.LRU
	*filters.FilterParams
}

func NewExactFilter(jobid string, client *redis.ClusterClient, param *filters.FilterParams, cacheSize int64) (filters.Filter, error) {
	lru := utils.NewLRU(uint64(cacheSize))
	return &ExactFilter{jobid, client, lru, param}, nil
}

func (filter *ExactFilter) Judge(data []byte) (exists bool, err error) {
	hashValue := md5.Sum([]byte(filter.tag + string(data)))

	_, exist := filter.lru.Get(hashValue)
	if exist {
		return true, err
	}

	success, _ := filter.client.SetNX(fmt.Sprintf("%x", hashValue), 1, time.Duration(filter.Expire)*1e9).Result()
	if !success {
		return true, err
	}
	filter.lru.Set(hashValue, 1)
	return false, nil
}

func (filter *ExactFilter) Reset(params *filters.FilterParams) (err error) {
	if params.Kind != filters.Exact {
		err = errors.New("kind not same")
		return
	}
	filter.Expire = params.Expire
	return nil
}

func (filter *ExactFilter) Euqal(params *filters.FilterParams) bool {
	return *filter.FilterParams == *params
}

func (filter *ExactFilter) Info() map[string]interface{} {
	result := make(map[string]interface{})
	result["info"] = filter.FilterParams
	result["rate"] = 0
	return result
}

func (filter *ExactFilter) Delete() {}
