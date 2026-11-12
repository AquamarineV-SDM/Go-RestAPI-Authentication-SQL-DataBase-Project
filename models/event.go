package models

import (
	"time"

	db "restapiaqua.com/database"
	//import db มาใช้ด้วย
)

//event is blueprint for store each data event

type Event struct {
	ID          int64
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	DateTime    time.Time `json:"date_time"`
	UserID      int64     //use when build login system
}

// method for saving the data  into Database
func (e *Event) Save() error {
	//คำสั่ง sql สำหรับเพิ่มข้อมูลใหม่
	//เครื่องหมาย ? คือตัวเเทนที่เราจะเอาข้อมูลมาเสียบทีหลัง ป้องกันการโดนเเฮก (SQL injection)
	query := `
	INSERT INTO events(name, description, location, date_time, user_id) 
	VALUES (?, ?, ?, ?, ?)
	`

	//สั่งให้ผู้จัดการ (DB) เตรียมคำสั่ง sql ให้พร้อม
	//stmt stands for statement(คำสั่ง)
	stmt, err := db.DB.Prepare(query) //db => ชื่อ ที่เรา import มา DB คือตัวเเปรใน db.go
	if err != nil {
		return err
	}

	//อย่าลืมปิด statement เมื่อทำงานจบ
	defer stmt.Close()

	//สั่ง execute คำสั่งที่เตรียมไว้ โดยเอาข้อมูลจาก struct มาเสียบเเทนเครื่องหมาย ? ทีละตัวตามละดับ
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	//ขอ ID ที่ฐานข้อมุลพึ่งสร้างให้อัตโนมัติ เเล้วเอามาใส่กลับเข้าใน struct ของเรา
	id, err := result.LastInsertId()
	e.ID = id

	return err
}

// not bind with anyone (no receiver)bc we need to retrive/fetch
func GetAllEvents() ([]Event, error) {
	//คำสั่ง SQL ดึงข้อมูลทั้งหมด จากตาราง events
	query := "SELECT * FROM events" //query synnomym => order

	//ใช้ DB.Query() เมื่อต้องการผลลัพธ์กลับมาหลายๆเเถว (ต่างจาก Exec ที่ไม่คืนค่าผลลัพธ์)
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //dont forget to close reading data

	//สร้างกล่อง (Slice) เตรียมไว้ใส่ข้อมูลที่จะดึงออกมา
	var events []Event

	//rows.Next() จะวนลูปอ่านข้อมูลทีละเเถว เรื่อยๆ จนหมด
	for rows.Next() {
		var event Event

		//.Scan => pcik the data in each colum store it into struct
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}

		//เอา event ที่อ่านข้อมูลครบเเล้ว ใส่เพิ่มเข้า block
		events = append(events, event)
	}
	//add  check that loop leak from err
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	//คำสั่ง SQL ใช้ where id = ? เพื่อกรองเอาเฉพาะเเถวที่ ID ตรงกัน
	query := "SELECT * FROM events WHERE id = ?"

	// use DB.QueryRow() เพราะเรามั่นใจว่า result จะมีเเค่ 1 เเถว
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	// return pointer of event back
	return &event, nil
}

func (e Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, date_time = ?
	WHERE id = ?
	`

	//เตรียมคำสั่ง SQL
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//นำข้อมูลใน struct มาเสียบเเทนที่ ? (dont forget that the last(e.ID) for WHERE)
	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)

	return err
}

func (e Event) Delete() error {
	//query SQL ลบข้อมูลเเถวที่ id ตรงกัน
	query := "DELETE FROM events WHERE id = ?"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//นำ ID ไปเสียบเเทนเครื่องหมาย
	_, err = stmt.Exec(e.ID)
	//Exec stands for execute as do the think like "insert"
	return err
}

// Register save userId registe with this event
func (e Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err

}

func (e Event) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}
