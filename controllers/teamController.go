package controllers

import (
	"github.com/gin-gonic/gin"
	"neft.web/models"
	"net/http"
)

type Teams struct {
	ts models.TeamService
}

func NewTeams(ts models.TeamService) *Teams {
	return &Teams{
		ts: ts,
	}
}

func (ts *Teams) RetrieveCompleteTeam(context *gin.Context) {
	team, err := ts.ts.AllTeamByID("1")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, team)
}