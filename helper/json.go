package helper

import (
	"encoding/json"
	"google-oauth/web"
	"net/http"
)

func WriteEncodeResponse(writer http.ResponseWriter, webResponse web.WebResponse) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(webResponse)
	if err != nil {
		panic(err)
	}
}