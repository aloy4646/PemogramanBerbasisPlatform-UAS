package controller

import (
	"kuis1/model"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()
var jwtKey = []byte(os.Getenv("JWT_TOKEN"))
var tokenName = os.Getenv("TOKEN_NAME")

type Claims struct {
	ID          int    `json:"Id"`
	Username    string `json:"Name"`
	Email       string `json:"Email"`
	Password    string `json:"Password"`
	Type        string `json:"Type"`
	DisableUser bool   `json:"DisableUser"`
	Activated   bool   `json:"Activated"`
	jwt.StandardClaims
}

func generateToken(w http.ResponseWriter, user model.Pengguna) {
	tokenExpiryTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		ID:          user.Id,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		Type:        user.Type,
		DisableUser: user.DisableUser,
		Activated:   user.Activated,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    signedToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})
}

func resetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
		Path:     "",
	})
}

func Authenticate(next http.HandlerFunc, accessType string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := validateUserToken(r, accessType)
		if !isValidToken {
			sendUnAuthorizedResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(r *http.Request, accessType string) bool {
	isAccessTokenValid, userType := validateTokenFromCookies(r)
	// fmt.Print(userType, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func validateTokenFromCookies(r *http.Request) (bool, string) {
	if cookie, err := r.Cookie(tokenName); err == nil {
		accessToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.Type
		}
	}
	return false, "GUEST"
}

func getUserFromCookies(r *http.Request) model.Pengguna {
	cookie, err := r.Cookie(tokenName)
	var user model.Pengguna
	if err == nil {
		accessToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			user = model.Pengguna{Id: accessClaims.ID, Username: accessClaims.Username, Email: accessClaims.Email, Password: accessClaims.Password,
				Type: accessClaims.Type, DisableUser: accessClaims.DisableUser, Activated: accessClaims.Activated}
			return user
		}
	}
	return user
}
