package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       int    `json:"_id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Post struct {
	ID               int    `json:"_id" bson:"_id"`
	User_ID          int    `json:"userId" bson:"userId"`
	Caption          string `json:"caption" bson:"caption"`
	Image_URL        string `json:"url" bson:"url"`
	Posted_Timestamp string `json:"timestamp" bson:"timestamp"`
}

var client *mongo.Client

func CreateUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("instagram").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)

}

func GetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var userdetails []User
	collection := client.Database("instagram").Collection("userdetails")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message" :"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		userdetails = append(userdetails, user)

	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message" :"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(userdetails)

}

func CreatePost(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)
	collection := client.Database("instagram").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)

}

func main() {
	fmt.Println("Hello")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, "mongodb://localhost:5000")
	http.HandleFunc("/users", CreateUsers)
	http.HandleFunc("/post", CreatePost)
	http.HandleFunc("/users/", GetUser)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}

}
