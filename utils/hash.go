package utils

import "golang.org/x/crypto/bcrypt"

//func convert password to Hash
func HashPassword(password string) (string, error) {
	//bcrypt.DefaultCost คือระดับความยากในการ Hash (ยิ่งมากยิ่งปลอดภัย เเต่ก็ประมวณผลนาน)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//func for check that password ตอนที่ login ตรงกับ Hash ในระบบไหม
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
