package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"neft.web/models"
)

type Teams struct {
	ts models.TeamDB
}

func NewTeams(ts models.TeamDB) *Teams {
	return &Teams{
		ts: ts,
	}
}

/*
// GET /team/:id
// Obtain the ID number of team from the body and search
*/
func (ts *Teams) RetrieveCompleteTeam(context *gin.Context) {

	// Retrieve the entire team data by the param of the url
	team, err := ts.ts.AllTeamByID(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, team)
}

/*
// PUT /team
// Create team from the body received
*/
func (ts *Teams) CreateTeam(context *gin.Context) {
	var team models.Team

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&team); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Create user with the data received
	if err := ts.ts.Create(&team); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Close connection with status 201 (resource created)
	context.AbortWithStatus(http.StatusCreated)
}
