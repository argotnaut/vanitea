package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/argotnaut/vanitea/utils"
	"github.com/creack/pty"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/mat/besticon/ico"
	"github.com/mattn/go-isatty"
	"github.com/nsf/termbox-go"
)

func GetTerminalSize() (width int, height int, err error) {
	cols := 80
	rows := 24

	OperatingSystem := runtime.GOOS
	switch OperatingSystem {
	case "windows":
		if err := termbox.Init(); err != nil {
			panic(err)
		}
		w, h := termbox.Size()
		termbox.Close()
		return w, h, nil
	default:
		if isatty.IsTerminal(os.Stdout.Fd()) {
			if rows, cols, err = pty.Getsize(os.Stdout); err != nil {
				return -1, -1, fmt.Errorf("failed to get the terminal size: %v", err)
			}
		}
	}
	return cols, rows, nil
}

const RESIZE_OFFSET_Y = 8
const RESIZE_FACTOR_Y = 2
const RESIZE_FACTOR_X = 1
const DEFAULT_TERM_COLS = 80
const DEFAULT_TERM_ROWS = 24

const ANSI_BG_TRANSPARENT_COLOR = "\x1b[0;39;49m"
const ANSI_BG_RGB_COLOR = "\x1b[48;2;%d;%d;%dm"
const ANSI_FG_TRANSPARENT_COLOR = "\x1b[0m "
const ANSI_FG_RGB_COLOR = "\x1b[38;2;%d;%d;%dm▄"
const ANSI_RESET = "\x1b[0m"

var InterpolationType = imaging.Lanczos

func decode(buf []byte) []image.Image {
	mime, err := mimetype.DetectReader(bytes.NewReader(buf))
	if err != nil {
		log.Panicf("failed to detect the mime type: %v", err)
	}

	allowed := []string{"image/gif", "image/png", "image/jpeg", "image/bmp", "image/x-icon"}
	if !mimetype.EqualsAny(mime.String(), allowed...) {
		log.Println(string(buf))
		log.Println(len(buf))
		log.Fatalf("invalid MIME type: %s", mime.String())
	}

	frames := make([]image.Image, 0)

	if mime.Is("image/gif") {
		gifImage, err := gif.DecodeAll(bytes.NewReader(buf))

		if err != nil {
			log.Panicf("failed to decode the gif: %v", err)
		}

		var lowestX int
		var lowestY int
		var highestX int
		var highestY int

		for _, img := range gifImage.Image {
			if img.Rect.Min.X < lowestX {
				lowestX = img.Rect.Min.X
			}
			if img.Rect.Min.Y < lowestY {
				lowestY = img.Rect.Min.Y
			}
			if img.Rect.Max.X > highestX {
				highestX = img.Rect.Max.X
			}
			if img.Rect.Max.Y > highestY {
				highestY = img.Rect.Max.Y
			}
		}

		imgWidth := highestX - lowestX
		imgHeight := highestY - lowestY

		overPaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
		draw.Draw(overPaintImage, overPaintImage.Bounds(), gifImage.Image[0], image.Point{}, draw.Src)

		for _, srcImg := range gifImage.Image {
			draw.Draw(overPaintImage, overPaintImage.Bounds(), srcImg, image.Point{}, draw.Over)
			frame := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
			draw.Draw(frame, frame.Bounds(), overPaintImage, image.Point{}, draw.Over)
			frames = append(frames, frame)
		}

		return frames
	} else {
		var frame image.Image
		var err error

		if mime.Is("image/x-icon") {
			frame, err = ico.Decode(bytes.NewReader(buf))
		} else {
			frame, _, err = image.Decode(bytes.NewReader(buf))
		}

		if err != nil {
			log.Panicf("failed to decode the image: %v", err)
		}

		imb := frame.Bounds()
		if imb.Max.X < 2 || imb.Max.Y < 2 {
			log.Fatal("the input image is to small")
		}

		frames = append(frames, frame)
	}

	return frames
}

func scale(frames []image.Image, scaleDivisor int) []image.Image {
	type data struct {
		i  int
		im image.Image
	}

	var err error

	cols := DEFAULT_TERM_COLS
	rows := DEFAULT_TERM_ROWS

	if isatty.IsTerminal(os.Stdout.Fd()) {
		rows, cols, err = GetTerminalSize()
		if err != nil {
			log.Panicf("failed to get the terminal size: %v", err)
		}
	}

	w := cols * RESIZE_FACTOR_X
	h := (rows - RESIZE_OFFSET_Y) * RESIZE_FACTOR_Y

	l := len(frames)
	r := make([]image.Image, l)
	c := make(chan *data, l)

	for i, f := range frames {
		go func(i int, f image.Image) {
			c <- &data{i, imaging.Fit(f, w/scaleDivisor, h/scaleDivisor, InterpolationType)}
		}(i, f)
	}

	for range r {
		d := <-c
		r[d.i] = d.im
	}

	return r
}

func escape(frames []image.Image) [][]string {
	type data struct {
		i   int
		str string
	}

	escaped := make([][]string, 0)

	for _, f := range frames {
		imb := f.Bounds()
		maxY := imb.Max.Y - imb.Max.Y%2
		maxX := imb.Max.X

		c := make(chan *data, maxY/2)
		lines := make([]string, maxY/2)

		for y := 0; y < maxY; y += 2 {
			go func(y int) {
				var sb strings.Builder

				for x := 0; x < maxX; x++ {
					r, g, b, a := f.At(x, y).RGBA()
					if a>>8 < 128 {
						sb.WriteString(ANSI_BG_TRANSPARENT_COLOR)
					} else {
						sb.WriteString(fmt.Sprintf(ANSI_BG_RGB_COLOR, r>>8, g>>8, b>>8))
					}

					r, g, b, a = f.At(x, y+1).RGBA()
					if a>>8 < 128 {
						sb.WriteString(ANSI_FG_TRANSPARENT_COLOR)
					} else {
						sb.WriteString(fmt.Sprintf(ANSI_FG_RGB_COLOR, r>>8, g>>8, b>>8))
					}
				}

				sb.WriteString(ANSI_RESET)
				sb.WriteString("\n")

				c <- &data{y / 2, sb.String()}
			}(y)
		}

		for range lines {
			line := <-c
			lines[line.i] = line.str
		}

		escaped = append(escaped, lines)
	}

	return escaped
}

func GetScaledImage(image []byte, scaleFactor int) string {
	output := escape(scale(decode(image), scaleFactor))
	if len(image) < 1 {
		return "[-]"
	}
	return strings.Join(output[0], "")
}

//go:embed Girl_With_a_Pearl_Earing.jpg
var topImageBytes []byte

//go:embed Volga_Boatmen.jpg
var bottomImageBytes []byte

func main() {
	// Get the top image as a byte array in the form of colored ascii art
	topImage := GetScaledImage(topImageBytes, 2)
	// Get the bottom image as a byte array in the form of colored ascii art
	bottomImage := GetScaledImage(bottomImageBytes, 1)
	// Create the scrollview
	width, height, err := GetTerminalSize()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%s\n", bottomImage)
	// fmt.Printf("%s\n", topImage)
	fmt.Printf("%s\n", utils.PlaceStacked(bottomImage, topImage, utils.CENTER, 0, 0, height, width))
}
