package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"neft.web/controllers"
	"neft.web/logger"
	"neft.web/middlewares"
	"neft.web/models"
)

var (
	isProd bool

	debug   bool
	debugDB bool
	ssl     bool
	route   string
)

func init() {
	flag.BoolVar(&isProd, "isProd", false, "This will ensure all pro vars are enabled")
	flag.BoolVar(&debug, "debug", false, "This will export all stats to file log.log")
	flag.BoolVar(&debugDB, "debugdb", false, "This will enable logs of db")
	flag.BoolVar(&ssl, "ssl", false, "This will require ssl to the ddbb connection")
	flag.StringVar(&route, "route", "log.txt", "This will create the log file in the desired route")
	gin.SetMode(gin.ReleaseMode)

}

func main() {
	flag.Parse()
	logger.InitLog(debug, route, "1.1.0")
	gin.DefaultWriter = logger.Wrt
	var sslmode string
	if ssl {
		sslmode = "require"
	} else {
		sslmode = "disable"
	}
	// Create connection with DB
	logger.Debug.Println("Creating connection with DB")

	services, err := models.NewServices(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("dbDirection"),
		5432,
		os.Getenv("dbUser"),
		os.Getenv("dbPsswd"),
		os.Getenv("dbName"),
		sslmode),
		debugDB)

	if err != nil {
		logger.Error.Println("error db", err)
		os.Exit(0)
	}

	defer services.Close()

  // Auto generate new tables or modifications in every start | Use DestructiveReset() to delete all data
	services.AutoMigrate()

	// Retrieve handlers struct
	logger.Debug.Println("Creating all services handlers")
	userC := controllers.NewUsers(services.User)
	rolesC := controllers.NewRoles(services.Rol)
	teamsC := controllers.NewTeams(services.Team)
	fieldsC := controllers.NewFields(services.Field)
  devicesC := controllers.NewDevices(services.Device)
	// Generate Router
	logger.Debug.Println("Creating gin router")
	r := initRouter(userC, rolesC, teamsC, fieldsC, devicesC)

	r.Use(middlewares.CORSMiddleware())
	logger.Info.Println("Runnig server")
	r.Run(":80")
}

// Generate a router with directions and middlewares
func initRouter(userC *controllers.Users,
	rolesC *controllers.Roles,
	teamsC *controllers.Teams,
	fieldsC *controllers.Fields,
  devicesC *controllers.Devices) *gin.Engine {
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

			// FIELD

			secured.GET("/field", fieldsC.RetrieveField)
			secured.PUT("/field", fieldsC.CreateField)
			secured.DELETE("/field", fieldsC.DeleteField)
			secured.PATCH("/field", fieldsC.UpdateField)
			secured.GET("/fields", fieldsC.RetrieveAllFields)
		}
    
	}
  api.Use(devicesC.RetrieveByMac)
	return router
}
