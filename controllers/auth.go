package controllers

import "net/http"

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

func register() {
	//i
}

func login(w http.ResponseWriter, r *http.Request) {

}
