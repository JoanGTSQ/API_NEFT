package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"neft.web/models"
)

type Teams struct {
	ts models.TeamService
}

func NewTeams(ts models.TeamService) *Teams {
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
