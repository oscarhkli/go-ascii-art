package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
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

func GenASCIIArt(img image.Image) string {
	l := len(grayscale)
	var sb strings.Builder
	for _, row := range CalcBrightnessNumbers(img) {
		for _, col := range row {
			s := string(grayscale[col*l/maxBrightness])
			for i := 0; i < repeat; i++ {
				fmt.Fprint(&sb, s)
			}
		}
		fmt.Fprintln(&sb)
	}

	return sb.String()
}

func CalcBrightnessNumbers(img image.Image) [][]int {
	height := img.Bounds().Dy()
	width := img.Bounds().Dx()
	brightnessNumbers := make([][]int, height)
	for i := range brightnessNumbers {
		brightnessNumbers[i] = make([]int, width)
	}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			brightnessNumbers[y][x] = (int(r) + int(g) + int(b)) / 3
		}
	}
	return brightnessNumbers
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("Image path isn't provided")
	}

	img, closer, err := LoadImage(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer closer()

	log.Println("Image loaded with size", img.Bounds().Size())

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

	n, err := f.WriteString(GenASCIIArt(img))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ASCI Art generated: %d bytes\n", n)

	f.Sync()
}
