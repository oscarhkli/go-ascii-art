package main

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestBrightnessNumbers(t *testing.T) {
	width := 2
	height := 3
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	red := color.RGBA{255, 0, 0, 0xff}
	cyan := color.RGBA{100, 200, 200, 0xff}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			switch {
			case x < width/2:
				img.Set(x, y, red)
			default:
				img.Set(x, y, cyan)
			}
		}
	}

	if got, want := CalcBrightnessNumbers(img, avg), [][]int{
		{21845, 42833},
		{21845, 42833},
		{21845, 42833},
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("CalcBrightnessNumbers(%v) = %v, want %v", "image with left half red and right half cyan", got, want)
	}
}

func TestBrightnessBlack(t *testing.T) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{1, 1}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	img.Set(0, 0, color.Black)

	if got, want := CalcBrightnessNumbers(img, avg), [][]int{{0}}; !reflect.DeepEqual(got, want) {
		t.Errorf("CalcBrightnessNumbers(%v) = %v, want %v", "image with single white pixel", got, want)
	}
}

func TestBrightnessWhite(t *testing.T) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{1, 1}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	img.Set(0, 0, color.White)

	if got, want := CalcBrightnessNumbers(img, avg), [][]int{{65535}}; !reflect.DeepEqual(got, want) {
		t.Errorf("CalcBrightnessNumbers(%v) = %v, want %v", "image with white black pixel", got, want)
	}
}

func TestCreateASCIIArtBlack(t *testing.T) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{1, 1}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	img.Set(0, 0, color.Black)

	if got, want := CreateASCIIArt(img, avg), "``\n"; got != want {
		t.Errorf("CreateASCIIArt(%v) = %v, want %v", "image with single black pixel", got, want)
	}
}

func TestCreateASCIIArtWhite(t *testing.T) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{1, 1}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	img.Set(0, 0, color.White)

	if got, want := CreateASCIIArt(img, avg), "@@\n"; got != want {
		t.Errorf("CreateASCIIArt(%v) = %v, want %v", "image with single white pixel", got, want)
	}
}
