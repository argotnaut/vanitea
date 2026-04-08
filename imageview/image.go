package imageview

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/mat/besticon/ico"
)

const ANSI_BG_TRANSPARENT_COLOR = "\x1b[0;39;49m"
const ANSI_BG_RGB_COLOR = "\x1b[48;2;%d;%d;%dm"
const ANSI_FG_TRANSPARENT_COLOR = "\x1b[0m "
const ANSI_FG_RGB_COLOR = "\x1b[38;2;%d;%d;%dm▄"
const ANSI_RESET = "\x1b[0m"

var InterpolationType = imaging.Lanczos

/*
Decodes a given slice of bytes into grid of color.Color values. The slice of bytes
is expected to represent an image of one of the following types:
gif, png, jpeg, bmp, or x-icon
*/
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

/*
Given an image and a target size, returns an image scaled to that size.

frames: []image.Image - A slice of grids of color.Color values, representing the image data

targetSize: tea.WindowSizeMsg - The desired size of the output image. Size here is measured

	in "pixels" where, because of the way the image is rendered with ASCII characters, a "pixel"
	is 1x the width and 0.5x the height of a cell in the terminal
*/
func scale(frames []image.Image, targetSize tea.WindowSizeMsg) []image.Image {
	type data struct {
		i  int
		im image.Image
	}

	l := len(frames)
	r := make([]image.Image, l)
	c := make(chan *data, l)

	for i, f := range frames {
		go func(i int, f image.Image) {
			c <- &data{i, imaging.Fit(f, targetSize.Width, targetSize.Height, InterpolationType)}
		}(i, f)
	}

	for range r {
		d := <-c
		r[d.i] = d.im
	}

	return r
}

/*
Given an image in the form of a slice of image.Image, return the ANSI-escaped 2D string slice
*/
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

/*
Given the current image dimensions and the bounding dimensions it needs
to fit into, this function returns the largest dimensions the image could
have while preserving its aspect ratio
*/
func rescaleImageToBounds(imageDimensions tea.WindowSizeMsg, bounds tea.WindowSizeMsg) (output tea.WindowSizeMsg) {
	resizeRatio := min(
		float64(bounds.Width)/float64(imageDimensions.Width),
		float64(bounds.Height)/float64(imageDimensions.Height),
	)
	output.Width = int(float64(imageDimensions.Height) * resizeRatio)
	output.Height = int(float64(imageDimensions.Width) * resizeRatio)
	return output
}

/*
Decodes a slice of bytes representing an image and returns a slice of
image.Image
*/
func decodeImageBytes(image []byte) (output []image.Image, err error) {
	if len(image) < 1 {
		return output, fmt.Errorf("image bytes slice was empty")
	}
	decodedImages := decode(image)
	return decodedImages, nil
}

/*
Returns the width and height dimensions of a given image
*/
func getImageDimensions(image image.Image) tea.WindowSizeMsg {
	return tea.WindowSizeMsg{
		Height: image.Bounds().Dy(),
		Width:  image.Bounds().Dx(),
	}
}

/*
Given an image and a target size, this function returns the image as
an ANSI-escaped ASCII pixel string
*/
func getScaledImage(frames []image.Image, size *tea.WindowSizeMsg) string {
	const ERROR_OUTPUT = "[-]"
	if len(frames) < 1 {
		return ERROR_OUTPUT
	}
	imageDimensions := getImageDimensions(frames[0])
	if size != nil {
		size.Height *= 2 // multiply height by two to convert from characters to "pixels"
		imageDimensions = rescaleImageToBounds(imageDimensions, *size)
	}
	rescaledImage := scale(
		frames,
		imageDimensions,
	)
	output := escape(rescaledImage)
	return strings.Join(output[0], "")
}
