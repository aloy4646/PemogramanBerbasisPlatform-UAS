package controller

import (
	"database/sql"
	"kuis1/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * from wallet"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 150, "Query error")
		return
	}

	var detailTransaction model.DetailTransaction
	var detailTransactions []model.DetailTransaction
	for rows.Next() {
		if err := rows.Scan(&detailTransaction.Id, &detailTransaction.Currency, &detailTransaction.Username, &detailTransaction.Password, &detailTransaction.DisableUser); err != nil {
			log.Println(err.Error())
			errorResponseMessage(w, 170, "Data error")
			return
		} else {
			detailTransaction.Transactions, err = getAllTransactionsFromWallet(detailTransaction.Id, db)
			detailTransactions = append(detailTransactions, detailTransaction)
			if err != nil {
				log.Println(err)
				errorResponseMessage(w, 150, "Query error")
				return
			}
		}
	}

	sendSuccessResponseWithData(w, detailTransactions)
}

func getAllTransactionsFromWallet(idWallet int, db *sql.DB) ([]model.Transaction, error) {
	query := "SELECT * from transaction WHERE idWallet = " + strconv.Itoa(idWallet)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var transaction model.Transaction
	var transactions []model.Transaction
	var datetimeuint []uint8
	var amountuint []uint8
	for rows.Next() {
		if err := rows.Scan(&transaction.Id, &transaction.IdWallet, &datetimeuint, &amountuint, &transaction.Description); err != nil {
			return nil, err
		} else {
			transaction.Datetime = string(datetimeuint)
			transaction.Amount = string(amountuint)
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

func InsertNewTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	idWallet := r.Form.Get("idWallet")
	datetime := r.Form.Get("datetime")
	amount := r.Form.Get("amount")
	description := r.Form.Get("description")

	resultQuery, errQuery := db.Exec("INSERT INTO transaction (idWallet, datetime, amount, description) VALUES (?,?,?,?)",
		idWallet,
		datetime,
		amount,
		description,
	)

	if errQuery != nil {
		log.Println(errQuery)
		errorResponseMessage(w, 400, "Query error, Insert failed")
		return
	}

	id, _ := resultQuery.LastInsertId()

	sendSuccessResponseWithData(w, model.Transaction{Id: int(id), IdWallet: idWallet, Datetime: datetime, Amount: amount, Description: description})
}

func UpdateWallet(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		errorResponseMessage(w, 100, "Parse error")
		return
	}

	id := r.Form.Get("id")
	currency := r.Form.Get("currency")
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	resultQuery, errQuery := db.Exec("UPDATE wallet SET currency=?, username=?, password=? WHERE id=?",
		currency,
		username,
		password,
		id,
	)

	if errQuery != nil {
		log.Println(errQuery)
		errorResponseMessage(w, 400, "Query error, Update failed")
		return
	}

	rowsAffected, _ := resultQuery.RowsAffected()
	responseFromRowsAffected(w, rowsAffected)
}

func DeleteWallet(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	resultQuery, errQuery := db.Exec("UPDATE wallet SET disableUser=1 WHERE id=?",
		id,
	)

	if errQuery != nil {
		log.Println(errQuery)
		errorResponseMessage(w, 400, "Query error, Update failed")
		return
	}

	rowsAffected, _ := resultQuery.RowsAffected()
	responseFromRowsAffected(w, rowsAffected)
}
