package main

import (
	"github.com/gin-gonic/gin"
	db "restapiaqua.com/database"
	"restapiaqua.com/routes"
	// นำเข้าแพ็กเกจ routes ที่เพิ่งสร้างใหม่
)

// valiable simulator that store tempolary event

func main() {
	db.InitDB() //สั่งเปิดฐานข้อมูลเเละสร้างตารางก่อนเป็นอันดับเเรก

	//create Server
	server := gin.Default()

	routes.RegisterRoutes(server)

	//call Server port at "8080"
	server.Run(":8080")
}
