package controllers

import (
	"context"
	"encoding/json"
	"flipr_assignment/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var noteCollection *mongo.Collection

func init() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connectionUri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	colName := os.Getenv("NOTES_COLLECTION")
	clientOption := options.Client().ApplyURI(connectionUri)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Printf("Failed to ping MongoDB: %v\n", err)
		return
	}
	
	fmt.Println("Connected to MongoDB")
	noteCollection = client.Database(dbName).Collection(colName)
	fmt.Println("Notes collection ready")
}


func GetNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	  options := &options.FindOptions{}
	  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	  defer cancel()
	  cur, err := noteCollection.Find(ctx, bson.D{{}}, options)
	  if err != nil {
		  log.Fatal(err)
	  }
	  var notes []primitive.M
	  for cur.Next(context.Background()) {
		  var note bson.M
		  err := cur.Decode(&note)
		  if err != nil {
			  log.Fatal(err)
		  }
		  notes = append(notes, note)
	  }
	  json.NewEncoder(w).Encode(notes)
  }

  func GetNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var note bson.M
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := noteCollection.FindOne(ctx, filter).Decode(&note)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(w).Encode("Song doesn't exist")
			return
		}
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(note)
}

func AddNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	var foundNote models.Note
	_ = json.NewDecoder(r.Body).Decode(&note)
	filter := bson.M{"title": note.Title, "description": note.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := noteCollection.FindOne(ctx, filter).Decode(&foundNote)
	if err == nil {
		json.NewEncoder(w).Encode("Note already exists")
		return
	}
	inserted, err := noteCollection.InsertOne(context.Background(), note)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(inserted)
}

func UpdateNote(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	var note models.Note
	_ = json.NewDecoder(r.Body).Decode(&note)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}

	// Create an update document with the $set operator and the new values
	update := bson.M{"$set": bson.M{}}

    if note.Title != "" {
        update["$set"].(bson.M)["title"] = note.Title
    }

    if note.Description != "" {
        update["$set"].(bson.M)["description"] = note.Description
    }
	options := options.Update().SetUpsert(false)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := noteCollection.UpdateOne(ctx, filter, update, options)
	if(err != nil){
		json.NewEncoder(w).Encode("the specified id doesn't exists")
	}
	json.NewEncoder(w).Encode(result)
}

func DeleteNote(w http.ResponseWriter, r *http.Request){

	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	deleteResult, err := noteCollection.DeleteOne(ctx, bson.M{"_id": id})
	if deleteResult.DeletedCount == 0 {
    log.Fatal("Error on deleting", err)
}
json.NewEncoder(w).Encode("note deleted")
}

