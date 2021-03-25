package rredis

/*
对redis进行简单的封装，实现一次实例，永久使用。
使用时不需要关心连接池和重连的问题。
*/

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"

	"github.com/dawei101/gor/rconfig"
	"github.com/dawei101/gor/rlog"
)

const (
	Nil     = redis.Nil
	DefName = "default"
)

var _redizz map[string]*redis.Client = make(map[string]*redis.Client)
var _redis_mu sync.RWMutex

type Config struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

//
// 注册一个为name的redis实例
//
// 建立无密码连接sample
//
//		Reg("name", "localhost:6379", "", 0)
//
func Reg(name, addr, password string, db int) error {
	ins := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if pong, err := ins.Ping(context.Background()).Result(); err != nil {
		rlog.Error(context.Background(), "could not connect to redis:", addr, db, "pong:", pong)
		return err
	}
	_redis_mu.Lock()
	defer _redis_mu.Unlock()
	_redizz[name] = ins
	return nil
}

//
// 返回 *redis.Client, 参见：  https://github.com/go-redis/redis
//
func Redis(name string) *redis.Client {
	_redis_mu.RLock()
	defer _redis_mu.RUnlock()
	ins, ok := _redizz[name]
	if !ok {
		panic("no redis instance named:" + name)
	}
	return ins
}

//
// 返回 *redis.Client, 参见：  https://github.com/go-redis/redis
//
func DefRedis() *redis.Client {
	return Redis(DefName)
}

func loadConfig() {
	configs := map[string]Config{}
	rconfig.DefConf().ValTo("rredis", &configs)

	for name, config := range configs {
		if err := Reg(name, config.Addr, config.Password, config.DB); err != nil {
			panic(err)
		}
	}
}
