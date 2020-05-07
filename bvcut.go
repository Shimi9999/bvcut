/* Black Video Cut */

package main

import (
  "fmt"
  "os"
  "os/exec"
  "strconv"

  "gocv.io/x/gocv"
  "github.com/cheggaaa/pb/v3"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Println("arg error")
    os.Exit(1)
  }

  filename := os.Args[1]

  video, err := gocv.VideoCaptureFile(filename)
  if err != nil {
    fmt.Println("file open error")
    os.Exit(1)
  }
  defer video.Close()

  imgframe := gocv.NewMat()
  defer imgframe.Close()

  fps := int64(video.Get(gocv.VideoCaptureFPS))
  totalFrame := video.Get(gocv.VideoCaptureFrameCount)
  totalSec := totalFrame / video.Get(gocv.VideoCaptureFPS)
  var perSec float64 = 0.5
  var intervalSec float64 = 5.0
  secCount := intervalSec
  blackSecs := []float64{0.0}
  progressBar := pb.StartNew(int(totalFrame))

  fmt.Println(fps, "fps, total", int(totalFrame), "frames")

  for {
    if !video.IsOpened() || secCount * float64(fps) >= totalFrame {
      break
    } else {
      video.Set(gocv.VideoCapturePosFrames, secCount * float64(fps))
      if video.Read(&imgframe) {
        if isBlackScreen(imgframe) {
          sec := video.Get(gocv.VideoCapturePosMsec) / 1000.0
          blackSecs = append(blackSecs, sec)
          secCount += intervalSec - perSec
        }

        frame := int64(video.Get(gocv.VideoCapturePosFrames))
        progressBar.SetCurrent(frame)
        //printMean(imgframe)
      }
      secCount += perSec
    }
  }
  progressBar.SetCurrent(int64(totalFrame))
  progressBar.Finish()
  blackSecs = append(blackSecs, totalSec)
  fmt.Println("BlackScreen sec:", blackSecs)

  for i := 1; i < len(blackSecs); i++ {
    outputFilename := "output_" + strconv.Itoa(i - 1) + ".mp4"
    err := createVideo(filename, outputFilename, blackSecs[i - 1], blackSecs[i])
    if err != nil {
      fmt.Println("cut video error:", err.Error())
      os.Exit(1)
    }
  }
}

func toMinute(sec float64) string {
  return strconv.Itoa(int(sec / 60.0)) + ":" + strconv.Itoa((int(sec) % 60))
}

func createVideo(srcFilename string, dstFilename string, startSec float64, endSec float64) error {
  start := strconv.FormatFloat(startSec, 'f', 8, 64)
  end := strconv.FormatFloat(endSec - startSec, 'f', 8, 64)
  fmt.Println(dstFilename + ": " + toMinute(startSec) + " ~ " + toMinute(endSec))
  // encode with CBR
  err := exec.Command("ffmpeg", "-ss", start, "-t", end, "-i", srcFilename,
    "-c:v", "libx264", "-b:v", "8000k", dstFilename).Run()
  return err
}

func isBlackScreen(img gocv.Mat) bool {
  imghsv := gocv.NewMat()
  defer imghsv.Close()
  gocv.CvtColor(img, &imghsv, gocv.ColorBGRToHSV)
  if imghsv.Mean().Val3 < 1 {
    return true
  }
  return false
}

func printMean(img gocv.Mat) {
  fmt.Println("bgr: ", img.Mean())

  imghsv := gocv.NewMat()
  defer imghsv.Close()
  gocv.CvtColor(img, &imghsv, gocv.ColorBGRToHSV)
  fmt.Println("hsv: ", imghsv.Mean())
  if imghsv.Mean().Val3 < 1 {
    fmt.Println("BLACK")
  }
}
