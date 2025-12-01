package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.budget-backend/cmd/api/handlers"
	"go.budget-backend/cmd/api/middlewares"
	"go.budget-backend/common"
	"go.budget-backend/internal/mailer"
)

type Application struct {
	logger        echo.Logger
	server        *echo.Echo
	handler       handlers.Handler
	appMiddleware middlewares.AppMiddleware
}

func initConfig(e *echo.Echo) {
	paths := []string{
		".env",
		filepath.Join("internal", ".env"),
	}

	var loaded bool
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err != nil {
				e.Logger.Fatalf("Failed loading .env from %s: %v", p, err)
			}
			e.Logger.Infof("Loaded .env from %s", p)
			loaded = true
			break
		}
	}

	if !loaded {
		e.Logger.Print("No .env found in expected locations; continuing without it")
	}
}

func main() {
	e := echo.New()

	initConfig(e)

	db, err := common.NewMySQL()
	if err != nil {
		e.Logger.Fatal("Error creating database. %s", err)
	}

	appMailer := mailer.NewMailer(e.Logger)
	h := handlers.Handler{
		DB:     db,
		Logger: e.Logger,
		Mailer: appMailer,
	}
	appMiddleware := middlewares.AppMiddleware{
		Logger: e.Logger,
		DB:     db,
	}

	app := Application{
		logger:        e.Logger,
		server:        e,
		handler:       h,
		appMiddleware: appMiddleware,
	}
	//e.Use(middleware.Logger(), middlewares.CustomMiddleware)
	app.routes(h)
	fmt.Println(app)

	port := os.Getenv("APP_PORT")
	e.Logger.Info("Port: %s", port)
	appAddress := fmt.Sprintf("localhost:%s", port)
	e.Logger.Fatal(e.Start(appAddress))
}
