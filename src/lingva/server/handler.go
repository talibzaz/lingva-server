package server

import (
	"net/http"
	"encoding/json"
	db "lingva/database"
	"os"
	"io"
)

type User struct {
	Username string	`json:"username"`
	Password string	`json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request){
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := db.Login(user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	j, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(j)
	return
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := db.AdminLogin(user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	j, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(j)
	return
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	imageName := r.URL.Path[len("/image/"):]

	img, err := os.Open("files/"+imageName)
	if err != nil {
		//log.Info("Image not available")
		//log.Error(err)
		w.Write([]byte("Image not available"))
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(w, img)
}