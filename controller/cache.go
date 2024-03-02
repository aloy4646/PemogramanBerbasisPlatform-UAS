package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"kuis1/model"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func GetUsersFromCache(r *http.Request) model.OTPModel {
	var userOTP model.OTPModel

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var ctx = context.Background()
	var user model.Pengguna
	user = getUserFromCookies(r)

	key := "userOTP" + strconv.Itoa(user.Id)
	fmt.Println(key)
	fmt.Println(key)

	value, err := client.Get(ctx, key).Result()
	if err != nil {
		log.Println("Get Error")
		log.Println(err)
		return userOTP
	}

	_ = json.Unmarshal([]byte(value), &userOTP)

	return userOTP
}

func SetUsersToCache(userOTP model.OTPModel) {
	converted, err := json.Marshal(userOTP)
	if err != nil {
		log.Println(err)
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var ctx = context.Background()

	key := "userOTP" + strconv.Itoa(userOTP.Pengguna.Id)
	fmt.Println(key)
	fmt.Println(key)

	err = client.Set(ctx, key, converted, 0).Err()
	if err != nil {
		log.Println("Set Error")
		log.Println(err)
		return
	} else {
		log.Println("Cache set")
	}
}
