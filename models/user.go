package models

import (
	"errors"

	db "restapiaqua.com/database"
	"restapiaqua.com/utils"
)

type User struct {
	ID       int64
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Method for saving user data
func (u *User) Save() error {

	//hash password that was sent by client
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO users(email, password)
	VALUES(?,?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId

	return err
}

func (u *User) ValidateCredentials() error {
	//find client from email in database
	query := "SELECT id, password FROM users WHERE email = ?"

	//QueryRow because we need just row
	row := db.DB.QueryRow(query, u.Email)

	var retrievePassword string
	//ดีง ID ออกมาใส่ใน struct u.ID and ดึงรหัสผ่านที่ hash เเล้วมาใส่ตัวเเปร
	err := row.Scan(&u.ID, &retrievePassword)
	if err != nil {
		return err //if cant find an email
	}

	//เอารหัสผ่าน that client texted => u.Password ไป compare with hashedpassword in database
	// function checkPasswordHash ตัวนี้เราสร้างไว้เเล้วใน utils folder
	passwordIsvalid := utils.CheckPasswordHash(u.Password, retrievePassword)

	if !passwordIsvalid {
		return errors.New("credentials invalid")
	}

	return nil
}
