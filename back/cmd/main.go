package main

import (
	"fmt"
	"strconv"

	"github.com/anotherhadi/eleakxir/backend/api"
	"github.com/anotherhadi/eleakxir/backend/server"
)

func main() {
	server := server.NewServer()
	fmt.Println("Starting the server.")

	api.Init(server)

	err := server.Router.Run(":" + strconv.Itoa(server.Settings.Port))
	if err != nil {
		panic(err)
	}
}
