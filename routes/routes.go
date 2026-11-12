package routes

import (
	"github.com/gin-gonic/gin"
	"restapiaqua.com/middlewares"
)

// รับตัวแปร server เข้ามาจัดการลงทะเบียน
func RegisterRoutes(server *gin.Engine) {
	//Public routes (anyone can use without using token)
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)
	server.POST("/signup", signup)
	server.POST("/login", login)

	//VIP route need to verify b Middleware
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)
}
