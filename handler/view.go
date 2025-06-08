package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func LoginView(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// fmt.Fprint(writer, "welcome ", name)
	http.ServeFile(writer, request, "./resources/login.html")
}

func RegisterView(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// fmt.Fprint(writer, "welcome ", name)
	http.ServeFile(writer, request, "./resources/register.html")
}