package db

import (
	"database/sql"
	//ใส่ _(Blank Identifier) ไว้ข้างหน้าเพื่อบอก GO ว่า “โหลดมาเฉยๆ ไม่ต้องเรียกใชตรงๆ”
	_ "modernc.org/sqlite"
)

// create DB for being datacase center
var DB *sql.DB

func InitDB() {
	var err error
	// สั่งเปิดหรือสร้าง database names api.db
	DB, err = sql.Open("sqlite", "api.db")
	if err != nil {
		panic("Could not connect to database") //if can't connect => STOP!!
	}

	//ตั้งค่าการเชื่อมต่อ
	DB.SetMaxIdleConns(10) //เปิดการเชื่อมต่อพร้อมกัน max 10 ทาง
	DB.SetMaxIdleConns(5)  //เลี้ยงการเชื่อมต่อที่ว่างไว้ 5 ทาง

	createTables()

}

func createTables() {

	//create users table (should create before event table)
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
	)
	` // UNIQUE keyword that "do not anyone use emil ซ้ำกัน"

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic("COULD not create users table: " + err.Error())
	}

	// query which create  event table
	//call SQL to create data table
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		date_time DATETIME NOT NULL,
		user_id INTEGER
	)
	`

	//สั่งให้ฐานข้อมูลรันคำสั่ง sql ด้านบน
	_, err = DB.Exec(createEventsTable)

	if err != nil {
		panic("Could not create events table:" + err.Error())
	}

	//RegistrationsTable
	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER,
			user_id INTEGER,
			FOREIGN KEY(event_id) REFERENCES events(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createRegistrationsTable)
	if err != nil {
		panic("Could not create registrations table.")
	}

}
