package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"sync"
	"time"
)

const XRESOLUTION = 2000
const YRESOLUTION = 2000
const FRAMES = 100
const MAX_ITERATIONS = 50

const MINX float64 = -2
const MAXX float64 = 2

const MINY float64 = -2
const MAXY float64 = 2

var colors [XRESOLUTION][YRESOLUTION]color.Color

var startTime = time.Now()

func main() {
	for frame := 0; frame < FRAMES; frame++ {
		drawFrame(frame)
		var timeLeft = (float64(FRAMES-frame) / float64(frame+1)) * time.Since(startTime).Seconds()
		fmt.Println("frame #"+strconv.Itoa(frame)+" done!", strconv.Itoa(int(timeLeft))+"s left")
	}
	fmt.Println("done! total time elapsed: ", time.Since(startTime))
}

func drawFrame(frame int) {
	var framestring = strconv.Itoa(frame)
	for len(framestring) < 3 {
		framestring = "0" + framestring
	}
	out, err := os.Create("./anim/out" + framestring + ".png")

	imgRect := image.Rect(0, 0, XRESOLUTION, YRESOLUTION)
	img := image.NewRGBA(imgRect)

	var wg sync.WaitGroup

	wg.Add(XRESOLUTION)
	for x := 0; x < XRESOLUTION; x++ {
		go getColumn(x, 2.0*math.Pi*float64(frame)/float64(FRAMES), &wg)
	}

	wg.Wait()

	for x := 0; x < XRESOLUTION; x++ {
		for y := 0; y < YRESOLUTION; y++ {
			img.Set(x, y, colors[x][y])
		}
	}

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println("errors encountered!")
		fmt.Println(err)
		os.Exit(1)
	}
}

func getColumn(x int, phi float64, group *sync.WaitGroup) {
	for y := 0; y < YRESOLUTION; y++ {
		colors[x][y] = getPixel(x, y, phi)
	}
	group.Done()
}

func getPixel(x int, y int, phi float64) color.Color {
	var c = complex(float64(x)/float64(XRESOLUTION)*(MAXX-MINX)+float64(MINX),
		float64(y)/float64(YRESOLUTION)*(MAXY-MINY)+float64(MINY))
	var z complex128 = 0
	var iterations uint8 = 0
	for cmplx.Abs(z) < 5 && iterations < MAX_ITERATIONS {
		z = cmplx.Tan(z+complex(phi, 0))*z*Polar(1, cmplx.Phase(c)+math.Pi) + c
		iterations++
	}
	return color.RGBA{R: iterations * 3, G: iterations * 4, B: iterations * 8, A: 255}
}

func Polar(radius, theta float64) complex128 {
	return complex(radius, 0) * cmplx.Exp(complex(0, theta))
}
