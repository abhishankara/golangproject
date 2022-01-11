package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/models"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connection = "mongodb+srv://abbhishekks004:1234mongo@cluster0.mrv2c.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
const dbName = "userinfo"
const colName = "userdata"

var collection *mongo.Collection
var db *sql.DB

func GetDB() *sql.DB {
	return db
}
func init() {
	var wg sync.WaitGroup
	wg.Add(2)
	go ConnectMongoDB(&wg)
	go ConnectMysql(&wg)
	wg.Wait()
	db = GetDB()
}

func ConnectMongoDB(wg *sync.WaitGroup) {
	defer wg.Done()
	clientOption := options.Client().ApplyURI(connection)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Connected to MongoDB")
	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection is ready")
}

func ConnectMysql(wg *sync.WaitGroup) {
	defer wg.Done()
	d, err := sql.Open("mysql", "root:1234sql@/usersql")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Connected to MySQL")
	db = d
	fmt.Println("MySQL is ready")

}

//Posting user information to MYSQL
func insertUserInfo(userInfo models.UserInfo) {
	query := "INSERT INTO `userinfo` (name, age, address, loginid, password) VALUES (?,?,?,?,?)"
	insert, err := db.Prepare(query)
	if err != nil {
		log.Panic(err)
	}
	defer insert.Close()
	_, err = insert.Exec(userInfo.Name, userInfo.Age, userInfo.Address, userInfo.Loginid, userInfo.Password)
	if err != nil {
		log.Panic(err)
	}
}

//Posting job information to MongoDB
func insertUserDetails(jobdata models.Jobdata) {
	inserted, err := collection.InsertOne(context.Background(), jobdata)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Inserted user data", inserted)
}

//Posting User information and job information to DB
func CreateUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	var user models.UserInfo
	_ = json.NewDecoder(r.Body).Decode(&user)
	insertUserInfo(user)
	insertUserDetails(user.Jobdetails)
	json.NewEncoder(w).Encode(user)
}

//Getting data from MYSQL
func getUserChan(ch chan models.UserInfo, name string) {
	var user models.UserInfo
	query := "SELECT * FROM userinfo WHERE name=?"
	row := db.QueryRow(query, name)
	err := row.Scan(&user.Name, &user.Age, &user.Address, &user.Loginid, &user.Password)
	if err != nil {
		log.Panic(err)
	}
	ch <- user
}

//Getting data from MYSQL
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	getChan := make(chan models.UserInfo)
	w.Header().Set("Content-type", "application/json")
	go getUserChan(getChan, params["name"])
	json.NewEncoder(w).Encode(<-getChan)
}

func alluserDetails(ch chan []primitive.M) {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Panic(err)
	}
	var userdetails []primitive.M
	for cur.Next(context.Background()) {
		var data bson.M
		err := cur.Decode(&data)
		if err != nil {
			log.Fatal(err)
		}
		userdetails = append(userdetails, data)
	}
	defer cur.Close(context.Background())
	ch <- userdetails
}

func GetJobInfo(w http.ResponseWriter, r *http.Request) {
	getChan := make(chan []primitive.M)
	w.Header().Set("Content-type", "application/json")
	go alluserDetails(getChan)
	json.NewEncoder(w).Encode(<-getChan)
}
