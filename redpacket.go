package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func sendRedPacket(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	amountSum := float32(amount)
	number, err := strconv.Atoi(r.FormValue("number"))
	if err != nil {
		w.WriteHeader(400)
		return
	}
	var amountList []float32
	for i := 0; i < number; i++ {
		current := rand.Float32() * float32(amountSum)
		amountList = append(amountList, current)
		amountSum -= current
	}
	log.Println(amountList)
}

func getRedPacket(w http.ResponseWriter, r *http.Request) {

}

func redPacket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getRedPacket(w, r)
	} else {
		sendRedPacket(w, r)
	}
}
