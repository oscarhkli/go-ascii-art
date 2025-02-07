package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	grayscale     string = "`.-':_,^=;><+!rc*/z?sLTv)J7(|Fi{C}fI31tlu[neoZ5Yxjya]2ESwqkP6h9d4VpOGbUAKXHm8RD#$Bg0MNWQ%&@"
	maxBrightness int    = 65535
	repeat        int    = 2
)

func LoadImage(imgPath string) (image.Image, func(), error) {
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		f.Close()
	}

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, closer, err
	}

	if format != "jpeg" && format != "png" && format != "gif" {
		return nil, closer, fmt.Errorf("unsupported format: %s", format)
	}

	return img, closer, nil
}

func CreateASCIIArt(img image.Image, bCalc BrightnessCalc) string {
	l := len(grayscale)
	var sb strings.Builder
	for _, row := range CalcBrightnessNumbers(img, bCalc) {
		for _, col := range row {
			s := string(grayscale[l*col/maxBrightness])
			for range repeat {
				fmt.Fprint(&sb, s)
			}
		}
		fmt.Fprintln(&sb)
	}

	return sb.String()
}

// Color to brightness in greyscale reference: https://alienryderflex.com/hsp.html
type BrightnessCalc func(int, int, int) int

func avg(r, g, b int) int {
	return (int(r) + int(g) + int(b)) / 3
}

func hsl(r, g, b int) int {
	return (max(r, g, b) + min(r, g, b)) / 2
}

func hsp(r, g, b int) int {
	return int(math.Sqrt(0.299*float64(r*r) + 0.587*float64(b*b) + 0.114*float64(b*b)))
}

func CalcBrightnessNumbers(img image.Image, bCalc BrightnessCalc) [][]int {
	height := img.Bounds().Dy()
	width := img.Bounds().Dx()
	brightnessNumbers := make([][]int, height)
	for i := range brightnessNumbers {
		brightnessNumbers[i] = make([]int, width)
	}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			brightnessNumbers[y][x] = bCalc(int(r), int(g), int(b))
		}
	}
	return brightnessNumbers
}

func main() {
	method := ""
	flag.StringVar(&method, "c", "hsl", "Brightness Calculation Method {avg | hsl | hsp}")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Image path isn't provided")
	} else if len(args) > 2 {
		flag.Usage()
		os.Exit(1)
	}

	img, closer, err := LoadImage(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer closer()

	log.Println("Image loaded with size", img.Bounds().Size())

	bCalcMap := map[string]BrightnessCalc{
		"avg": avg,
		"hsl": hsl,
		"hsp": hsp,
	}

	bCalc, ok := bCalcMap[method]
	if !ok {
		log.Println("Invalid Brightness Calculation Method is selected. Use default <hsp>")
		bCalc = hsp
	} else {
		log.Println("Brightness Calculation Method:", method)
	}

	asciiArt := CreateASCIIArt(img, bCalc)

	outputFilePath := "out/ascii.txt"
	if len(args) > 1 {
		outputFilePath = args[1]
	}

	dirPath := filepath.Dir(outputFilePath)
	err = os.MkdirAll(dirPath, 0750)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	size, err := f.WriteString(asciiArt)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ASCI Art generated: %d bytes\n", size)

	f.Sync()
}
