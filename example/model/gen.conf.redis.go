package model

//! conf.redis
import (
	"errors"
	"fmt"
	"strings"

	"github.com/ezbuy/redis-orm/orm"
)

var (
	_redis_store *orm.RedisStore
)

const (
	PAIR = "pair"
	HASH = "hash"
	SET  = "set"
	ZSET = "zset"
	GEO  = "geo"
	LIST = "list"

	ERROR_SPLIT = "#-#"
)

type Object interface {
	GetClassName() string
	GetStoreType() string
	GetPrimaryName() string
	GetIndexes() []string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

func RedisSetUp(cf *RedisConfig) {
	store, err := orm.NewRedisClient(cf.Host, cf.Port, cf.Password, 0)
	if err != nil {
		panic(err)
	}
	_redis_store = store
}

func Redis() *orm.RedisStore {
	return _redis_store
}

// 处理error，把一个error变成error数组
func SplitError(err error) []error {
	ss := strings.Split(err.Error(), ERROR_SPLIT)
	result := make([]error, len(ss))
	for i, s := range ss {
		result[i] = errors.New(s)
	}
	return result
}

//! util functions
func keyOfObject(obj Object, keys ...string) string {
	if len(keys) > 0 {
		suffix := strings.Join(keys, ":")
		return fmt.Sprintf("%s:%s:object:%s", obj.GetStoreType(), obj.GetClassName(), suffix)
	}
	return keyOfClass(obj)
}

func keyOfClass(obj Object, keys ...string) string {
	switch obj.GetStoreType() {
	case PAIR:
		return pairOfClass(obj.GetClassName(), keys...)
	case HASH:
		return hashOfClass(obj.GetClassName(), keys...)
	case SET:
		return setOfClass(obj.GetClassName(), keys...)
	case ZSET:
		return zsetOfClass(obj.GetClassName(), keys...)
	case GEO:
		return geoOfClass(obj.GetClassName(), keys...)
	case LIST:
		return listOfClass(obj.GetClassName(), keys...)
	}
	return ""
}

func pairOfClass(class string, keys ...string) string {
	if len(keys) > 0 {
		suffix := strings.Join(keys, ":")
		return fmt.Sprintf("%s:%s:%s", PAIR, class, suffix)
	}
	return fmt.Sprintf("%s:%s", PAIR, class)
}

func hashOfClass(class string, keys ...string) string {
	if len(keys) > 0 {
		suffix := strings.Join(keys, ":")
		return fmt.Sprintf("%s:%s:%s", HASH, class, suffix)
	}
	return fmt.Sprintf("%s:%s", HASH, class)
}

func setOfClass(class string, keys ...string) string {
	if len(keys) > 0 {
		suffix := strings.Join(keys, ":")
		return fmt.Sprintf("%s:%s:%s", SET, class, suffix)
	}
	return fmt.Sprintf("%s:%s", SET, class)
}

func zsetOfClass(class string, keys ...string) string {
	if len(keys) > 0 {
		suffix := strings.Join(keys, ":")
		return fmt.Sprintf("%s:%s:%s", ZSET, class, suffix)
	}
	return fmt.Sprintf("%s:%s", ZSET, class)
}

func geoOfClass(class string, keys ...string) string {
	if len(keys) > 0 {
		suffix := strings.Join(keys, ":")
		return fmt.Sprintf("%s:%s:%s", GEO, class, suffix)
	}
	return fmt.Sprintf("%s:%s", GEO, class)
}

func listOfClass(class string, keys ...string) string {
	if len(keys) > 0 {
		suffix := strings.Join(keys, ":")
		return fmt.Sprintf("%s:%s:%s", LIST, class, suffix)
	}
	return fmt.Sprintf("%s:%s", LIST, class)
}
