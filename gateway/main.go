package main

import (
	_ "microservice/gateway/load"
	"microservice/gateway/app"
)

func main() {
	app.API()
}
