package main

import (
	serv "main/pkg/server"
)

func main() {
	serv.CreateServer().ListenAndServe()
}
