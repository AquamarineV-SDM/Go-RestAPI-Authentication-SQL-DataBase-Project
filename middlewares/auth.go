package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"restapiaqua.com/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	//ฝากค่า userId เก็บไว้ในกระเป๋าของ context เพื่อให้ func ถัดไปหยิบไปใช้ได้
	context.Set("userId", userId)

	//allow ให้วิ่งไปทำ func หลัก เช่น createEvent ต่อได้
	context.Next()
}
