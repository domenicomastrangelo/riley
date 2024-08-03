package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"riley/internal/config"
	"riley/internal/handlers"
	"riley/internal/handlers/middlewares"
	"riley/internal/sql"
)

func main() {
	sqlDatabase := sql.Connect(config.LoadConfig())
	defer func() {
		err := sqlDatabase.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	hndl := handlers.Handler{
		SQLDatabase: sqlDatabase,
		Config:      config.LoadConfig(),
		Logger:      logger,
	}

	http.Handle("GET /list", middlewares.DefaultMiddlewares(hndl.List))
	http.Handle("POST /upload", middlewares.DefaultMiddlewares(hndl.Upload))
	http.Handle("POST /download", middlewares.DefaultMiddlewares(hndl.Download))
	http.Handle("POST /delete", middlewares.DefaultMiddlewares(hndl.Delete))
	http.Handle("POST /login", middlewares.RateLimiter(hndl.Login))

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
