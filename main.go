package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

func help(w http.ResponseWriter, r *http.Request) {
	log.Println("help")
	fmt.Fprintln(w, "help")
}

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func locker(key, token string, expire time.Duration) (got bool, err error) {
	setResult, err := redisClient.SetNX("fetch:locker", token, 5*time.Second).Result()
	if err != nil {
		return false, err
	}
	return setResult, nil
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
	return res.(int64), nil
}

func fetch(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	id := r.FormValue("id")
	uuid, err := uuid.NewV4()
	commodity := fmt.Sprintf("commodity:%v", id)
	key := fmt.Sprintf("fetchlocker:%v", id)
	if err != nil {
		return
	}
	got, err := locker(key, string(uuid[:]), 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	if got {
		defer unlock(key, string(uuid[:]))
		log.Println(redisClient.Get(commodity).Result())
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("get a request")
		fmt.Fprintln(w, "default")
	})

	http.HandleFunc("/fetch", fetch)

	http.HandleFunc("/help", help)
	http.ListenAndServe(":10010", nil)
}
