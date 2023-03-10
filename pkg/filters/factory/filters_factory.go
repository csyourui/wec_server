package filtersfactory

import (
	"errors"
	"github.com/csyourui/wechat_server/pkg/filters"
	bloomfilter "github.com/csyourui/wechat_server/pkg/filters/bloom"
	exactfilter "github.com/csyourui/wechat_server/pkg/filters/exact"
	"github.com/go-redis/redis/v7"
)

func CreateFilter(jobid string, client *redis.ClusterClient, param *filters.FilterParams, cacheSize int64) (filters.Filter, error) {
	switch param.Kind {
	case filters.None:
		return filters.NewDefaultFilter()
	case filters.Bloom:
		return bloomfilter.NewRqBloomFilter(jobid, client, param)
	case filters.Exact:
		return exactfilter.NewExactFilter(jobid, client, param, cacheSize)
	}
	return nil, errors.New("wrong param kind")
}
