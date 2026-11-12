package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// secretkey for regulate Token (do not disclose/reveal)
const secretKey = "supersecret"

func GenerateToken(email string, userId int64) (string, error) {
	// create Payload (ข้อมูลที่จะฝังลงไปใน token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(), //กำหนดให้ Token expired in 2hrs
	})
	//เซ็นต์ชื่อกำกับ token ด้วยกุญเเตลับของเรา เเล้วส่งกลับเป็นข้อความ string
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {
	// parse ถอดรหัส token and varify
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		//เช็คก่อนว่าใช้วิธีเซ็นชื่อ (HMAC) ตรงกับตอนที่เราสร้างไหม (ป้องกันคนเเอบใช้วิธีอื่น)
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		//ถ้าวิธีเซ็นชื่อถูกต้อง ให้ส่งกุญเเจกลับไปใช้ไขรหัส
		return []byte(secretKey), nil
	})

	//if parse ไม่ผ่าน example token มั่ว or exp
	if err != nil {
		return 0, errors.New("could not parse token")

	}

	//verify status =>(VALID)
	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return 0, errors.New("invalid token")
	}

	//เเกะเอาข้อมูล (Claims) ที่เราฝังไว้ข้างในตอนเเรกออกมา
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")

	}

	//ดีงค่า userId ออกมา
	//(notice: ตอนที่ package jwt ดึงตัวเลขกลับมา มันจะมองเป็น float64 เสมอ เราเลยต้องแปลงเป็น int64 ให้ตรงกับ model ของเรา)
	userId := int64(claims["userId"].(float64))
	return userId, nil
}
