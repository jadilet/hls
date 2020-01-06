package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

const indexFile = "index.m3u8"

var errNoVideo = errors.New("no video found")

// Video struct
type Video struct {
	Name string `redis:"name"` // file name
	Path string `redis:"path"` // file path
	Dir  string `redis:"dir"`  // file directory
	Ext  string `redis:"ext"`  // extension
	Key  string `redis:"key"`  // key for redis
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// check the file exists "index.m3u8"
func (video *Video) indexFileExist() bool {
	return fileExists(filepath.Join(video.Dir, indexFile))
}

// return video index.m3u8 file path
func (video *Video) indexFilePath() string {
	return filepath.Join(video.Dir, indexFile)
}

// return video segment "index*.ts" file path
func (video *Video) tsFilePath(fileName string) string {
	return filepath.Join(video.Dir, fileName)
}

//check the ts file exists "index*.ts"
func (video *Video) tsFileExist(tsFile string) bool {
	return fileExists(filepath.Join(video.Dir, tsFile))
}

func findVideo(key string) (*Video, error) {
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", strings.ToUpper(key)))

	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, errNoVideo
	}

	var video Video
	err = redis.ScanStruct(values, &video)

	if err != nil {
		return nil, err
	}

	return &video, nil
}
