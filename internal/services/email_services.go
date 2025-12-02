package services

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

var maxValueCode = 999999
var minValueCode = 111111

type EmailService interface {
	SendCodeToEmail(email, code string) error
	CreateCode() string
	ValidateEmail(email string) error
}

type MyEmailService struct {
}

func CreateEmailService() *MyEmailService {
	return &MyEmailService{}
}

func (s *MyEmailService) SendCodeToEmail(email, code string) error {
	//при добавлении редиса добавить ограничение по кол-ву писем
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

func (s *MyEmailService) CreateCode() string {
	codeInt := rand.IntN(maxValueCode-minValueCode) + minValueCode
	return strconv.Itoa(codeInt)
}

func (s *MyEmailService) ValidateEmail(email string) error {
	if strings.Contains(email, "@") && len(email) > 3 {
		return nil
	}
	return errors.New("unacceptable email")
}
