package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"restapiaqua.com/models"
)

//func for managing when someone call through this route

func getEvents(context *gin.Context) {

	events, err := models.GetAllEvents()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events. Try again later."})

		return
	}

	//send back as JSON form whith password 200 (statusOK refers to "seccess")
	context.JSON(http.StatusOK, events)
}

// func for manage data that sender want to create events
func createEvent(context *gin.Context) {
	//check บัตรผ่าน (ดึงค่าจาก header ที่ชื่อ Authorization)
	// token := context.Request.Header.Get("Authorization")

	// if token == "" {
	// 	context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
	// 	return
	// }

	// userId, err := utils.VerifyToken(token)
	// if err != nil {
	// 	fmt.Println("Token Error:", err)
	// 	context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
	// 	return
	// }

	// create emtry var from event we created before
	var event models.Event
	//get the JSON data ที่ผุ้ใช้เเนบมา (request body) จับยัดใส่ตัวเเปร event ของเรา
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}
	//หยิบค่า userId ที่ middleware เเปะป้ายฝากมาให้ (ใช้ GetInt64 เพราะเราเป็น int64)
	userId := context.GetInt64("userId")
	//create ID simulator to give this new event (เดี่ยวเราใช้ database มันจะทำให้เอง)
	event.UserID = userId

	//call method save() to save into SQLite
	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save event. Try again later."})
		fmt.Println(err) //add this to see error on the terminal
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})

}

// ดึงเเค่งานเดว
func getEvent(context *gin.Context) {
	// context.Param("id") จะดึงค่าจาก URL
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
		return
	}

	//call above func
	event, err := models.GetEventByID(eventId)
	if err != nil {

		//ให้ server ปริ๊น error ลง terminal to show for dev
		fmt.Println(err)

		//404 error check
		if err.Error() == "sql: no rows in result set" {
			context.JSON(http.StatusNotFound, gin.H{"massage": "Event not found."})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}
	context.JSON(http.StatusOK, event)
}

// function for PUT (http request)
func updateEvent(context *gin.Context) {
	//1. retrieve the ID from URL
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	//หยิบ userId ของคนที่ login มาจากกระเป๋า context
	userId := context.GetInt64("userId")

	//Is any error occur? (Use the func we have had before)
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not Found."})
		return
	}

	//check that editor(userId) เป็นคนเดียวกับเข้าของ event (event.UserID) มั้ย
	if event.UserID != userId {
		context.JSON(http.StatusNotFound, gin.H{"message": "Not authorized to update event."})
		return
	}

	//receive the data from client
	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	//บังคับให้ id ข้อมูลใหม่ ตรงกับ id ใน url (to protect clients change Id)
	updatedEvent.ID = eventId

	//update in database
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated successgully!"})

}

func deleteEvent(context *gin.Context) {
	//1. retreive ID value from URL
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse "})
		return
	}

	//pick userId from context
	userId := context.GetInt64("userId")

	// to check that the data accually have before removing
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not Found."})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to delete event."})
		return
	}

	//delete
	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete event."})
		return
	}

	//send task verify deleting
	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func registerForEvent(context *gin.Context) {
	//ดึง ID of Event from URL (ex. /events/1/register)
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	//retrieve userId of user who login from context (at Middleware เเปะมาให้)
	userId := context.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	//สั่งลงทะเบียน by call func Register() ที่เราพึ่งเขียนไว้ใน models
	err = event.Register(userId)
	if err != nil {
		fmt.Println("SQL Error:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register for event."})
		return
	}

	//ส่งข้อความตอบกลับว่า success
	context.JSON(http.StatusCreated, gin.H{"message": "Registered!"})
}

func cancelRegistration(context *gin.Context) {
	//retrieve ID
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	//retrieve login id
	userId := context.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel registration."})
		return
	}

	//reply that "canel success"
	context.JSON(http.StatusOK, gin.H{"message": "Cancelled"})
}
