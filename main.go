package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/morgulbrut/color256"
	"github.com/morgulbrut/wpgen/sketch"
)

func main() {

	var (
		sourceImg       string
		totalCycleCount int
		query           string
		params          sketch.UserParams
	)
	flag.StringVar(&sourceImg, "i", "", "Input file (if not specified the program grabs one from the internet")
	flag.IntVar(&totalCycleCount, "cyclecount", 5000, "")
	flag.IntVar(&params.DestHeight, "h",1080, "Height")
	flag.IntVar(&params.DestWidth, "w", 1920, "Width")
	flag.Float64Var(&params.StrokeRatio, "strokeratio", 0.75, "")
	flag.Float64Var(&params.InitialAlpha, "initialalpha", 0.1, "")
	flag.Float64Var(&params.StrokeReduction, "strokereduction", 0.002, "")
	flag.Float64Var(&params.AlphaIncrease, "alphaincrease", 0.06, "")
	flag.Float64Var(&params.StrokeInversionThreshold, "strokeinversiontreshold", 0.05, "")
	flag.IntVar(&params.StrokeJitter, "j", 200, "Stroke jitter (Width / 10 seems a good starting point)")
	flag.IntVar(&params.MinEdgeCount, "min", 3, "Minimal edge count for the polygons")
	flag.IntVar(&params.MaxEdgeCount, "max", 8, "Maximal edge count for the polygons")
	flag.Float64Var(&params.RotationJitter, "r", 0.5, "Rotation jitter")
	flag.StringVar(&params.Shape, "s", "polygon", "Shape of the elements, can be circle, square, hexagon roundedsquare or polygon")
	flag.BoolVar(&params.Fill, "nofill", false, "Don't fill the elements")
	flag.BoolVar(&params.Stroke, "nostroke", false, "Don't draw strokes around the elements")
	flag.StringVar(&query, "q", "blue", "Query for the image")

	flag.Parse()
	var err error
	var img image.Image
	if sourceImg == "" {
		img, err = loadRandomUnsplashImage(params.DestWidth, params.DestHeight,query)
	} else {
		img, err = loadImage(sourceImg)
	}
	if err != nil {
		log.Panicln(err)
	}

	s := sketch.NewSketch(img, params)

	rand.Seed(time.Now().Unix())
	color256.PrintHiCyan("Generating image")

	for i := 0; i < totalCycleCount; i++ {
		s.Update()
	}

	fn := tempFileName(params.Shape+"_", ".png")
	color256.PrintHiCyan("\nWriting %s", fn)
	saveOutput(s.Output(), fn)
}

func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	return filepath.Join(prefix + hex.EncodeToString(randBytes) + suffix)
}

func loadRandomUnsplashImage(width, height int,query string) (image.Image, error) {
	// url := fmt.Sprintf("https://api.unsplash.com/random/%dx%d/?client_id=%s", width, height,accessKey)
	url := fmt.Sprintf("https://loremflickr.com/%d/%d/%s", width, height,query)
	color256.PrintHiCyan("Fetching %s", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	return img, err
}

func loadImage(src string) (image.Image, error) {
	file, _ := os.Open(src)
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func saveOutput(img image.Image, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode to `PNG` with `DefaultCompression` level
	// then save to file
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}
