package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
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

func fetch(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	id := r.FormValue("id")
	uuid := fmt.Sprint(uuid.Must(uuid.NewV4()))
	commodity := fmt.Sprintf("commodity:%v", id)
	key := fmt.Sprintf("fetchlocker:%v", id)
	if err != nil {
		return
	}
	cancelChan := make(chan struct{})
	resChan, tick := locker(key, uuid, 5*time.Second, cancelChan), time.Tick(10*time.Second)
	select {
	case resp := <-resChan:
		if resp != nil {
			w.WriteHeader(500)
			return
		}
		defer unlock(key, uuid)
		break
	case <-tick:
		log.Println("cancel")
		cancelChan <- struct{}{}
		w.WriteHeader(504)
		return
	}
	value, err := redisClient.Get(commodity).Result()
	if err != nil {
		log.Fatal("key not exist:", key, err)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	intValue--
	setRes, err := redisClient.Set(commodity, intValue, 0).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(setRes)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("get a request")
		fmt.Fprintln(w, "default")
	})

	http.HandleFunc("/fetch", fetch)
	http.HandleFunc("/redpacket", redPacket)

	http.HandleFunc("/help", help)
	log.Fatal(http.ListenAndServe(":10010", nil))
}
