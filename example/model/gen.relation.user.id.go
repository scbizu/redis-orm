package model

import (
	"fmt"
	"github.com/ezbuy/redis-orm/orm"
	redis "gopkg.in/redis.v5"
	"strings"
	"time"
)

var (
	_ time.Time
	_ fmt.Formatter
	_ strings.Reader
	_ orm.VSet
)

//! relation
type UserId struct {
	Key   string `db:"key" json:"key"`
	Value int32  `db:"value" json:"value"`
}

func (relation *UserId) GetClassName() string {
	return "UserId"
}

func (relation *UserId) GetIndexes() []string {
	idx := []string{}
	return idx
}

func (relation *UserId) GetStoreType() string {
	return "list"
}

type _UserIdRedisMgr struct {
	*orm.RedisStore
}

func UserIdRedisMgr(stores ...*orm.RedisStore) *_UserIdRedisMgr {
	if len(stores) > 0 {
		return &_UserIdRedisMgr{stores[0]}
	}
	return &_UserIdRedisMgr{_redis_store}
}

func (m *_UserIdRedisMgr) NewUserId(key string) *UserId {
	return &UserId{
		Key: key,
	}
}

//! pipeline
type _UserIdRedisPipeline struct {
	*redis.Pipeline
	Err error
}

func (m *_UserIdRedisMgr) BeginPipeline(pipes ...*redis.Pipeline) *_UserIdRedisPipeline {
	if len(pipes) > 0 {
		return &_UserIdRedisPipeline{pipes[0], nil}
	}
	return &_UserIdRedisPipeline{m.Pipeline(), nil}
}

//! redis relation list
func (m *_UserIdRedisMgr) ListLPush(relation *UserId) error {
	return m.LPush(listOfClass("UserId", "UserId", relation.Key), relation.Value).Err()
}

func (m *_UserIdRedisMgr) ListRPush(relation *UserId) error {
	return m.RPush(listOfClass("UserId", "UserId", relation.Key), relation.Value).Err()
}

func (m *_UserIdRedisMgr) ListLPop(key string) (*UserId, error) {
	str, err := m.LPop(listOfClass("UserId", "UserId", key)).Result()
	if err != nil {
		return nil, err
	}

	relation := m.NewUserId(key)
	if err := orm.StringScan(str, &relation.Value); err != nil {
		return nil, err
	}

	return relation, nil
}

func (m *_UserIdRedisMgr) ListRPop(key string) (*UserId, error) {
	str, err := m.RPop(listOfClass("UserId", "UserId", key)).Result()
	if err != nil {
		return nil, err
	}

	relation := m.NewUserId(key)
	if err := orm.StringScan(str, &relation.Value); err != nil {
		return nil, err
	}

	return relation, nil
}

func (m *_UserIdRedisMgr) ListLRange(key string, start, stop int64) ([]*UserId, error) {
	strs, err := m.LRange(listOfClass("UserId", "UserId", key), start, stop).Result()
	if err != nil {
		return nil, err
	}

	relations := make([]*UserId, 0, len(strs))
	for _, str := range strs {
		relation := m.NewUserId(key)
		if err := orm.StringScan(str, &relation.Value); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}
	return relations, nil
}

func (m *_UserIdRedisMgr) ListLRem(relation *UserId) error {
	return m.LRem(listOfClass("UserId", "UserId", relation.Key), 0, relation.Value).Err()
}

func (m *_UserIdRedisMgr) ListLLen(key string) (int64, error) {
	return m.LLen(listOfClass("UserId", "UserId", key)).Result()
}

func (m *_UserIdRedisMgr) ListLDel(key string) error {
	return m.Del(listOfClass("UserId", "UserId", key)).Err()
}

func (m *_UserIdRedisMgr) Clear() error {
	strs, err := m.Keys(listOfClass("UserId", "UserId", "*")).Result()
	if err != nil {
		return err
	}
	if len(strs) > 0 {
		return m.Del(strs...).Err()
	}
	return nil
}

func (m *_UserIdRedisMgr) Load(db DBFetcher) error {

	return fmt.Errorf("yaml importSQL unset.")

}

func (m *_UserIdRedisMgr) AddBySQL(db DBFetcher, sql string, args ...interface{}) error {
	objs, err := db.FetchBySQL(sql, args...)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if err := m.ListLPush(obj.(*UserId)); err != nil {
			return err
		}
	}

	return nil
}
func (m *_UserIdRedisMgr) DelBySQL(db DBFetcher, sql string, args ...interface{}) error {
	objs, err := db.FetchBySQL(sql, args...)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if err := m.ListLRem(obj.(*UserId)); err != nil {
			return err
		}
	}
	return nil
}

type _UserIdDBMgr struct {
	db orm.DB
}

func UserIdDBMgr(db orm.DB) *_UserIdDBMgr {
	if db == nil {
		panic(fmt.Errorf("UserIdDBMgr init need db"))
	}
	return &_UserIdDBMgr{db: db}
}

func (m *_UserIdDBMgr) FetchBySQL(q string, args ...interface{}) (results []interface{}, err error) {
	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("UserId fetch error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result UserId
		err = rows.Scan(&(result.Key), &(result.Value))
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("UserId fetch result error: %v", err)
	}
	return
}
