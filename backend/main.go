package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/bootstrap"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	godotenv.Load()
}

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	
	// configure log rotation on a file
	rotator := &lumberjack.Logger{
		Filename:   "app.log", 	// log file path
		MaxSize:    1,        	// MB
		MaxBackups: 0,         	// how many backup files
		MaxAge:     28,        	// days to retain
		Compress:   true,      	// compress old log files (optional)
	}
	handler := slog.NewJSONHandler(rotator, &slog.HandlerOptions{
		Level:    slog.LevelInfo, // this is the level at which logging starts printing
		AddSource: true,					// adds the source file and line number to the log
	})
	// makes the log defined the main logger of the application
	slog.SetDefault(slog.New(handler))

	app, err := bootstrap.InitializeServer(ctx)
	if err != nil {
		panic("failed to initialize server: " + err.Error())
	}

  app.Run("0.0.0.0:" + port)
}
