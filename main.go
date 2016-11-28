package main

import (
	"net/http"
	"encoding/json"
	"github.com/pidgeonhole/owl/lib"
	"log"
)

func handleJob(w http.ResponseWriter, r *http.Request) {
	var job owl.Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	defer r.Body.Close()

	results, err := owl.RunJob(job)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	json.NewEncoder(w).Encode(results)
	return
}

func main() {
	http.HandleFunc("/", handleJob)
	log.Fatal(http.ListenAndServe("0.0.0.0:3001", nil))
}
