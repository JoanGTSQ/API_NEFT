package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"neft.web/auth"
	"neft.web/logger"
	"neft.web/models"
)

type Users struct {
	us models.UserService
}

type LoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CompletePsswdReset struct {
	Token    string
	Password string
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		us: us,
	}
}

// GET /users
// Return all users in a JSON
func (us *Users) RetrieveAllUsers(context *gin.Context) {

	// Retrieve all users data
	users, err := us.us.GetAllUsers()
	if err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})
		return
	}

	// Close connection returning code 200 and JSON with all users
	sendAnswer(http.StatusOK, context, users)
}

// GET /user
// Obtain the remmember hash from the JWT token and return it in JSON
func (us *Users) RetrieveUser(context *gin.Context) {

	// Obtain data from JWT token
	tokenNeft := context.GetHeader("neftAuth")
	claims, err := auth.ReturnClaims(tokenNeft)

	if err != nil {
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}

	// Search the user from the claims by remmember hash
	user, err := us.us.ByRemember(claims.RemmemberHash)
	if err != nil {
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}

	// Return JSON with user and code 200
	sendAnswer(http.StatusOK, context, user)
}

// DELETE /users/:id
// Obtain user data, search by ID and delete it, return code 202
func (us *Users) DeleteUser(context *gin.Context) {

	// Try to delete the user
	if err := us.us.Delete(context.Param("id")); err != nil {
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}

	// Close connection with status 202 (resource deleted)
	sendAnswer(202, context, nil)
}

// PATCH /users
// Obtain user data, search by ID and update it
func (us *Users) UpdateUser(context *gin.Context) {

	var user models.User

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&user); err != nil {
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}

	// Try to update the user
	if err := us.us.Update(&user); err != nil {
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}

	// Close connection with status 200 (status ok)
	sendAnswer(http.StatusOK, context, nil)

}

// PUT /auth
// Obtain user data, register it in the database and return a JWT TOKEN and 201 code
func (us *Users) RegisterUser(context *gin.Context) {
	var user models.User

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&user); err != nil {
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}
	// Create user with the data received
	if err := us.us.Create(&user); err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})
		return
	}

	// Generate  JWT Token
	tokenString, err := auth.GenerateJWT(user.RememberHash, user.RolID)
	if err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})
		return
	}

	// Insert token in the header and return a 201 Code
	context.Header("neftAuth", tokenString)
	sendAnswer(http.StatusCreated, context, nil)
}

// POST /users
// Retrieve data user from body and register it in the bbdd
func (us *Users) CreateUser(context *gin.Context) {
	var user models.User

	// Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&user); err != nil {
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}
	// Create user with the data received
	if err := us.us.Create(&user); err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})
		return
	}

	// Close connection with status 201 (resource created)
	sendAnswer(http.StatusCreated, context, nil)
}

// POST /auth
// Obtain login data (email,password), authenticate it and return jwt token in header
func (us *Users) Login(context *gin.Context) {
	var form LoginStruct

	// Obtain the body in the request and parse to the LoginStruct
	if err := context.ShouldBindJSON(&form); err != nil {
		logger.Debug.Println(context.Writer)
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}

	// Try to auth with the inserted data and return an error or a user
	userAuth, err := us.us.Authenticate(form.Email, form.Password)
	if err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})
		return
	}

	// Generate  JWT Token
	tokenString, err := auth.GenerateJWT(userAuth.RememberHash, userAuth.RolID)
	if err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})
		return
	}

	// Insert token in the header and return a 200 Code
	context.Header("neftAuth", tokenString)

	sendAnswer(http.StatusOK, context, nil)
}

// GET /user/:id/recover
// Initiate the reset password from the id param and return a token
func (us *Users) InitiateReset(context *gin.Context) {

	token, err := us.us.InitiateReset(context.Param("id"))

	if err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})
		return
	}

	sendAnswer(http.StatusOK, context, token)
}

// POST /user/:id/recover
// Obtain token and new password from the body and complete the reset
func (us *Users) CompleteReset(context *gin.Context) {
	var form CompletePsswdReset

	// Obtain the body in the request and parse to the LoginStruct
	if err := context.ShouldBindJSON(&form); err != nil {
		sendAnswer(http.StatusBadRequest, context, gin.H{"error": err.Error()})
		return
	}

	// Complete the reset password from the body
	user, err := us.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})
		return
	}

	// Generate  JWT Token
	tokenString, err := auth.GenerateJWT(user.RememberHash, user.RolID)
	if err != nil {
		sendAnswer(http.StatusInternalServerError, context, gin.H{"error": err.Error()})

		context.Abort()
		return
	}

	// Insert token in the header and return a 200 Code
	context.Header("neftAuth", tokenString)
	sendAnswer(http.StatusOK, context, nil)
}
