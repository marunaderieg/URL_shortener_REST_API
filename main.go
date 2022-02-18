package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thanhpk/randstr"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
)

type Request struct {
	Url string `json:"url"`
}

var database *sql.DB

/*Saves the JSON input of form "url":"example.com/mylongurl" into a sqlite database. Each url is assigned a random
string of 8 characters which serves as a short url. The short url is also used as the primary key in the
database. One long url can have multiple short urls redirecting to it. Each short url is unique.*/
func shorten(w http.ResponseWriter, r *http.Request) {
	// check if url is valid
	var request Request
	_ = json.NewDecoder(r.Body).Decode(&request)
	_, err := url.ParseRequestURI(request.Url)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		response := `{"errormessage": "` + request.Url + ` is not a valid url. Maybe you forgot to add the protocol?"}`
		w.WriteHeader(400)
		w.Write([]byte(response))
		return
	}
	/* Generate a random string of 8 characters, which will serve as the shortened url. This random string is also the
	primary key in the database. Try to add the new "short_url - url" pair to the database. This will return an error
	if the generated short url already exists in the database. Repeat until the database insert has been successful
	or until 62**8 iterations have been performed. */
	i := math.Pow(62, 8)
	for i > 0 {
		random := randstr.String(8)
		//insert row
		statement, _ := database.Prepare("INSERT INTO urls (short_url, url) VALUES (?,?)")
		_, err := statement.Exec(random, request.Url)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			response := `{"short_url": "http://127.0.0.1:8080/` + random + `"}`
			w.WriteHeader(200)
			w.Write([]byte(response))
			return
		} else {
			i -= 1
		}
	}
	/* Error handling in case a lot of entries are saved to the DB. Consider doing autoincrement of the primary key
	or create shortened urls with more than 8 characters. */
	w.Header().Set("Content-Type", "application/json")
	response := `{"errormessage": "too many entries are stored in the database"}`
	w.WriteHeader(500)
	w.Write([]byte(response))
	log.Print("Too many entries are stored in the database.")
}

/* Performs a 302 redirect from "http://127.0.0.1:8080/short_url" to the corresponding long url.
An error is returned in case the provided short url is not present in the database. */
func redirect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	rows, _ := database.Query("SELECT url FROM urls WHERE short_url=?", params["id"])
	var myUrl string
	for rows.Next() {
		rows.Scan(&myUrl)
	}
	if myUrl != "" {
		http.Redirect(w, r, myUrl, 302)
	} else {
		w.Header().Set("Content-Type", "application/json")
		response := `{"errormessage": "passed short_url does not exist in database"}`
		w.WriteHeader(400)
		w.Write([]byte(response))
	}
}

func main() {
	// Create Database
	os.Remove("sqlite-database.db")
	database, _ = sql.Open("sqlite3", "./sqlite-database.db")
	statement, _ :=
		database.Prepare("CREATE TABLE urls (short_url TEXT PRIMARY KEY, url TEXT)")
	statement.Exec()
	// Init Router
	r := mux.NewRouter()
	// Route Handlers / Endpoints
	r.HandleFunc("/{id}", redirect).Methods("GET")
	r.HandleFunc("/shorten", shorten).Methods("POST")
	// Error Logging
	log.Fatal(http.ListenAndServe(":8080", r))
}
