package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sangianpatrick/go-redis-distlock-demo/sync"

	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"github.com/go-redis/redis/v8"

	"github.com/gorilla/mux"
)

const (
	acquiredLockMessage    = "Acquired the lock."
	notAcquiredLockMessage = "Resource is locked. Acquired by another request."
	statusOK               = "OK"
	statusLocked           = "LOCKED"
)

var distlock sync.DistributedLock

type result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	port := os.Getenv("PORT")
	redisHost := os.Getenv("REDIS_HOST")

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})

	pool := goredis.NewPool(redisClient)

	distlock = sync.NewRedsyncAdapter(pool)

	router := mux.NewRouter()

	router.HandleFunc("/healthcheck", handlerHealthCheck).Methods(http.MethodGet)
	router.HandleFunc("/low", handlerWithLowLatency).Methods(http.MethodGet)
	router.HandleFunc("/high", handlerWithHighLatency).Methods(http.MethodGet)

	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Service is running properly.")
}

func handlerWithLowLatency(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mutex := distlock.NewMutex("lock:global", 5, time.Second*1, time.Second*30)

	w.Header().Add("Content-Type", "application/json")

	if err := mutex.Lock(ctx); err != nil {
		res := result{
			Success: false,
			Message: notAcquiredLockMessage,
			Status:  statusLocked,
		}
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(res)
		return
	}
	defer mutex.Unlock(ctx)

	time.Sleep(time.Millisecond * 500)
	res := result{
		Success: true,
		Message: acquiredLockMessage,
		Status:  statusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func handlerWithHighLatency(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mutex := distlock.NewMutex("lock:global", 5, time.Second*1, time.Second*30)

	w.Header().Add("Content-Type", "application/json")

	if err := mutex.Lock(ctx); err != nil {
		res := result{
			Success: false,
			Message: notAcquiredLockMessage,
			Status:  statusLocked,
		}
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(res)
		return
	}
	defer mutex.Unlock(ctx)

	time.Sleep(time.Second * 20)
	res := result{
		Success: true,
		Message: acquiredLockMessage,
		Status:  statusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}
