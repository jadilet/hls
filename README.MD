## HLS video streaming and transcoding

### `segmentation` part:
 <p>This part is responsible for transcoding and segmenting using ffmpeg.</p>
  <p>It monitors the directory, and start
transcoding if new files show up.</p>
<p>It uses this directory for the video files "/var/videos".</p>

### `streaming` part:
<p>This part is responsible for HLS streaming.</p>
 <p>It can be accessible by the following link for the streaming the video using VLC player.</p>
 <p>http://address:8080/VIDEOFILENAME_EXTENSION/index.m3u8</p>
 <p>For instance: http://address:8080/bbb_mov/index.m3u8</p>
 
 ### INSTALL:
 Please install the `Redis server` and `FFmpeg` applications.<br>
 `Redis server` used for saving the video information. <br>
 `FFmpeg` is used for transcoding and segmenting.
 ### RUN:
 Please create a video directory in this path `/var/videos` and set permission for writing and reading. By default the application uses this directory for transcoding and streaming the video. <br>
 
 For transcoding and segmenting:<br>
 `*hls/segementation$ go run main.go` <br>
 
 For HLS streaming: <br>
 `*hls/streaming$ go run .`<br>
