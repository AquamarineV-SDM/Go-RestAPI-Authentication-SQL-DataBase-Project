package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"restapiaqua.com/models"
	"restapiaqua.com/utils"
)

func signup(context *gin.Context) {
	var user models.User

	//รับข้อมูล JSON จากผุ้ใช้มาเเปลงใส่ struct
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	//save into database
	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}

func login(context *gin.Context) {
	var user models.User

	//receive email and password from request
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	//valify it
	err = user.ValidateCredentials()
	if err != nil {
		//if wrong
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticate user."})
		return
	}

	//when login passed  call function Token Generating
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
		return
	}

	//send the token back to client (เเทนที่จะส่งเเค่ข้อความว่า login successfully)
	context.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}
