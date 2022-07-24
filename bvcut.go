/* Black Video Cut */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/cheggaaa/pb/v3"
	"gocv.io/x/gocv"
)

func main() {
	var (
		per = flag.Float64("per", 0.5, "per sec")
		enc = flag.Bool("enc", true, "encode")
	)
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: bvcut [-per <float>][enc=<bool>] <videopath>")
		os.Exit(1)
	}

	filename := flag.Arg(0)

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
	var perSec float64 = *per
	var intervalSec float64 = 5.0
	secCount := intervalSec
	blackSecs := []float64{0.0}
	progressBar := pb.StartNew(int(totalFrame))

	fmt.Println(fps, "fps, total", int(totalFrame), "frames")

	for {
		if !video.IsOpened() || secCount*float64(fps) >= totalFrame {
			break
		} else {
			video.Set(gocv.VideoCapturePosFrames, secCount*float64(fps))
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
		outputFilename := "output_" + strconv.Itoa(i-1) + ".mp4"
		err := createVideo(filename, outputFilename, blackSecs[i-1], blackSecs[i], *enc)
		if err != nil {
			fmt.Println("cut video error:", err.Error())
			os.Exit(1)
		}
	}
}

func toMinute(sec float64) string {
	return fmt.Sprintf("%d:%02d", int(sec/60.0), (int(sec) % 60))
}

func createVideo(srcFilename string, dstFilename string, startSec float64, endSec float64, enc bool) error {
	start := strconv.FormatFloat(startSec, 'f', 8, 64)
	end := strconv.FormatFloat(endSec-startSec, 'f', 8, 64)
	fmt.Println(dstFilename + ": " + toMinute(startSec) + " ~ " + toMinute(endSec))

	var err error
	if enc {
		// encode with enc
		err = exec.Command("ffmpeg", "-ss", start, "-t", end, "-i", srcFilename,
			/*"-r", "60",*/ "-vsync", "1",
			"-c:v", "libx264", "-b:v", "12M",
			"-c:a", "libfdk_aac", "-b:a", "128k", dstFilename).Run()
	} else {
		err = exec.Command("ffmpeg", "-ss", start, "-t", end, "-i", srcFilename,
			"-c", "copy", dstFilename).Run()
	}
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
