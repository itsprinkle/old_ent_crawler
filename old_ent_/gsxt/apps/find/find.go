package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gsxt/credit"
	"gsxt/gsxt/creditd"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	name, keyword := os.Args[1], os.Args[2]
	var num = 1
	var err error
	for {
		if len(os.Args) > 3 {
			num, err = strconv.Atoi(os.Args[3])
			if err != nil {
				num = 1
			}
		}
		credit := credit.MustGet(name)
		_, infos, err := creditd.Get(credit, keyword, 3, 2, num)
		if err != nil {
			log.Fatal(err)
		}
		for _, info := range infos {
			body, err := json.Marshal(info)
			if err == nil {
				println(string(body))
			}
			println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			body, err = json.Marshal(creditd.ToV1(info))
			if err == nil {
				println(string(body))
			}
		}
		time.Sleep(time.Second * 10)
	}
}
