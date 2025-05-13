package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	client := newRedisClient()
	key := r.URL.Query().Get("key")
	field := r.URL.Query().Get("field")
	value := r.URL.Query().Get("value")
	if key == "" || field == "" || value == "" {
		http.Error(w, "Key, field, and value are required", http.StatusBadRequest)
		return
	}

	fmt.Println("Setting hash field:", field, "with value:", value)
	err := client.HSet(ctx, key, field, value).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error setting hash field:", err)
		return
	}
	fmt.Fprintf(w, "Hash field set successfully")
	fmt.Println("Hash field set successfully")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	client := newRedisClient()
	key := r.URL.Query().Get("key")
	field := r.URL.Query().Get("field")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	if field == "" {
		fmt.Println("Getting String key:", key)
		val, err := client.Get(ctx, key).Result()
		if err == redis.Nil {
			http.Error(w, "Key not found", http.StatusNotFound)
			fmt.Println("Key not found:", key)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error getting String key:", err)
			return
		}
		fmt.Fprintf(w, "Value: %s", val)
		fmt.Println("Value:", val)
	} else {
		fmt.Println("Getting hash field:", field)
		val, err := client.HGet(ctx, key, field).Result()
		if err == redis.Nil {
			http.Error(w, "Field not found", http.StatusNotFound)
			fmt.Println("Field not found:", field)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("Error getting hash field:", err)
			return
		}
		fmt.Fprintf(w, "Field: %s, Value: %s", field, val)
		fmt.Println("Field:", field, "Value:", val)
	}

}

func main() {
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}
