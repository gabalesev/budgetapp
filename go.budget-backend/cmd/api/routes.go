package main

import (
	"go.budget-backend/cmd/api/handlers"
)

func (app *Application) routes(handler handlers.Handler) {
	apiGroup := app.server.Group("/api/v1")
	{
		apiGroup.POST("/register", handler.RegisterHandler)
		apiGroup.POST("/login", handler.LoginHandler)
		apiGroup.POST("/forgot_password", handler.ForgotPasswordHandler)
		apiGroup.POST("/reset_password", handler.ResetPasswordHandler)
	}

	profileRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthMiddleware)
	{
		profileRoutes.GET("/authenticated/user", handler.GetAuthenticatedUserHandler)
		profileRoutes.PATCH("/change_password", handler.UpdateUserPasswordHandler)
		profileRoutes.GET("/:email", handler.GetUserHandler)

	}

	categoriesRoutes := apiGroup.Group("/categories", app.appMiddleware.AuthMiddleware)
	{
		categoriesRoutes.GET("", handler.GetAllCategoriesHandler)
		categoriesRoutes.POST("", handler.CreateCategoryHandler)
		categoriesRoutes.DELETE("/:id", handler.DeleteCategoryByIDHandler)
		categoriesRoutes.GET("/:id", handler.GetCategoryByIDHandler)
	}

	budgetsRoutes := apiGroup.Group("/budgets", app.appMiddleware.AuthMiddleware)
	{
		budgetsRoutes.GET("", handler.GetAllBudgetsHandler)
		budgetsRoutes.POST("", handler.CreateBudgetHandler)
		budgetsRoutes.PUT("/:id", handler.UpdateBudgetHandler)
		//budgetsRoutes.GET("/:id", handler.GetBudgetByIDHandler)
	}

	app.server.GET("/", handler.HealthCheck)
	app.server.GET("/health", handler.HealthCheck)

}
