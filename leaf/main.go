package main

import (
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	h := &LeafSvr{}
	h.init()

	h.Run()
}
