package controllers

import (
	"github.com/gin-gonic/gin"
	"neft.web/auth"
	"neft.web/models"
	"net/http"
)

type Users struct {
	us models.UserService
}



type LoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
	}
  
  // Close connection returning code 200 and JSON with all users
	context.JSON(http.StatusOK, answerV1(users))
}

// GET /user 
// Obtain the remmember hash from the JWT token and return it in JSON
func (us *Users) RetrieveUser(context *gin.Context) {
  
  // Obtain data from JWT token
  tokenNeft := context.GetHeader("neftAuth")
  claims, err := auth.ReturnClaims(tokenNeft)
  
  if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
  }
  
  // Search the user from the claims by remmember hash
  user, err := us.us.ByRemember(claims.RemmemberHash)
  if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
  }

  // Return JSON with user and code 200
  context.JSON(http.StatusOK, answerV1(user))  
}


// DELETE /users
// Obtain user data, search by ID and delete it, return code 202
func (us *Users) DeleteUser(context *gin.Context) {
  var user models.User
  
  // Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

  // Try to delete the user
  if err := us.us.Delete(user.ID); err != nil{
    context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

  // Close connection with status 202 (resource deleted)
  context.JSON(202, answerV1(nil))
  
}

// POST /users
// Obtain user data, register it in the database and return a JWT TOKEN and 201 code
func (us *Users) RegisterUser(context *gin.Context) {
	var user models.User
  
  // Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
  // Create user with the data received
	if err := us.us.Create(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
  
	// Generate  JWT Token
	tokenString, err := auth.GenerateJWT(user.RememberHash, user.RolID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
  
	// Insert token in the header and return a 201 Code
  context.Header("neftAuth", tokenString)
	context.JSON(http.StatusCreated, answerV1(nil))
}

// POST /register
// Retrieve data user from body and register it in the bbdd
func (us *Users) CreateUser(context *gin.Context) {
	var user models.User

  // Obtain the body in the request and parse to the user
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
  // Create user with the data received
	if err := us.us.Create(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
  
  // Close connection with status 201 (resource created)
	context.JSON(http.StatusCreated, answerV1(nil))
}

// POST /login
// Obtain login data (email,password), authenticate it and return jwt token in header
func (us *Users) Login(context *gin.Context) {
	var form LoginStruct
  
  // Obtain the body in the request and parse to the LoginStruct
	if err := context.ShouldBindJSON(&form); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
  
  // Try to auth with the inserted data and return an error or a user 
	userAuth, err := us.us.Authenticate(form.Email, form.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
  
  // Generate  JWT Token
	tokenString, err := auth.GenerateJWT(userAuth.RememberHash, userAuth.RolID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
  
	// Insert token in the header and return a 200 Code
  context.Header("neftAuth", tokenString)
	context.JSON(http.StatusOK, answerV1(nil))
}
