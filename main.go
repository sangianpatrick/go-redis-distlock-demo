package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"github.com/go-redis/redis/v8"

	"github.com/gorilla/mux"
)

const (
	acquiredLockMessage    = "Acquired the lock."
	notAcquiredLockMessage = "Resource is locked. Acquired by another request."
)

var rs *redsync.Redsync

type result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {
	port := os.Getenv("PORT")

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	pool := goredis.NewPool(redisClient)

	rs = redsync.New(pool)

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

	mx := rs.NewMutex("lock:global",
		redsync.WithExpiry(time.Second*30),
		redsync.WithTries(5),
		redsync.WithRetryDelay(time.Second*1),
	)

	if err := mx.LockContext(ctx); err != nil {
		res := result{
			Success: false,
			Message: notAcquiredLockMessage,
		}
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(res)
		return
	}
	defer mx.UnlockContext(ctx)

	time.Sleep(time.Millisecond * 500)
	res := result{
		Success: true,
		Message: acquiredLockMessage,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func handlerWithHighLatency(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mx := rs.NewMutex("lock:global",
		redsync.WithExpiry(time.Second*30),
		redsync.WithTries(5),
		redsync.WithRetryDelay(time.Second*1),
	)

	if err := mx.LockContext(ctx); err != nil {
		res := result{
			Success: false,
			Message: notAcquiredLockMessage,
		}
		w.WriteHeader(http.StatusLocked)
		json.NewEncoder(w).Encode(res)
		return
	}
	defer mx.UnlockContext(ctx)

	time.Sleep(time.Second * 20)
	res := result{
		Success: true,
		Message: acquiredLockMessage,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}
