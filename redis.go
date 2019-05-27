package main

import (
	"errors"
	"log"
	"time"
)

func locker(key, token string, expire time.Duration) (err error) {
	const maxCount = 5
	count := 0
	for {
		setResult, err := redisClient.SetNX(key, token, expire).Result()
		if err != nil {
			return err
		}
		if setResult {
			log.Println("lock--->")
			break
		}
		if count++; count > maxCount {
			return errors.New("timeout")
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func unlock(key, token string) (int64, error) {
	const script = `
	if redis.call("get",KEYS[1]) == ARGV[1]
	then
		return redis.call("del",KEYS[1])
	else
		return 0
	end
	`
	res, err := redisClient.Eval(script, []string{key}, token).Result()
	if err != nil {
		return 0, err
	}
	log.Println("<---unlock", res.(int64))
	return res.(int64), nil
}
