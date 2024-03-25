package handlers

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Viet-ph/xss-vulnerable/database"
	"github.com/Viet-ph/xss-vulnerable/response"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	var (
		email    string
		password string
	)
	if r.Method == "POST" {
		email = r.FormValue("username")
		password = r.FormValue("password")
	} else {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	//Check user credential
	for _, user := range db.Users {
		if user.Email == email {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
				accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					Subject:   strconv.Itoa(user.Id),
				})

				signedAccessToken, _ := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

				http.SetCookie(w, &http.Cookie{
					Name:  "jwt_token",
					Value: signedAccessToken,
					Path:  "/",
				})
				http.Redirect(w, r, "/search", http.StatusSeeOther)

			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "jwt_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Redirect to the login page after logging out
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			// No session cookie; redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		reqToken := cookie.Value

		token, err := jwt.ParseWithClaims(reqToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			response.RespondWithError(w, http.StatusUnauthorized, "Unauthorized Access")
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ExtractIdFromToken(token *jwt.Token) (string, error) {
	id, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return id, nil
}
