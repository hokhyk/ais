package libs

import (
	"log"

	"github.com/go-redis/redis/v7"
)

//Redis 物件參數
type Redis struct {
	Client *redis.Client
}

//New 建構式
func (r Redis) New(addr, password string, db int) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping().Result()

	log.Println(pong)

	if err != nil {
		log.Panicln(err)
	}

	r.Client = client

	return &r
}

//Get 取得Redis資料
func (r *Redis) Get(key string) string {
	val, err := r.Client.Get(key).Result()
	if err != nil {
		log.Panicln(err)
		val = ""
	}
	return val
}

//Set 設定Redis資料
func (r *Redis) Set(key, value string) bool {
	err := r.Client.Set(key, value, 0).Err()
	status := true

	if err != nil {
		log.Panicln(err)
		status = false
	}

	return status
}

//Remove 移除Redis資料
func (r *Redis) Remove(key string) bool {
	err := r.Client.Del(key).Err()
	status := true

	if err != nil {
		log.Panicln(err)
		status = false
	}

	return status
}
