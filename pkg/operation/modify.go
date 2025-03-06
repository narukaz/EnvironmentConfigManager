package operation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	Client *mongo.Client
}

type Employee struct {
	ID     string `bson:"_id" json:"id"`
	Name   string `json:"name" bson:"name"`
	Phone  int64  `json:"phone" bson:"phone"`
	Gender string `json:"gender" bson:"gender"`
}

func (h *Handler) GetAllEmployee(w http.ResponseWriter, r *http.Request) {
	var data []map[string]interface{}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	collection := h.Client.Database("test").Collection("employee")
	res, _ := collection.Find(ctx, bson.M{})

	for res.Next(ctx) {
		var doc map[string]interface{}
		if err := res.Decode(&doc); err != nil {
			log.Fatal(err)
		}
		data = append(data, doc)
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "failed to encode data", http.StatusBadRequest)
		log.Fatal("failed to encode data ")
		return
	}

}

func (h *Handler) InsertItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method is not allowed", http.StatusBadRequest)
		return
	}
	var data Employee

	err := json.NewDecoder(r.Body).Decode(&data)
	data.ID = primitive.NewObjectID().Hex()
	if err != nil {
		fmt.Println("error while parsing data: ", err)
		http.Error(w, "error while parsing data", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	collection := h.Client.Database("test").Collection("employee")
	res, _ := collection.InsertOne(ctx, data)
	if !res.Acknowledged {
		http.Error(w, "error Inserting in Database", http.StatusBadRequest)
		return
	}
	w.Write([]byte("successfully created"))

}

func (h *Handler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var req struct {
		ID string `bson:"_id" json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("decode request body:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	// fmt.Println(conv)

	collection := h.Client.Database("test").Collection("employee")
	res, err := collection.DeleteOne(ctx, bson.M{"_id": req.ID})
	if err != nil {
		fmt.Println("unable to write to database")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if res.DeletedCount == 0 {
		fmt.Println("unable to write to database")
		http.Error(w, "unable to write to database", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int64{"delete count": res.DeletedCount})
}
