package freetype

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"

	"github.com/dhindustries/graal"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	//fontlib "golang.org/x/image/font"
)

func loadTrueTypeFont(api graal.Api, path string, size uint, ranges Ranges) (graal.Image, map[rune]*graal.Glyph, error) {
	file, err := api.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	defer api.Release(file)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	font, err := truetype.Parse(data)
	if err != nil {
		return nil, nil, err
	}

	glyphsCount := ranges.Length()
	cols := uint32(math.Sqrt(float64(glyphsCount))) + 1
	rows := (glyphsCount / cols) + 1
	bounds := font.Bounds(fixed.Int26_6(size))
	glyphWidth := uint32(bounds.Max.X - bounds.Min.X)
	glyphHeight := uint32(bounds.Max.Y - bounds.Min.Y)
	imageWidth := pow2(glyphWidth * rows)
	imageHeight := pow2(glyphHeight * cols)
	fmt.Printf("count: %d\n", glyphsCount)
	fmt.Printf("cols: %d\n", cols)
	fmt.Printf("rows: %d\n", rows)
	fmt.Printf("glyphWidth: %d\n", glyphWidth)
	fmt.Printf("glyphHeight: %d\n", glyphHeight)
	fmt.Printf("imageWidth: %d\n", imageWidth)
	fmt.Printf("imageHeight: %d\n", imageHeight)

	target := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))
	context := freetype.NewContext()
	context.SetDPI(72)
	context.SetFont(font)
	context.SetFontSize(float64(size))
	context.SetClip(target.Bounds())
	context.SetDst(target)
	context.SetSrc(image.White)

	glyphs := make(map[rune]*graal.Glyph, glyphsCount)
	var index uint32
	var x, y int

	for _, rang := range ranges {
		for ch := rang[0]; ch <= rang[1]; ch++ {
			if _, ok := glyphs[ch]; !ok {
				metrics := font.HMetric(fixed.Int26_6(size), font.Index(ch))
				glyph := &graal.Glyph{}
				left := float64(x)
				top := float64(y - int(glyphHeight/2))
				right := left + float64(glyphWidth)
				bottom := top + float64(glyphHeight)
				//fmt.Printf("char '%c' (left: %f, top: %f, right: %f, bottom: %f)\n", ch, left, top, right, bottom)

				glyph.Advance = float64(metrics.AdvanceWidth) / float64(size)
				glyph.Left = left / float64(imageWidth)
				glyph.Top = top / float64(imageHeight)
				glyph.Right = right / float64(imageWidth)
				glyph.Bottom = bottom / float64(imageHeight)
				glyph.Size = [2]float64{float64(glyphWidth) / float64(size), float64(glyphHeight) / float64(size)}
				glyphs[ch] = glyph

				pt := freetype.Pt(x, y+int(context.PointToFixed(float64(size))>>8))
				context.DrawString(string(ch), pt)

				if index%cols == 0 {
					x = 0
					y += int(glyphHeight)
				} else {
					x += int(glyphWidth)
				}
				index++
			}
		}
	}

	f, err := os.Create("font.jpeg")
	if err != nil {
		panic(err)
	}
	if err = jpeg.Encode(f, target, nil); err != nil {
		panic(err)
	}
	f.Close()
	//panic("font written")

	img, err := api.Image(target)
	if err != nil {
		return nil, nil, err
	}
	return img, glyphs, nil
}
