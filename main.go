package main

import (
	"ToDo/adopter/gateway/mysql"
	"ToDo/cmd/api"
)

func main() {
	mysql.Init()
	api.Execute()
}
