package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"neft.web/logger"
	"neft.web/models"
)

type Fields struct {
	fd models.FieldDB
}

func NewFields(fd models.FieldDB) *Fields {
	return &Fields{
		fd: fd,
	}
}

/*
// GET /field
// Obtain the ID number of field from the body and search
*/
func (fd *Fields) RetrieveField(context *gin.Context) {

	field := &models.Field{}

	// Obtain the body in the request and parse to the field
	if err := context.ShouldBindJSON(field); err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// Retrieve the entire field data by the param of the url
	err := fd.fd.ByID(field)

	if err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, field)
}

/*
// GET /field
// Obtain the ID number of field from the body and search
*/
func (fd *Fields) RetrieveAllFields(context *gin.Context) {

	// Retrieve the entire field data by the param of the url
	fields, err := fd.fd.AllFields()

	if err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, fields)
}

/*
// PATCH /field
// Retrieve the field in the body and update it
*/
func (fd *Fields) UpdateField(context *gin.Context) {

	field := &models.Field{}

	// Obtain the body in the request and parse to the field
	if err := context.ShouldBindJSON(field); err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Retrieve the entire field data by the param of the url
	if err := fd.fd.Update(field); err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.AbortWithStatus(http.StatusOK)
}

/*
// PUT /field
// Create field from the body received
*/
func (fd *Fields) CreateField(context *gin.Context) {
	var field models.Field

	// Obtain the body in the request and parse to the field
	if err := context.ShouldBindJSON(&field); err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Create field with the data received
	if err := fd.fd.Create(&field); err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Close connection with status 201 (resource created)
	context.AbortWithStatus(http.StatusCreated)
}

/*
// DELETE /field/:id
// Obtain field data, search by ID and delete it, return code 202
*/
func (fd *Fields) DeleteField(context *gin.Context) {
	field := &models.Field{}

	// Obtain the body in the request and parse to the field
	if err := context.ShouldBindJSON(field); err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// Try to delete the field
	if err := fd.fd.Delete(field); err != nil {
		logger.Warning.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Close connection with status 202 (resource deleted)
	context.AbortWithStatus(202)

}
