package controller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"kuis1/model"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	//wallet
	resultQuery, errQuery := db.Exec("INSERT INTO pengguna (username, email, password, type, disableUser, activated) VALUES (?,?,?,'USER', 0, 0)",
		username,
		email,
		password,
	)

	if errQuery != nil {
		log.Println(errQuery)
		errorResponseMessage(w, 400, "Query error, Insert failed")
		return
	}

	lastId, _ := resultQuery.LastInsertId()

	var user model.Pengguna = model.Pengguna{Id: int(lastId), Username: username, Email: email, Password: password, Type: "USER", DisableUser: false, Activated: false}
	go SendWelcomeEmail(user)
	generateToken(w, user)

	//response
	sendSuccessResponseWithData(w, user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	var pengguna model.Pengguna
	if err := db.QueryRow("SELECT * from pengguna where username=? AND password=?",
		username, password).Scan(&pengguna.Id, &pengguna.Username, &pengguna.Email, &pengguna.Password, &pengguna.Type, &pengguna.DisableUser, &pengguna.Activated); err != nil {
		log.Println(err.Error())
		errorResponseMessage(w, 170, "False")
		return
	}

	generateToken(w, pengguna)

	var response model.SuccessResponse
	response.Status = 200
	response.Message = "Login success"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	successResponseMessage(w)
}

func GenerateOTP(w http.ResponseWriter, r *http.Request) {
	otp := generateOTP()
	fmt.Println(otp)
	fmt.Println(otp)
	fmt.Println(otp)

	var user model.Pengguna
	user = getUserFromCookies(r)
	var userOTP model.OTPModel
	userOTP = model.OTPModel{Pengguna: user, OTP: otp}
	SetUsersToCache(userOTP)

	go SendOTPEmail(userOTP)
	successResponseMessage(w)
}

func generateOTP() string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	message := r.Form.Get("message")

	var user model.Pengguna
	user = getUserFromCookies(r)
	var userMessage model.Message
	userMessage = model.Message{Pengguna: user, Message: message}

	go sendEmailWithMessage(userMessage)
	successResponseMessage(w)
}

func AktivasiAkun(w http.ResponseWriter, r *http.Request) {
	var userOTP model.OTPModel
	userOTP = GetUsersFromCache(r)

	if userOTP.OTP == "" {
		errorResponseMessage(w, 404, "OTP not found")
		return
	}

	/////////////////////////////////////////////////////////////////
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}
	/////////////////////////////////////////////////////////////////

	otp := r.Form.Get("otp")

	fmt.Println(otp + " - " + userOTP.OTP)
	if otp != userOTP.OTP {
		errorResponseMessage(w, 406, "OTP is wrong")
		return
	}

	successResponseMessage(w)
}

func ResendOTP(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	email := r.Form.Get("email")

	row, err := db.Query("SELECT * from pengguna where email=?", email)
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 150, "Query error")
		return
	}

	var user model.Pengguna
	for row.Next() {
		if err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Type, &user.DisableUser, &user.Activated); err != nil {
			log.Println(err.Error())
			errorResponseMessage(w, 170, "Data error")
			return
		}
	}

	var userOTP model.OTPModel
	otp := generateOTP()
	userOTP = model.OTPModel{Pengguna: user, OTP: otp}
	SetUsersToCache(userOTP)

	go SendOTPEmail(userOTP)
	successResponseMessage(w)
}
