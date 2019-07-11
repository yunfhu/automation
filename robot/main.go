package main

import (
	"log"
	"net/http"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

const url = "http://10.21.19.77:31380/productpage"

func visitBook(seq int) {
	log.Printf("robot %d online", seq)
	defer wg.Done()
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Panic("visit productPage got err")
		}
		log.Printf("visit productPage,got repsonse code:%s", resp.Status)
		defer resp.Body.Close()
	}
}

func main() {
	runtime.GOMAXPROCS(2)
	robotNum := 16
	wg.Add(robotNum)
	for i := 0; i < robotNum; i++ {
		go visitBook(i)
	}
	wg.Wait()

}
