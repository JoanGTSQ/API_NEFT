package controllers

import (
	"github.com/gin-gonic/gin"
	"neft.web/models"
	"net/http"
)

type Roles struct {
	rs models.RolService
}

func NewRoles(rs models.RolService) *Roles {
	return &Roles{
		rs: rs,
	}
}
func (rs *Roles) RetrieveAllRoles(context *gin.Context) {
	roles, err := rs.rs.RetrieveAll()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, roles)
}

type UserRol struct {
	Id    int `json:"id"`
	Count int `json:"count"`
}

func (rs *Roles) RetrieveUsersOfRol(context *gin.Context) {
	rolesUsers, err := rs.rs.ReturnUserRol()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, rolesUsers)
}
