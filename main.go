package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kabukky/httpscerts"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	Routes()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type UserInformation struct {
	ID         int
	First_Name string
	Last_Name  string
	Email      string
	Gender     string
	Age        int
}

type ResponseStatus struct {
	Status  string
	Message string
}

func Routes() {
	r := mux.NewRouter()
	r.HandleFunc("/getUser", getUserInformation)
	r.HandleFunc("/createUser", createUser)
	r.HandleFunc("/updateUser", updateUser)
	r.HandleFunc("/deleteUser", deleteUser)
	r.HandleFunc("/getAgeRange", getAgeRange)
	r.HandleFunc("/getEditDataUser", getEditDataUser)
	http.Handle("/", r)
	err := httpscerts.Check("cert.pem", "key.pem")
	// If they are not available, generate new ones.http.ListenAndServe(":8080", nil)
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", "localhost")
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}
	http.ListenAndServe(":8080", nil)
}

func getUserInformation(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "../user.db")
	checkErr(err)
	var userInfo []UserInformation

	// select data from users table
	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)
	var ID int
	var First_Name string
	var Last_Name string
	var Email string
	var Gender string
	var Age int
	for rows.Next() {
		err = rows.Scan(&ID, &First_Name, &Last_Name, &Email, &Gender, &Age)
		checkErr(err)
		userInfo = append(userInfo, UserInformation{
			ID:         ID,
			First_Name: First_Name,
			Last_Name:  Last_Name,
			Email:      Email,
			Gender:     Gender,
			Age:        Age,
		})
	}
	rows.Close() //good habit to close
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(userInfo)
}

func createUser(w http.ResponseWriter, r *http.Request) {

	var ID int
	var First_Name string
	var Last_Name string
	var Email string
	var Gender string
	var Age string

	db, err := sql.Open("sqlite3", "../user.db")
	checkErr(err)
	rows, err := db.Query("SELECT * FROM users ORDER BY id DESC")
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&ID, &First_Name, &Last_Name, &Email, &Gender, &Age)
		checkErr(err)
		break
	}

	rows.Close() //good habit to close
	// insert
	stmt, err := db.Prepare("INSERT INTO users(id,first_name, last_name, email,gender,age) values(?,?,?,?,?,?)")
	checkErr(err)

	ID++
	First_Name = r.FormValue("First_Name")
	Last_Name = r.FormValue("Last_Name")
	Email = r.FormValue("Email")
	Gender = r.FormValue("Gender")
	Age = r.FormValue("Age")
	res, err := stmt.Exec(ID, First_Name, Last_Name, Email, Gender, Age)

	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err == nil || res != nil {
		response := ResponseStatus{
			Status:  "success",
			Message: "เพิ่มข้อมูลสำเร็จ",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		response := ResponseStatus{
			Status:  "error",
			Message: "เพิ่มข้อมูลไม่สำเร็จ",
		}
		json.NewEncoder(w).Encode(response)
	}

}
func updateUser(w http.ResponseWriter, r *http.Request) {

	var ID string
	var First_Name string
	var Last_Name string
	var Email string
	var Gender string
	var Age string

	ID = r.FormValue("ID")
	First_Name = r.FormValue("First_Name")
	Last_Name = r.FormValue("Last_Name")
	Email = r.FormValue("Email")
	Gender = r.FormValue("Gender")
	Age = r.FormValue("Age")

	db, err := sql.Open("sqlite3", "../user.db")
	checkErr(err)

	// update
	stmt, err := db.Prepare("UPDATE users SET first_name=?, last_name=?, email=?,gender=?,age=? WHERE id=?")
	checkErr(err)

	res, err := stmt.Exec(First_Name, Last_Name, Email, Gender, Age, ID)
	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err == nil || res != nil {
		response := ResponseStatus{
			Status:  "success",
			Message: "แก้ไขข้อมูลสำเร็จ",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		response := ResponseStatus{
			Status:  "error",
			Message: "แก้ไขข้อมูลไม่สำเร็จ",
		}
		json.NewEncoder(w).Encode(response)
	}

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var ID string
	ID = r.FormValue("ID")

	db, err := sql.Open("sqlite3", "../user.db")
	checkErr(err)
	// delete
	stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
	checkErr(err)

	res, err := stmt.Exec(ID)
	checkErr(err)
	db.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err == nil || res != nil {
		response := ResponseStatus{
			Status:  "success",
			Message: "ลบข้อมูลสำเร็จ",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		response := ResponseStatus{
			Status:  "error",
			Message: "ลบข้อมูลไม่สำเร็จ",
		}
		json.NewEncoder(w).Encode(response)
	}

}
func getAgeRange(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "../user.db")
	checkErr(err)
	var userInfo []UserInformation

	var ID int
	var First_Name string
	var Last_Name string
	var Email string
	var Gender string
	var Age int

	var Age_Start = r.FormValue("Age_Start")
	var Age_End = r.FormValue("Age_End")
	// select data where 1 - 2

	stmt, err := db.Prepare("SELECT * FROM users WHERE age BETWEEN ? AND ?")
	checkErr(err)
	rows, err := stmt.Query(Age_Start, Age_End)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&ID, &First_Name, &Last_Name, &Email, &Gender, &Age)
		checkErr(err)
		userInfo = append(userInfo, UserInformation{
			ID:         ID,
			First_Name: First_Name,
			Last_Name:  Last_Name,
			Email:      Email,
			Gender:     Gender,
			Age:        Age,
		})
	}
	rows.Close() //good habit to close
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(userInfo)
}

func getEditDataUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "../user.db")
	checkErr(err)
	var userInfo []UserInformation

	var ID int
	var First_Name string
	var Last_Name string
	var Email string
	var Gender string
	var Age int

	var UserID = r.FormValue("ID")

	stmt, err := db.Prepare("SELECT * FROM users WHERE id=?")
	checkErr(err)
	rows, err := stmt.Query(UserID)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&ID, &First_Name, &Last_Name, &Email, &Gender, &Age)
		checkErr(err)
		userInfo = append(userInfo, UserInformation{
			ID:         ID,
			First_Name: First_Name,
			Last_Name:  Last_Name,
			Email:      Email,
			Gender:     Gender,
			Age:        Age,
		})
	}
	rows.Close() //good habit to close
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(userInfo)
}
