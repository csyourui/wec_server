package bloomfilter

import (
	"errors"
	"fmt"
	"github.com/csyourui/wechat_server/pkg/filters"
	"github.com/go-redis/redis/v7"
	"strconv"
)

const LimitCount = 0.2
const BloomCacheSize = 1000000

func GetTaskHashKey(prefix, id string) string {
	return prefix + ":" + id
}

type RqBloom struct {
	tag         string
	client      *redis.ClusterClient
	bloomFilter *BloomFilter
	*filters.FilterParams
	lc uint //limit count
}

func NewRqBloomFilter(jobid string, client *redis.ClusterClient, param *filters.FilterParams) (filters.Filter, error) {
	if param.N < 1000 || param.N > 10000000000 || param.P < 0.00001 || param.P > 0.3 {
		return nil, errors.New("should 1000 <= n <= 10000000000, 0.00001 <= p <= 0.3")
	}
	key := GetTaskHashKey(filters.RedisDeduplicatorPrefix, jobid)
	bloomFilter, err := NewBloomFilter(Redis, client, key, uint(param.N), param.P, BloomCacheSize)
	if err == nil {
		return &RqBloom{
			jobid,
			client,
			bloomFilter,
			param,
			EstimateCount(param.N, param.P, LimitCount),
		}, err
	}
	return nil, err

}

func (rqb *RqBloom) Append(value []byte) (err error) {
	key := GetTaskHashKey(filters.RedisDeduplicatorCount, rqb.tag)
	if err = rqb.bloomFilter.Append(value); err == nil {
		count, e := rqb.client.HIncrBy(key, rqb.tag, 1).Result()
		if e == nil && uint(count) > rqb.lc {
			err = rqb.Reset(&filters.FilterParams{
				filters.Bloom,
				rqb.N,
				rqb.P,
				0})
		} else {
			err = e
		}
	}
	return
}

func (filter *RqBloom) Euqal(params *filters.FilterParams) bool {
	return *filter.FilterParams == *params
}

func (filter *RqBloom) Info() map[string]interface{} {
	result := make(map[string]interface{})
	result["info"] = filter.FilterParams
	rate, _ := filter.GetRate()
	result["rate"] = rate
	return result
}

func (rqb *RqBloom) Exists(value []byte) (exists bool, err error) {
	return rqb.bloomFilter.Exists(value)
}

func (rqb *RqBloom) Judge(value []byte) (exists bool, err error) {
	key := GetTaskHashKey(filters.RedisDeduplicatorCount, rqb.tag)
	if exists, err = rqb.bloomFilter.ExistsAndAppend(value); err == nil && !exists {
		count, e := rqb.client.HIncrBy(key, rqb.tag, 1).Result()
		if e == nil && uint(count) > rqb.lc {
			err = rqb.Reset(&filters.FilterParams{
				filters.Bloom,
				rqb.N,
				rqb.P,
				0})
		} else {
			err = e
		}
	}
	return
}

func (rqb *RqBloom) Reset(param *filters.FilterParams) (err error) {
	if param.Kind != rqb.Kind {
		err = errors.New("kind not same")
		return
	}
	key := GetTaskHashKey(filters.RedisDeduplicatorCount, rqb.tag)
	if rqb.bloomFilter == nil {
		var bloomFilter *BloomFilter
		bloomFilter, err := NewBloomFilter(Redis, rqb.client, key, param.N, param.P, BloomCacheSize)
		if err == nil {
			rqb.bloomFilter = bloomFilter
		}
	} else {
		err = rqb.bloomFilter.Reset(param.N, param.P)
	}

	if err == nil {
		pipe := rqb.client.TxPipeline()
		pipe.HSet(key, rqb.tag, 0)
		_, err = pipe.Exec()
	}
	if err == nil {
		rqb.N = param.N
		rqb.P = param.P
		rqb.lc = EstimateCount(rqb.N, rqb.P, LimitCount)
	}

	if err != nil {
		err = errors.New(err.Error())
	}

	return err
}

func (rqb *RqBloom) GetRate() (float64, error) {
	key := GetTaskHashKey(filters.RedisDeduplicatorCount, rqb.tag)
	count, err := rqb.client.HGet(key, rqb.tag).Result()
	if err != nil {
		return 0, err
	}
	m, k := EstimateParameters(rqb.N, rqb.P)
	intcount, _ := strconv.ParseUint(count, 10, 64)
	r := EstimateRate(m, k, uint(intcount))
	return r, nil
}

func (rqb *RqBloom) GetNP() (uint, float64) {
	return rqb.N, rqb.P
}

func (rqb *RqBloom) Delete() {
	_, k := EstimateParameters(rqb.N, rqb.P)
	key := GetTaskHashKey(filters.RedisDeduplicatorPrefix, rqb.tag)
	for i := uint(0); i < k; i++ {
		rqb.client.Del(fmt.Sprintf("%s.%d", key, i+1)).Result()
	}
	rqb.client.Del(GetTaskHashKey(filters.RedisDeduplicatorCount, rqb.tag)).Result()
}
