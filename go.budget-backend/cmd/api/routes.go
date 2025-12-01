package main

import (
	"go.budget-backend/cmd/api/handlers"
)

func (app *Application) routes(handler handlers.Handler) {
	apiGroup := app.server.Group("/api/v1")
	{
		apiGroup.POST("/register", handler.RegisterHandler)
		apiGroup.POST("/login", handler.LoginHandler)
	}

	profileRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthMiddleware)
	{
		profileRoutes.GET("/authenticated/user", handler.GetAuthenticatedUserHandler)
		profileRoutes.PATCH("/change_password", handler.UpdateUserPassword)
		profileRoutes.GET("/:email", handler.GetUserHandler)

	}

	app.server.GET("/", handler.HealthCheck)
	app.server.GET("/health", handler.HealthCheck)

}
