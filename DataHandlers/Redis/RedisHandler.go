package DataHandlers

import (
	"encoding/json"
	"github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
	"time"
)

type RedisHandler struct {
	client *redis.Client
	config RedisConfiguration
}

type RedisConfiguration struct {
	Db               int
	Credentials      string
	ConnectionString string
}

type CacheHandler interface {
	Init(config RedisConfiguration)
	Lock(key string) (newLock *lock.Lock, err error)
	UnLock(lockedLock *lock.Lock)
	Get(key string) (string, error)
	Set(key string, val interface{}, duration time.Duration) error
	Delete(key string) error
}


func (self *RedisHandler) init(){
	self.client = redis.NewClient(&redis.Options{
		Addr:     self.config.ConnectionString,
		Password: self.config.Credentials,
		DB:       self.config.Db})
}
func (self *RedisHandler) Init(config RedisConfiguration) {
	self.config = config

}

func (self *RedisHandler) Lock(key string) (newLock *lock.Lock, err error) {
	self.init()
	defer self.client.Close()
	newLock, err = lock.ObtainLock(self.client, key, &lock.LockOptions{
		LockTimeout: time.Duration(5) * time.Second,
		WaitRetry:   time.Duration(300) * time.Microsecond,
		WaitTimeout: time.Duration(100) * time.Microsecond,
	})
	return
}

func (self *RedisHandler) UnLock(lockedLock *lock.Lock) {
	self.init()
	defer self.client.Close()
	if lockedLock.IsLocked() {
		lockedLock.Unlock()
	}
}

func (self *RedisHandler) Get(key string) (string, error) {
	self.init()
	defer self.client.Close()
	res := self.client.Get(key)
	if res.Err() != nil {
		return "", res.Err()
	}
	if res.Val() != "" {
		return res.Val(), nil
	}
	return "", nil
}

func (self *RedisHandler) Set(key string, val interface{}, duration time.Duration) error {
	self.init()
	defer self.client.Close()
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	status := self.client.Set(key, string(jsonBytes), duration)
	return status.Err()
}

func (self *RedisHandler) Delete(key string) error {
	self.init()
	status := self.client.Del(key)
	defer self.client.Close()
	return status.Err()
}
