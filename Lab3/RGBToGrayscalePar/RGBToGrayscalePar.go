package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"runtime"
	"time"
)

func grayScaleWorker(img image.Image, gray *image.Gray, bounds image.Rectangle, done chan bool, from int, to int) {
	//Coefficients of R G B to GrayScale
	redCoefficient := 0.299
	greenCoefficient := 0.587
	blueCoefficient := 0.114

	for y := from; y < to; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Retrieving contents of a pixel
			originalColor := img.At(x, y)

			//Extracting R, G, B values
			r, g, b, _ := originalColor.RGBA()

			//Calculating GrayScale
			grayValue := uint8(redCoefficient*float64(r>>8) + greenCoefficient*float64(g>>8) + blueCoefficient*float64(b>>8))
			newColor := color.Gray{Y: grayValue}
			gray.Set(x, y, newColor)
		}
	}

	done <- true
}

func main() {
	var fileNameR, fileNameW string
	numThreads := runtime.GOMAXPROCS(0)
	//Input and Output files using command line arguments
	// if len(os.Args) != 3 {
	// 	log.Println("Usage: go run RGBtoGrayScale.go <file to read> <file to write>")
	// 	os.Exit(1)
	// }
	// fileNameR = os.Args[1]
	// fileNameW = os.Args[2]

	fileNameR = "original.jpg"
	fileNameW = "new.jpg"

	//Reading Input file to an image
	file, err := os.Open(fileNameR)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	//Start timing
	start := time.Now()

	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	block := bounds.Max.Y / numThreads

	done := make(chan bool, 1)

	for i := 0; i < numThreads; i++ {
		from := i * block
		to := i*block + block
		if i == (numThreads - 1) {
			to = bounds.Max.Y
		}

		go grayScaleWorker(img, gray, bounds, done, from, to)
	}

	for i := 0; i < numThreads; i++ {
		<-done
	}

	//Stop timing
	elapsedTime := time.Since(start)

	//Saving the modified image to Output file
	outputFile, err := os.Create(fileNameW)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, gray, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done...")
	log.Println("Time taken:", elapsedTime)
}
