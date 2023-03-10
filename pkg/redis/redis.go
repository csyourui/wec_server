package redis

import (
	"fmt"
	"github.com/csyourui/wechat_server/pkg/log"
	"os"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

// NewRedisClient create redis client and connect
func NewRedisSingleClient(conf *viper.Viper) (*redis.Client, error) {

	var err error
	client := redis.NewClient(
		&redis.Options{
			Addr:     conf.GetString("redis.addr"),
			Password: conf.GetString("redis.password"),
			DB:       conf.GetInt("redis.db"),
		})

	_, err = client.Ping().Result()
	if err == nil {
		log.Logger.Infof("connect redis succ %s\n", client.Options().Addr)
		debug := false
		if conf.IsSet("debug") {
			debug = conf.GetBool("debug")
		}
		if debug {
			if info, err := client.Info().Result(); err == nil {
				fmt.Fprintf(os.Stderr, "Redis Info:\n%s\n", info)
			}
		}
	} else {
		log.Logger.Errorf("connect redis fail %s %s", client.Options().Addr, err)
	}
	return client, err
}

// NewRedisClient create redis client and connect
func NewRedisClusterClient(conf *viper.Viper) (*redis.ClusterClient, error) {
	var err error
	client := redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:    []string{conf.GetString("redis.addr")},
			Password: conf.GetString("redis.password"),
		})

	_, err = client.Ping().Result()
	if err == nil {
		log.Logger.Infof("connect redis succ %s\n", client.Options().Addrs)
		debug := false
		if conf.IsSet("debug") {
			debug = conf.GetBool("debug")
		}
		if debug {
			if info, err := client.Info().Result(); err == nil {
				fmt.Fprintf(os.Stderr, "Redis Info:\n%s\n", info)
			}
		}
	} else {
		log.Logger.Errorf("connect redis fail %s %s", client.Options().Addrs, err)
	}
	return client, err
}

// ParsedRedisInfo TODO
type ParsedRedisInfo map[string]map[string]string

// ParseRedisInfo TODO
func ParseRedisInfo(info string) ParsedRedisInfo {
	out := ParsedRedisInfo{}
	var b map[string]string
	var e int
	for len(info) > 0 {
		if strings.HasPrefix(info, "# ") {
			e = strings.Index(info, "\n")
			if e > 0 {
				s := strings.Index(info[:e], " ")
				if s > 0 {
					k := strings.ToLower(strings.TrimSpace(info[s:e]))
					b = map[string]string{}
					out[k] = b
				}
			}
		} else {
			e = strings.Index(info, "\n")
			if e > 0 {
				t := strings.Index(info[:e], ":")
				if t > 0 {
					key := strings.TrimSpace(info[:t])
					value := strings.TrimSpace(info[t+1 : e])
					b[key] = value
				}
			}
		}
		if e < 0 {
			break
		}
		info = info[e+1:]
	}
	return out
}

// ExtractRedisInfo TODO
func ExtractRedisInfo(info string, key string) (string, string) {
	t := strings.Index(info, key)
	if t < 0 {
		return "", ""
	}
	info = info[t:]
	t = strings.Index(info, ":")
	if t < 0 {
		return "", ""
	}
	k := info[:t]
	info = info[t+1:]

	t = strings.Index(info, "\r\n")
	if t < 0 {
		return "", ""
	}
	v := info[:t]
	return k, v
}

type scanFunc func(key string, cursor uint64, match string, count int64) *redis.ScanCmd

// ScanProcessFunc TODO
type ScanProcessFunc func(keys []string) error

// Scanner TODO
type ClusterScanner struct {
	client *redis.ClusterClient
	key    string
	cursor uint64
	match  string
	count  int64
	scan   scanFunc
}

// Scanner TODO
type SingleScanner struct {
	client *redis.Client
	key    string
	cursor uint64
	match  string
	count  int64
	scan   scanFunc
}

func NewRedisScanner(client interface{}, scanType string, key string, match string, count int64) interface{} {
	switch client.(type) {
	case *redis.ClusterClient:
		return NewRedisClusterScanner(client.(*redis.ClusterClient), scanType, key, match, count).(*ClusterScanner)
	case *redis.Client:
		return NewRedisSingleScanner(client.(*redis.Client), scanType, key, match, count).(*SingleScanner)
	}
	return nil
}

// NewRedisSingleScanner TODO
func NewRedisSingleScanner(client *redis.Client, scanType string, key string, match string, count int64) interface{} {
	scanner := &SingleScanner{
		client: client,
		key:    key,
		cursor: 0,
		match:  match,
		count:  count,
	}
	st := strings.ToLower(scanType)
	switch st {
	case "scan":
		scanner.scan = func(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
			return client.Scan(cursor, match, count)
		}
	case "hscan":
		scanner.scan = client.HScan
	case "sscan":
		scanner.scan = client.SScan
	case "zscan":
		scanner.scan = client.ZScan
	default:
		panic(fmt.Errorf("invalid scan type %s", scanType))
	}
	return scanner
}

// NewRedisClusterScanner TODO
func NewRedisClusterScanner(client *redis.ClusterClient, scanType string, key string, match string, count int64) interface{} {
	scanner := &ClusterScanner{
		client: client,
		key:    key,
		cursor: 0,
		match:  match,
		count:  count,
	}
	st := strings.ToLower(scanType)
	switch st {
	case "scan":
		scanner.scan = func(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
			return client.Scan(cursor, match, count)
		}
	case "hscan":
		scanner.scan = client.HScan
	case "sscan":
		scanner.scan = client.SScan
	case "zscan":
		scanner.scan = client.ZScan
	default:
		panic(fmt.Errorf("invalid scan type %s", scanType))
	}
	return scanner
}

// Scan TODO
func (scanner *SingleScanner) Scan(f ScanProcessFunc) (err error) {
	for {
		var keys []string
		var cursor uint64
		keys, cursor, err = scanner.scan(scanner.key, scanner.cursor, scanner.match, scanner.count).Result()
		if err == nil {
			err = f(keys)
			if err != nil {
				break
			}
		} else {
			break
		}
		scanner.cursor = cursor
		if scanner.cursor == 0 {
			break
		}
	}
	return
}

// Scan TODO
func (scanner *ClusterScanner) Scan(f ScanProcessFunc) (err error) {
	for {
		var keys []string
		var cursor uint64
		keys, cursor, err = scanner.scan(scanner.key, scanner.cursor, scanner.match, scanner.count).Result()
		if err == nil {
			err = f(keys)
			if err != nil {
				break
			}
		} else {
			break
		}
		scanner.cursor = cursor
		if scanner.cursor == 0 {
			break
		}
	}
	return
}
