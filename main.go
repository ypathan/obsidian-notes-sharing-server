package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/sessions"
)


func signToken(username string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(username))
	return hex.EncodeToString(mac.Sum(nil))
}

func checkSignToken(receivedToken string, storedToken string) bool {
	if receivedToken == storedToken {
		return true
	}
	return false
}

func main() {
	port := ":8080"
	dir := "/Users/myousuf/dev/obs-notes/obsdian-notes/dist/"

	if os.Getenv("env") == "prod"  {
		dir = "/data/"
	}

	var secret = os.Getenv("secret")
	var store = sessions.NewCookieStore([]byte(secret))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dir, filepath.Clean(r.URL.Path))

		// Get file info
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		if info.IsDir() {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, path)
	})

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/admin.html")
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == os.Getenv("username") && password == os.Getenv("password") {
			secret_key := os.Getenv("secret")
			token := signToken(username, secret_key)

			tokenCookie := &http.Cookie{
				Name:     "session_token",
				Value:    token,
				Path:     "/admin",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}

			usernameCookie := &http.Cookie{
				Name:     "username",
				Value:    username,
				Path:     "/admin",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}

			session, _ := store.Get(r, "user-session")
			session.Values[username] = token
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, tokenCookie)
			http.SetCookie(w, usernameCookie)
			http.Redirect(w, r, "/admin", 302)
			return
		}

	})

	fileDir := http.Dir(dir)
	fileServer := http.StripPrefix("/admin/", http.FileServer(fileDir))

	mux.Handle("/admin/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		received_token, err := r.Cookie("session_token")
		if err != nil {
			//TODO return 
			fmt.Println("Error getting token from cookie")
			return
		}
		username, err := r.Cookie("username")
		if err != nil {
			//TODO return
			fmt.Println("Error getting username from cookie")
			return
		}

		session, _ := store.Get(r, "user-session")
		token, ok := session.Values[username.Value].(string)
		if !ok {
			http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
			return
		}

		if checkSignToken(received_token.Value, token) {
			fileServer.ServeHTTP(w, r)
		}
	
		//TODO graceful return
		return
	}))

	println("Server running at http://localhost" + port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("Error starting server")
	}


}
