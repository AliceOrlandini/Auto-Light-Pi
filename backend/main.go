package main

import (
	"os"

	"github.com/AliceOrlandini/Auto-Light-Pi/wire"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	port := os.Getenv("PORT")
	app, err := wire.InitializeServer()
	if err != nil {
		panic("failed to initialize server: " + err.Error())
	}
  app.Run("0.0.0.0:" + port)
}
