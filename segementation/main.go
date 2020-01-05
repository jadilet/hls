package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/fsnotify/fsnotify"
	"github.com/jadilet/hls/segementation/cmd"
	"github.com/jadilet/hls/segementation/segment"
)

// CheckSegment walks video folder for segmenting
func CheckSegment() {
	conn := segment.Pool.Get()
	defer conn.Close()

	var videos []segment.Video

	err := filepath.Walk(segment.VideoPath, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)

		if segment.IsSupportedExtension(ext) {
			video := segment.Video{}
			video.Ext = ext
			video.Set(path)

			exists, err := redis.Bool(conn.Do("EXISTS", strings.ToUpper(video.Key)))

			if err != nil {
				log.Println(err)
			}

			if !exists {
				videos = append(videos, video)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, video := range videos {
		os.Mkdir(video.Dir, os.ModePerm)
		go video.Segment()
	}
}

// MonitorVideo monitors the video folder
func MonitorVideo() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				ext := filepath.Ext(event.Name)

				// start segmenting the video if a new video file found
				if (event.Op == fsnotify.Create) && segment.IsSupportedExtension(ext) {
					video := segment.Video{}
					video.Ext = ext
					video.Set(event.Name)
					// add a new folder with name "VIDEO_FILE_NAME + _ + FILE_EXTENSION" for the video segmenting
					os.Mkdir(video.Dir, os.ModePerm)

					go video.Segment()
				} else if (event.Op == fsnotify.Remove) && segment.IsSupportedExtension(ext) { // remove the segemnet files if a segmented file removed
					video := segment.Video{}
					video.Ext = ext
					video.Set(event.Name)

					go video.RemoveAll()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(segment.VideoPath)
	if err != nil {
		log.Println(err)
	}
	<-done
}

// main
func main() {
	segment.Pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	// check ffmpeg command exists
	ffmpegPath, err := exec.LookPath("ffmpeg")

	if err != nil {
		log.Println("ffmpeg command not found")
		return
	}
	cmd.FfmpegPath = ffmpegPath

	log.Println("ffmpeg command found in ", cmd.FfmpegPath)

	CheckSegment()
	MonitorVideo()
}
