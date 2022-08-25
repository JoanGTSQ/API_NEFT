package controllers

import (
	"github.com/gin-gonic/gin"
	"neft.web/auth"
	"neft.web/models"
	"net/http"
  "time"
  
)

type Users struct {
	us models.UserService
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		us: us,
	}
}

func (us *Users) RetrieveAllUsers(context *gin.Context) {
	users, err := us.us.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
	}
  for _, user := range users {
    user.DOB, err = time.Parse("01-02-2006", user.DOB.Format("01-02-2006"))
    if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
    return
	}
  }
	context.JSON(http.StatusOK, answerV1(users))
}
func (us *Users) RetrieveUser(context *gin.Context) {
  tokenNeft := context.GetHeader("neftAuth")
  
  claims, err := auth.ReturnClaims(tokenNeft)
  
  if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
  }

  userComplete, err := us.us.ByRemember(claims.RemmemberHash)
  if err != nil {
    context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
  }
  context.JSON(http.StatusOK, userComplete)  
}


//
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

  context.AbortWithStatus(202)
  
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
	context.AbortWithStatus(http.StatusCreated)
}
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

	context.AbortWithStatus(http.StatusCreated)
}



type LoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	context.AbortWithStatus(http.StatusOK)
}
