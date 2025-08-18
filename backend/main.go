package main

import (
	"context"
	"os"

	"github.com/AliceOrlandini/Auto-Light-Pi/wire"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	app, err := wire.InitializeServer(ctx)
	if err != nil {
		panic("failed to initialize server: " + err.Error())
	}
  app.Run("0.0.0.0:" + port)
}
