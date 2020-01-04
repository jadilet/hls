package segment

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/jadilet/hls/segementation/cmd"
)

// Pool of redis connection
var Pool *redis.Pool

// VideoPath video path
const VideoPath = "/var/videos"

// SupportedExtensions supported video extensions
var SupportedExtensions = map[string]string{
	".mp4": "mp4",
	".mov": "mov",
}

// Video struct
type Video struct {
	Name string // file name
	Path string // file path
	Dir  string // file directory
	Ext  string // extension
	Key  string // key for redis
}

// Set the Video struct values by the video file path
func (video *Video) Set(videoPath string) {

	video.Name = filepath.Base(videoPath)
	video.Path = videoPath

	videoFolder := strings.TrimSuffix(video.Name, video.Ext) + "_" +
		SupportedExtensions[video.Ext]

	video.Key = videoFolder
	video.Dir = filepath.Join(VideoPath, videoFolder)
}

// Segment the video with ffmeg
func (video Video) Segment() {
	conn := Pool.Get()
	defer conn.Close()

	cmd := exec.Command(cmd.FfmpegPath, cmd.Params(video.Path)...)
	cmd.Dir = video.Dir
	err := cmd.Start()
	log.Printf("started segmenting the video %s", video.Path)

	if err != nil {
		log.Println(err)
	}

	err = cmd.Wait()

	if err != nil {
		os.RemoveAll(video.Dir)
		log.Println("finished segmenting with error ", err)
	} else {
		serialized, err := json.Marshal(&video)

		if err != nil {
			log.Println(err)
		}

		log.Println("finished segementing successfully ", string(serialized))

		_, err = conn.Do("SET", video.Key, serialized)

		if err != nil {
			log.Println(err)
		}
	}
}

// RemoveAll removes all segemented video
func (video Video) RemoveAll() {
	conn := Pool.Get()
	defer conn.Close()

	exists, err := redis.Int(conn.Do("EXISTS", video.Key))

	if err != nil {
		log.Println(err)
	}

	if exists == 1 { // found
		log.Println("started removing the segment file: ", video.Dir)
		_, err = conn.Do("DEL", video.Key)

		if err != nil {
			log.Println(err)
		}

		err = os.RemoveAll(video.Dir)

		if err != nil {
			log.Println(err)
		}
		log.Println("finished removing the segment file: ", video.Dir)
	}
}

// IsSupportedExtension check video extension is supported
func IsSupportedExtension(ext string) bool {
	if _, ok := SupportedExtensions[ext]; ok {
		return ok
	}

	return false
}
