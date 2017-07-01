package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	//for extracting service credentials from VCAP_SERVICES
	//"github.com/cloudfoundry-community/go-cfenv"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
)

const (
	DEFAULT_PORT = "8080"
)

var index = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/index.html",
))

type StatusCodeDescription struct {
	Code        int
	Description string
}

type Response struct {
	Error    string
	Response StatusCodeDescription
}

func helloworld(w http.ResponseWriter, req *http.Request) {
	index.Execute(w, nil)
}

func writeJSON(w http.ResponseWriter, code int, obj interface{}) (err error) {
	var bytes []byte

	if bytes, err = json.Marshal(obj); err != nil {
		log.Println(err.Error())
	}

	w.WriteHeader(code)
	_, _ = w.Write(bytes)
	return
}

func describeHTTPStatusCode(w http.ResponseWriter, req *http.Request) {
	var response Response

	vars := mux.Vars(req)
	statuscode, err := strconv.Atoi(vars["statuscode"])

	if err != nil {
		response.Error = err.Error()
		log.Println(err.Error())
		_ = writeJSON(w, http.StatusBadRequest, response)
		return
	}

	description := http.StatusText(statuscode)

	scd := StatusCodeDescription{Code: statuscode, Description: description}
	response.Response = scd

	var bytes []byte

	if bytes, err = json.Marshal(response); err != nil {
		response.Error = err.Error()
		log.Println(err.Error())
		_ = writeJSON(w, http.StatusInternalServerError, response)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(bytes); err != nil {
		response.Error = err.Error()
		log.Println(err.Error())
		_ = writeJSON(w, http.StatusInternalServerError, response)
		return
	}

}

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}

	r := mux.NewRouter()
	r.HandleFunc("/", helloworld)
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/status/{statuscode}", describeHTTPStatusCode)
	http.Handle("/", r)

	log.Printf("Starting app on port %+v\n", port)
	http.ListenAndServe(":"+port, nil)
}
