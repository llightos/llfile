package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

var RedisCtx context.Context
var RedisDB *redis.Client

func init() {
	RedisCtx = context.Background()
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "175.178.40.245:6380",
		Password: "peng1275960183@",
		DB:       0,
	})
}

func RedisVal(uid string) (val string, ok bool) {
	val = ""
	ok = true //默认是存在的
	var err error
	val, err = RedisDB.Get(uid).Result()
	if err != nil {
		if err == redis.Nil {
			ok = false //不存在
			return
		}
		fmt.Println("ABOUT REDIS:", err)
	}
	return
}

func RedisAdd(key, val string) error {
	err := RedisDB.Set(key, val, 0).Err()
	return err
}

func RedisDEL(key string) error {
	_, err := RedisDB.Del(key).Result()
	return err
}

func LPush(key string, val ...string) {
	RedisDB.LPush(key, val)
}

func RPop(key string) {
	RedisDB.RPop(key)
}

func Range(key string) []string {
	result, err := RedisDB.LRange(key, 0, -1).Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	return result
}
