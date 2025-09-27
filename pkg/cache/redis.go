package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/onec-tech/bot/pkg/logger"
)

type ICache interface {
	Get(key string, in interface{}, hstore bool, field ...string) (err error)
	HGetAll(key string) (data map[string]string, err error)
	Store(key string, data interface{}, duration time.Duration, hstore bool, field ...string) (err error)
	StoreNX(key string, data interface{}, duration time.Duration) (updated bool, err error)
	Reset(key string) (err error)
	ResetMany(keys ...string) (err error)
	HRemove(key string, field ...string) error
	ResetByParent(parentKey string) (err error)
	ExistKey(key string) (ok bool)
	HExistKey(key, field string) (ok bool)
	KeyList(parentKey string) []string
	Push(key string, data []byte) error
	Pop(key string, in interface{}) error // LPOP
	LRange(key string, start, stop int64, in interface{}) error
	LRem(key string, count int64, value interface{}) (int64, error)
	RPop(key string, in interface{}) error
	BRPop(key string, timeout time.Duration, in interface{}) error
	BLPop(key string, timeout time.Duration, in interface{}) error
	HLen(key string) (uint64, error)
	LLen(key string) (uint64, error)
	Incr(key string) (uint64, error)
	Decr(key string) (uint64, error)
	RateLimit(key string, count int64, expire ...time.Duration) (ok bool)
	Expire(key string, expire time.Duration) error
	SAdd(key string, members ...interface{}) (int64, error) // O(N), N = len(members)
	SRem(key string, members ...interface{}) (int64, error) // O(N), N = len(members)
	SIsMember(key string, member interface{}) (bool, error) // O(1)
	SMembers(key string, in interface{}) error              // O(N)
}

type Cache struct {
	redis redis.UniversalClient
	log   logger.Logger
}

func (c *Cache) SetCacheImplementation(cacheClient redis.UniversalClient) {
	c.redis = cacheClient
}

func (c *Cache) SetLogger(l logger.Logger) {
	c.log = l
}

func (c *Cache) HGetAll(key string) (data map[string]string, err error) {
	return c.redis.HGetAll(context.Background(), key).Result()
}

func (c *Cache) Get(key string, in interface{}, hstore bool, field ...string) (err error) {
	var data string
	if hstore {
		if data, err = c.hGet(key, field[0]); err != nil {
			return
		}
	} else {
		if data, err = c.get(key); err != nil {
			return
		}

	}

	err = json.Unmarshal([]byte(data), &in)
	if err != nil {
		return fmt.Errorf("seriliazarion error %s", err)
	}

	return nil
}

func (c *Cache) Store(key string, data interface{}, duration time.Duration, hstore bool, field ...string) (err error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("seriliazarion error %s", err)
	}

	if hstore {
		c.setHash(key, serializedData, field[0], duration)
	} else {
		c.setSimple(key, serializedData, duration)
	}

	return nil
}

func (c *Cache) StoreNX(key string, data interface{}, duration time.Duration) (updated bool, err error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Errorf("seriliazarion error %s", err)
	}

	updated, err = c.setSimpleNX(key, serializedData, duration)

	if err != nil {
		return updated, err
	}

	return updated, nil
}

func (c *Cache) Reset(key string) (err error) {
	_, err = c.redis.Del(context.Background(), key).Result()
	return err
}

func (c *Cache) ResetMany(keys ...string) (err error) {
	_, err = c.redis.Del(context.Background(), keys...).Result()
	return err
}

func (c *Cache) ResetByParent(parentKey string) (err error) {
	keys := c.redis.Keys(context.Background(), parentKey).Val()
	_, err = c.redis.Del(context.Background(), keys...).Result()
	return
}

func (c *Cache) KeyList(parentKey string) []string {
	return c.redis.Keys(context.Background(), parentKey).Val()
}

func (c *Cache) ExistKey(key string) (ok bool) {
	if exists := c.redis.Exists(context.Background(), key).Val(); exists > 0 {
		ok = true
		return
	}

	return
}

func (c *Cache) HExistKey(key, field string) (ok bool) {
	return c.redis.HExists(context.Background(), key, field).Val()
}

func (c *Cache) Push(key string, data []byte) error {
	cmd := c.redis.LPush(context.Background(), key, data)
	return cmd.Err()
}

func (c *Cache) Pop(key string, in interface{}) error {
	bytes, err := c.redis.LPop(context.Background(), key).Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(bytes), &in)
	if err != nil {
		return fmt.Errorf("seriliazarion error %s", err)
	}

	return nil
}

func (c *Cache) LRange(key string, start, stop int64, in interface{}) error {
	result, err := c.redis.LRange(context.Background(), key, start, stop).Result()
	if err != nil {
		return err
	}

	jsonStr := "[" + strings.Join(result, ",") + "]"

	err = json.Unmarshal([]byte(jsonStr), in)
	if err != nil {
		return fmt.Errorf("serialization error %s", err)
	}

	return nil
}

func (c *Cache) LRem(key string, count int64, value interface{}) (int64, error) {
	return c.redis.LRem(context.Background(), key, count, value).Result()
}

func (c *Cache) RPop(key string, in interface{}) error {
	bytes, err := c.redis.RPop(context.Background(), key).Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(bytes), &in)
	if err != nil {
		return fmt.Errorf("seriliazarion error %s", err)
	}

	return nil
}

func (c *Cache) BRPop(key string, timeout time.Duration, in interface{}) error {

	res, err := c.redis.BRPop(context.Background(), timeout, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(res[1]), &in)
	if err != nil {
		return fmt.Errorf("seriliazarion error %s", err)
	}

	return nil
}

func (c *Cache) BLPop(key string, timeout time.Duration, in interface{}) error {

	res, err := c.redis.BLPop(context.Background(), timeout, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(res[1]), &in)
	if err != nil {
		return fmt.Errorf("seriliazarion error %s", err)
	}

	return nil
}

func (c *Cache) HLen(key string) (uint64, error) {
	return c.redis.HLen(context.Background(), key).Uint64()
}

func (c *Cache) LLen(key string) (uint64, error) {
	return c.redis.LLen(context.Background(), key).Uint64()
}

func (c *Cache) Incr(key string) (uint64, error) {
	return c.redis.Incr(context.Background(), key).Uint64()
}

func (c *Cache) Decr(key string) (uint64, error) {
	return c.redis.Decr(context.Background(), key).Uint64()
}

func (c *Cache) Expire(key string, expire time.Duration) (err error) {
	return c.redis.Expire(context.Background(), key, expire).Err()
}

func (c *Cache) HRemove(key string, field ...string) error {
	return c.redis.HDel(context.Background(), key, field...).Err()
}

func (c *Cache) RateLimit(key string, count int64, expire ...time.Duration) (ok bool) {
	ctx := context.Background()
	res := c.redis.Incr(ctx, key)
	if res.Err() != nil {
		c.log.Errorf("RateLimit.Incr: %s", res.Err().Error())
		return false
	}

	if len(expire) > 0 {
		c.redis.Expire(ctx, key, expire[0])
	}

	if res.Val() >= count {
		return false
	}

	return true
}

func (c *Cache) SAdd(key string, members ...interface{}) (int64, error) {
	return c.redis.SAdd(context.Background(), key, members...).Result()
}

func (c *Cache) SRem(key string, members ...interface{}) (int64, error) {
	return c.redis.SRem(context.Background(), key, members...).Result()
}

func (s *Cache) SIsMember(key string, member interface{}) (bool, error) {
	return s.redis.SIsMember(context.Background(), key, member).Result()
}

func (s *Cache) SMembers(key string, in interface{}) error {
	result, err := s.redis.SMembers(context.Background(), key).Result()
	if err != nil {
		return err
	}

	for i, v := range result {
		if _, er := strconv.Atoi(v); er == nil {
			continue
		}
		result[i] = fmt.Sprintf("\"%s\"", v)
	}

	jsonStr := "[" + strings.Join(result, ",") + "]"

	err = json.Unmarshal([]byte(jsonStr), in)
	if err != nil {
		return fmt.Errorf("serialization error %s", err)
	}

	return nil
}
