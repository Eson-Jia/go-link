package main

import (
	"log"
	"time"
)

func locker(key, token string, expire time.Duration, cancelCh <-chan struct{}) (reponseChan chan error) {
	ticker := time.Tick(1 * time.Second)
	reponseChan = make(chan error)
	go func() {
		for {
			setResult, err := redisClient.SetNX(key, token, expire).Result()
			if err != nil {
				log.Fatal(err)
				reponseChan <- err
				break
			}
			if setResult {
				reponseChan <- nil
				break
			}
			select {
			case <-ticker:
				log.Println("tick")
			case <-cancelCh:
				close(reponseChan)
				return
			}
		}
	}()
	return reponseChan
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
