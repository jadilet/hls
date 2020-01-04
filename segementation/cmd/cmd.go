package cmd

// FfmpegPath binary path
var FfmpegPath string

// Params return ffmpeg command parameter
func Params(videoPath string) []string {
	return []string{
		"-i",
		videoPath,
		"-profile:v",
		"baseline",
		"-level",
		"3.0",
		"-s",
		"640x360",
		"-start_number",
		"0",
		"-hls_time",
		"10",
		"-hls_list_size",
		"0",
		"-f",
		"hls",
		"index.m3u8",
	}
}
