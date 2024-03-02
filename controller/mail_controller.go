package controller

import (
	"bytes"
	"fmt"
	"kuis1/model"
	"text/template"

	gomail "gopkg.in/gomail.v2"
)

func SendWelcomeEmail(user model.Pengguna) {
	mail := gomail.NewMessage()

	template := "controller/welcomeMessage.html"

	result, _ := parseTemplate(template, user)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Testing Send Email")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email Delivered to ", user.Email)
	}
}

func SendOTPEmail(userOTP model.OTPModel) {
	mail := gomail.NewMessage()

	template := "controller/sendOTP.html"

	result, _ := parseTemplate(template, userOTP)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", userOTP.Pengguna.Email)
	mail.SetHeader("Subject", "Testing Send Email")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email Delivered to ", userOTP.Pengguna.Email)
	}
}

func sendEmailWithMessage(message model.Message) {
	mail := gomail.NewMessage()

	template := "controller/sendMessage.html"

	result, _ := parseTemplate(template, message)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", message.Pengguna.Email)
	mail.SetHeader("Subject", "Testing Send Email")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email Delivered to ", message.Pengguna.Email)
	}
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	// mengubah text html ke dalam bentuk byte
	t, err := template.ParseFiles(templateFileName)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	buff := new(bytes.Buffer)

	//render struct ke dalam file html (td ada var name di dalam htmlnya)
	if err = t.Execute(buff, data); err != nil {
		fmt.Println(err)
		return "", err
	}

	return buff.String(), nil
}
