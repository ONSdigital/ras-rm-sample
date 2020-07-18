package main

import (
	"net/http"

	"github.com/ONSdigital/ras-rm-sample/file-uploader/file"
	"github.com/ONSdigital/ras-rm-sample/file-uploader/health"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/samples/{type}/fileupload", file.ProcessFile)
	router.HandleFunc("/readiness", health.Readiness)
	router.HandleFunc("/liveness", health.Liveness)
	http.ListenAndServe(":8080", router)
}
