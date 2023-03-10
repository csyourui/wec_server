package bloomfilter

import (
	"fmt"
	"github.com/csyourui/wechat_server/pkg/log"

	"github.com/go-redis/redis/v7"
)

type filter = string

// RedisStorage is a struct representing the Redis backend for the bloom filters.
type RedisStorage struct {
	client   *redis.ClusterClient
	bf       *BloomFilter
	filters  []filter
	hashIter uint
	size     uint
}

// NewRedisStorage creates a Redis backend storage to be used with the bloom filters.
func NewRedisStorage(client *redis.ClusterClient, bf *BloomFilter, key string, hashIter uint, size uint) (*RedisStorage, error) {
	var err error
	store := RedisStorage{client: client, bf: bf, filters: make([]filter, hashIter), hashIter: hashIter, size: size}
	keys := make([]string, hashIter)
	for i := uint(0); i < hashIter; i++ {
		keys[i] = store.getRedisKey(i + 1)
	}
	exists, err := client.Exists(keys...).Result()
	if err != nil {
		return nil, err
	}
	//redis will not be inited, if keys has existed
	if uint(exists) != store.hashIter {
		if err = store.Init(hashIter, store.size); err != nil {
			return nil, err
		}
	} else {
		store.filters = keys
	}
	return &store, nil
}

func (s RedisStorage) getRedisKey(i uint) string {
	log.Logger.Debugf("getRedisKey: {%s}.%d", s.bf.key, i)
	return fmt.Sprintf("{%s}.%d", s.bf.key, i)
}

// Init takes care of settings every bit to 0 in the Redis bitset.
func (s *RedisStorage) Init(hashIter, size uint) (err error) {
	pipe := s.client.TxPipeline()
	defer pipe.Close()
	var k uint

	for k = 0; k < hashIter; k++ {
		filterKey := s.getRedisKey(k + 1)
		pipe.Del(filterKey)
		pipe.SetBit(filterKey, int64(size), 0)
		s.filters[k] = filterKey
	}
	if cmders, err := pipe.Exec(); err != nil {
		for k = 0; k < hashIter; k++ {
			s.filters[k] = ""
			perm, err := cmders[k].(*redis.StatusCmd).Result()
			log.Logger.Infof("bloom filters set err %v %v", perm, err)
			perm2, err2 := cmders[k+1].(*redis.IntCmd).Result()
			log.Logger.Infof("bloom filters set err %v %v", perm2, err2)
		}
		pipe.Discard()
	} else {
		s.hashIter = hashIter
		s.size = size
	}
	return err
}

// Append appends the bit, which is to be saved, to the Redis.
func (s *RedisStorage) Append(value []byte) error {
	pipe := s.client.TxPipeline()
	defer pipe.Close()
	hash1, hash2 := hashValue(&value)
	for index, fkey := range s.filters {
		pipe.SetBit(fkey, int64(getOffset(hash1, hash2, uint(index+1), s.size)), 1)
	}
	_, err := pipe.Exec()
	if err != nil {
		pipe.Discard()
	}
	return err
}

// Exists checks if the given bit exists in the Redis backend.
func (s *RedisStorage) Exists(value []byte) (ret bool, err error) {
	pipe := s.client.TxPipeline()
	defer pipe.Close()
	hash1, hash2 := hashValue(&value)
	res := make([]*redis.IntCmd, s.hashIter)
	for index, fkey := range s.filters {
		res[index] = pipe.GetBit(fkey, int64(getOffset(hash1, hash2, uint(index+1), s.size)))
	}
	if _, err := pipe.Exec(); err != nil {
		pipe.Discard()
		return false, err
	}
	ret = true
	for _, v := range res {
		if v.Val() != 1 {
			ret = false
			break
		}
	}
	return ret, nil
}

// ExistsAndAppend checks and append
func (s *RedisStorage) ExistsAndAppend(value []byte) (ret bool, err error) {
	pipe := s.client.TxPipeline()
	defer pipe.Close()
	hash1, hash2 := hashValue(&value)
	res := make([]*redis.IntCmd, s.hashIter)
	for index, fkey := range s.filters {
		res[index] = pipe.SetBit(fkey, int64(getOffset(hash1, hash2, uint(index+1), s.size)), 1)
	}
	if _, err := pipe.Exec(); err != nil {
		pipe.Discard()
		return false, err
	}
	ret = true
	for _, v := range res {
		if v.Val() != 1 {
			ret = false
			break
		}
	}
	return ret, nil
}
