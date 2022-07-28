package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/faiz-tessact/go-postgres/pkg/websocket"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Notification struct {
	NotificationId      string `json:"notificationid"`
	NotificationTitle   string `json:"title"`
	UserId              string `json:"userid"`
	NotificationContent string `json:"content"`
	NotificationModel   string `json:"model"`
}

type JsonResponse struct {
	Type    string         `json:"type"`
	Data    []Notification `json:"data"`
	Message string         `json:"message"`
}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

func main() {
	printMessage("Getting Notifications...")
	handleRequests()
}

func handleRequests() {
	pool := websocket.NewPool()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/notifications/", GetNofications).Methods("GET")
	// router.HandleFunc("/notifications/", CreateNotifications).Methods("POST")
	// router.HandleFunc("/notifications/{notificationid}", DeleteNotifications).Methods("DELETE")
	go pool.Start()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Println("WebSocket error")
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

func GetNofications(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting notifications...")

	rows, err := db.Query("SELECT * FROM notifications")

	checkErr(err)

	var notifications []Notification

	for rows.Next() {
		var id int
		var notificationId string
		var notificationTitle string
		var userId string
		var notificationContent string
		var notificationModel string

		err = rows.Scan(&id, &notificationId, &notificationTitle, &userId, &notificationContent, &notificationModel)

		checkErr(err)

		notifications = append(notifications, Notification{NotificationId: notificationId, NotificationTitle: notificationTitle, UserId: userId, NotificationContent: notificationContent, NotificationModel: notificationModel})
	}

	var response = JsonResponse{Type: "success", Data: notifications}

	json.NewEncoder(w).Encode(response)
}

// func CreateMovie(w http.ResponseWriter, r *http.Request) {
// 	movieID := r.FormValue("movieid")
// 	movieName := r.FormValue("moviename")

// 	var response = JsonResponse{}

// 	if movieID == "" || movieName == "" {
// 		response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
// 	} else {
// 		db := setupDB()

// 		printMessage("Inserting movie into DB")

// 		fmt.Println("Inserting new movie with ID: " + movieID + " and name: " + movieName)

// 		var lastInsertID int
// 		err := db.QueryRow("INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", movieID, movieName).Scan(&lastInsertID)
// 		checkErr(err)

// 		response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
// 	}

// 	json.NewEncoder(w).Encode(response)
// }

// func DeleteMovie(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)

// 	movieID := params["movieid"]

// 	var response = JsonResponse{}

// 	if movieID == "" {
// 		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
// 	} else {
// 		db := setupDB()

// 		printMessage("Deleting movie from DB")

// 		_, err := db.Exec("DELETE FROM movies where movieID = $1", movieID)

// 		// check errors
// 		checkErr(err)

// 		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
// 	}

// 	json.NewEncoder(w).Encode(response)
// }
