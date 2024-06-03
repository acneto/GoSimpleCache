package main

import (
	"github.com/acneto/simple_cache/cache/domain"
	"github.com/acneto/simple_cache/cache/server"
	"log"
)

func main() {
	svr := server.NewServer(domain.NewCache(), ":3000")
	err := svr.Start()
	if err != nil {
		log.Fatalf("can't start the server %+v\n", err)
	}
}
