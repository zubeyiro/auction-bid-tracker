package main

import (
	"github.com/zubeyiro/auction-bid-tracker/src/server/data"
	"github.com/zubeyiro/auction-bid-tracker/src/server/rest"
)

func main() {
	db := data.Init()
	rest.StartRestServer(db)
}
