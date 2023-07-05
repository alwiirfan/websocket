package main

import (
	"log"
	"server/database"
	"server/internal/user"
	"server/router"
)

func main() {
	dbConn, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connrction %s", err)
	}

	userRepo := user.NewRepository(dbConn.GetDatabase())
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	router.InitRoute(userHandler)
	router.Start("0.0.0.0:8080")
}
