package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"neft.web/controllers"
	"neft.web/middlewares"
	"neft.web/models"
)

func main() {

	// Create connection with DB
	services, err := models.NewServices(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("dbDirection"),
		5432,
		os.Getenv("dbUser"),
		os.Getenv("dbPsswd"),
		os.Getenv("dbName")))

	if err != nil {
		fmt.Println("error db", err)
		os.Exit(0)
	}

	defer services.Close()

	// Auto generate new tables or modifications in every start | Use DestructiveReset() to delete all data
	services.AutoMigrate()

	// Retrieve handlers struct
	userC := controllers.NewUsers(services.User)
	rolesC := controllers.NewRoles(services.Rol)
	teamsC := controllers.NewTeams(services.Team)
	// Generate Router
	r := initRouter(userC, rolesC, teamsC)

	r.Use(middlewares.CORSMiddleware())
	r.Run()
}

// Generate a router with directions and middlewares
func initRouter(userC *controllers.Users, rolesC *controllers.Roles, teamsC *controllers.Teams) *gin.Engine {
	router := gin.Default()

	api := router.Group("/v1")
	{
		api.PUT("/auth", userC.RegisterUser)
		api.POST("/auth", userC.Login)

		secured := api.Group("/secured").Use(middlewares.RequireAuth())
		{
			// USER
			secured.GET("/user", userC.RetrieveUser)
			secured.POST("/user", userC.CreateUser)
			secured.PATCH("/user", userC.UpdateUser)
			secured.DELETE("/user/:id", userC.DeleteUser)
			secured.GET("/user/:id/recover", userC.InitiateReset)
			secured.POST("/user/:id/recover", userC.CompleteReset)
			secured.GET("/users", userC.RetrieveAllUsers)

			// ROL
			secured.GET("/roles", rolesC.RetrieveAllRoles)
			secured.GET("/roleUser", rolesC.RetrieveUsersOfRol)

			// TEAM
			secured.GET("/team", teamsC.RetrieveCompleteTeam)
			secured.PUT("/team", teamsC.CreateTeam)
			secured.PATCH("/team", teamsC.UpdateTeam)
			secured.DELETE("/team", teamsC.DeleteTeam)
			secured.GET("/teams", teamsC.RetrieveAllTeams)
		}
	}
	return router
}
