package main

import (
	"fmt"
	"kuis1/controller"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/transactions", controller.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions", controller.InsertNewTransaction).Methods("POST")
	router.HandleFunc("/wallet", controller.UpdateWallet).Methods("PUT")
	router.HandleFunc("/wallet/{id}", controller.DeleteWallet).Methods("DELETE")
	router.HandleFunc("/login", controller.Login).Methods("POST")
	router.HandleFunc("/logout", controller.Logout).Methods("POST")

	//UAS
	router.HandleFunc("/users", controller.RegisterUser).Methods("POST")
	router.HandleFunc("/users/generateotp", controller.Authenticate(controller.GenerateOTP, "USER")).Methods("GET")
	router.HandleFunc("/users/sendemail", controller.Authenticate(controller.SendEmail, "USER")).Methods("POST")
	
	//lupa nambahin sesuatu
	router.HandleFunc("/users/aktivasi", controller.Authenticate(controller.AktivasiAkun, "USER")).Methods("POST")

	//Tadi salah nama endpoint
	router.HandleFunc("/admin/resend", controller.Authenticate(controller.ResendOTP, "ADMIN")).Methods("POST")

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://ithb.ac.id"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	http.Handle("/", handler)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
