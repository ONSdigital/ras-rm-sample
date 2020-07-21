package main

import (
	"net/http"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/routes"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/samples/{type}/fileupload", routes.ProcessFile)
	router.HandleFunc("/readiness", routes.Readiness)
	router.HandleFunc("/liveness", routes.Liveness)
	http.ListenAndServe(":8080", router)
}
