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
    err := client.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Key set successfully")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
    client := newRedisClient()
    val, err := client.Get(ctx, "key").Result()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Key: %s", val)
}

func main() {
    http.HandleFunc("/set", setHandler)
    http.HandleFunc("/get", getHandler)
    http.ListenAndServe(":8080", nil)
}
