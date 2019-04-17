# ffmpeg

## 把X.U..flac转成320kbps的mp3

```bash
cd <ffmpegPath>
ffmpeg.exe -i "\01. X.U..flac"  -map 0:a -b:a 320k "\ X.U. .mp3"
```

## 从01. X.U..flac的56秒开始截取38秒音频,输出到X.mp3

```bash
ffmpeg.exe  -ss 56  -t   38   -i "\01. X.U..flac"  -map 0:a -b:a 320k "d:\ X.U.ls .mp3"
```

