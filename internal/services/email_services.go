package services

import (
	"math/rand/v2"

	"gopkg.in/gomail.v2"
)

var maxValueCode = 999999
var minValueCode = 111111

func SendCodeToEmail(email, code string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", "adnagulovadel@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Code aut")
	m.SetBody("text/plain", code)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"adnagulovadel@gmail.com",
		"gzflnfrmifbgmbij",
	)
	return d.DialAndSend(m)
}

func CreateCode() int {
	return rand.IntN(maxValueCode-minValueCode) + minValueCode
}
