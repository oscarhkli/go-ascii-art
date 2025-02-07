package main

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestLoadImage(t *testing.T) {
	imgPath := "sample.jpg"
	img, closer, err := LoadImage(imgPath)
	defer closer()

	if img == nil && err != nil {
		t.Errorf("LoadImage(%v) = (img: %v, err: %v), want (img: %v, err: %v)", imgPath, nil, "not <nil>", "not <nil>", err)
	}
}

func TestLoadImageNotFound(t *testing.T) {
	imgPath := "abc.jpg"
	img, _, err := LoadImage(imgPath)

	if img != nil && err == nil {
		t.Errorf("LoadImage(%v) = (img: %v, err: %v), want (img: %v, err: %v)", imgPath, "not <nil>", err, nil, "not <nil>")
	}
}

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

	got := CalcBrightnessNumbers(img)
	want := [][]int{
		{21845, 42833},
		{21845, 42833},
		{21845, 42833},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CalcBrightnessNumbers(%v) = %v, want %v", "image with left half red and right half cyan", got, want)
	}
}
