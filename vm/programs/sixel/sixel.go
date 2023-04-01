package sixel

import (
	"fmt"
	"image"
	"log"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/mattn/go-sixel"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"golang.org/x/image/draw"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(img2sixel))
}

func drawScaler(name string) draw.Scaler {
	switch name {
	case "nearest":
		return draw.NearestNeighbor
	case "bilinear":
		return draw.BiLinear
	case "approx-bilinear":
		return draw.ApproxBiLinear
	case "cubic", "catmull-rom":
		return draw.CatmullRom
	default:
		return nil
	}
}

var scalerNames = []string{
	"nearest",
	"bilinear",
	"approx-bilinear",
	"cubic",
	"catmull-rom (same as cubic)",
}

var img2sixel = cli.App{
	Name:      "img2sixel",
	Usage:     "convert JPG/PNG to SIXEL",
	UsageText: `img2sixel [options...] <file>`,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "width",
			Usage: "width of the output image in pixels without upscaling",
			Value: 400,
		},
		&cli.IntFlag{
			Name:  "height",
			Usage: "height of the output image in pixels without upscaling",
			Value: 400,
		},
		&cli.BoolFlag{
			Name: "keep-aspect",
			Usage: "keep the aspect ratio of the original image, " +
				"making --width and --height the maximum values",
			Value: true,
		},
		&cli.StringFlag{
			Name: "scaler",
			Usage: "scaler to use when resizing the image. " +
				"one of: " + strings.Join(scalerNames, ", "),
			Value: "approx-bilinear",
			Action: func(ctx *cli.Context, v string) error {
				if drawScaler(v) == nil {
					return fmt.Errorf("invalid interpolator: %q", v)
				}
				return nil
			},
		},
		&cli.BoolFlag{
			Name:  "dither",
			Usage: "enable dithering",
		},
		&cli.IntFlag{
			Name:  "colors",
			Usage: "number of colors to use in the palette, 2-4096",
			Value: 256,
		},
	},
	Action: func(c *cli.Context) error {
		env := vm.EnvironmentFromContext(c.Context)

		path := c.Args().First()
		if path == "" {
			return &vm.UsageError{Usage: "Usage: " + c.App.UsageText}
		}

		colors := c.Int("colors")
		if colors < 2 || colors > 4096 {
			return errors.Errorf("invalid number of colors: %d", colors)
		}

		f, err := env.Open(path)
		if err != nil {
			return errors.Wrap(err, "open")
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return errors.Wrap(err, "image decode")
		}

		if c.Int("width") != 0 || c.Int("height") != 0 {
			scaler := drawScaler(c.String("scaler"))

			img, err = resize(img, scaler, c.Int("width"), c.Int("height"), c.Bool("keep-aspect"))
			if err != nil {
				return errors.Wrap(err, "resize")
			}
		}

		sixelEnc := sixel.NewEncoder(env.Terminal.Stdout)
		sixelEnc.Dither = c.Bool("dither")
		sixelEnc.Colors = colors
		return sixelEnc.Encode(img)
	},
}

func widthFromHeight(original image.Rectangle, height int) int {
	return original.Dx() * height / original.Dy()
}

func heightFromWidth(original image.Rectangle, width int) int {
	return original.Dy() * width / original.Dx()
}

func resize(img image.Image, scaler draw.Scaler, maxW, maxH int, keepAspect bool) (image.Image, error) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	if keepAspect {
		if w <= maxW && h <= maxH {
			return img, nil
		}
		if w > maxW {
			w = maxW
			h = heightFromWidth(img.Bounds(), w)
		}
		if h > maxH {
			h = maxH
			w = widthFromHeight(img.Bounds(), h)
		}
		log.Println("resized to", w, "x", h)
	} else {
		if w <= maxW && h <= maxH {
			return img, nil
		}
		if maxW == 0 || maxH == 0 {
			return nil, errors.New("invalid --width or --height")
		}
		w = maxW
		h = maxH
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	scaler.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst, nil
}
