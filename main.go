package main

import (
	"os"
	"github.com/tttmaximttt/catalogservice/service"
	"github.com/cloudfoundry-community/go-cfenv"
	"fmt"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	appEnv, err := cfenv.Current()
	if err != nil {
		fmt.Println("CF Environment not detected.")
	}

	server := service.NewServerFromCFEnv(appEnv)
	server.Run(":" + port)
}