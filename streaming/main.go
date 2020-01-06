package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

// VideoPath video path
const VideoPath = "/var/videos"

func hlsM3U8Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	video, err := findVideo(vars["videoName"])

	if err == errNoVideo {
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if video.indexFileExist() {

		http.ServeFile(w, r, video.indexFilePath())
		w.Header().Set("Content-Type", "application/x-mpegURL")
	} else {
		http.NotFound(w, r)
	}
}

func hlsTsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	video, err := findVideo(vars["videoName"])

	if err == errNoVideo {
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	segName := vars["segName"]

	if video.tsFileExist(segName) {
		http.ServeFile(w, r, video.tsFilePath(segName))
		w.Header().Set("Content-Type", "video/MP2T")
	} else {
		http.NotFound(w, r)
	}
}

func handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{videoName}/index.m3u8", hlsM3U8Handler).Methods("GET")
	router.HandleFunc("/{videoName}/{segName:index[0-9]+.ts}", hlsTsHandler).Methods("GET")

	return router
}

func main() {
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	http.Handle("/", handlers())
	http.ListenAndServe(":8080", nil)
}
