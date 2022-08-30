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
// GET /team
// Obtain the ID number of team from the body and search
*/
func (ts *Teams) RetrieveCompleteTeam(context *gin.Context) {

	team := &models.Team{}

	// Obtain the body in the request and parse to the team
	if err := context.ShouldBindJSON(team); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Retrieve the entire team data by the param of the url
	team, err := ts.ts.AllTeamByID(team)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, team)
}

/*
// PATCH /team
// Retrieve the team in the body and update it
*/
func (ts *Teams) UpdateTeam(context *gin.Context) {

	team := &models.Team{}

	// Obtain the body in the request and parse to the team
	if err := context.ShouldBindJSON(team); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Retrieve the entire team data by the param of the url
	if err := ts.ts.Update(team); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.AbortWithStatus(http.StatusOK)
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

/*
// DELETE /team/:id
// Obtain team data, search by ID and delete it, return code 202
*/
func (ts *Teams) DeleteTeam(context *gin.Context) {
	team := &models.Team{}

	// Obtain the body in the request and parse to the team
	if err := context.ShouldBindJSON(team); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// Try to delete the team
	if err := ts.ts.Delete(team); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Close connection with status 202 (resource deleted)
	context.AbortWithStatus(202)

}

/*
// GET /teams
// Return a JSON with all teams and preloaded users
*/
func (ts *Teams) RetrieveAllTeams(context *gin.Context) {

	teams, err := ts.ts.AllTeams()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, teams)
}
