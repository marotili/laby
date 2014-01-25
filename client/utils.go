// Copyright (c) 2012 by Lecture Hall Games Authors.
// All source files are distributed under the Simplified BSD License.

package main

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"os"
	"unsafe"
)

type Vector struct {
	x float32
	y float32
}

func (v Vector) Normalized() Vector {
	length := float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
	if length == 0 {
		return Vector{0, 0}
	}
	return Vector{v.x / length, v.y / length}
}

func (v *Vector) Normalize() {
	length := float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
	if length == 0 {
		v.x = 0
		v.y = 0
	} else {
		v.x = v.x / length
		v.y = v.y / length
	}
}

func (v Vector) Length() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
}

func (v Vector) CrossProd(v2 Vector) float32 {
	return v.x*v2.y - v.y*v2.x
}

func (lhs Vector) Mul(rhs Vector) float32 {
	return lhs.x*rhs.x + lhs.y*rhs.y
}

func (lhs Vector) MulScalar(rhs float32) Vector {
	return Vector{lhs.x * rhs, lhs.y * rhs}
}

func (lhs Vector) DivScalar(rhs float32) Vector {
	return lhs.MulScalar(1 / rhs)
}

func (lhs Vector) Add(rhs Vector) Vector {
	return Vector{lhs.x + rhs.x, lhs.y + rhs.y}
}

func (lhs Vector) Sub(rhs Vector) Vector {
	return lhs.Add(rhs.MulScalar(-1))
}

func (v Vector) Rotate(angle float32) Vector {
	return Vector{
		(v.x*float32(math.Cos(float64(angle))) -
			v.y*float32(math.Sin(float64(angle)))),
		(v.x*float32(math.Sin(float64(angle))) +
			v.y*float32(math.Cos(float64(angle)))),
	}
}

func (v Vector) Project(v2 Vector) (Vector, float32) {
	dot := v.Mul(v2) / v2.Length()
	return v2.Normalized().MulScalar(dot), dot
}

type Sprite struct {
	img    *image.RGBA
	tex    gl.Texture
	width  float32
	height float32
}

func uploadTexture(img *image.RGBA) gl.Texture {
	gl.Enable(gl.TEXTURE_2D)
	tex := gl.GenTexture()
	tex.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, 4, img.Bounds().Max.X, img.Bounds().Max.Y,
		0, gl.RGBA, gl.UNSIGNED_BYTE, img.Pix)
	gl.Disable(gl.TEXTURE_2D)
	return tex
}

func NewSprite(path string, width, height float32) (*Sprite, error) {
	img, err := LoadImageRGBA(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tex := uploadTexture(img)
	return &Sprite{img, tex, width, height}, nil
}

func NewSpriteFromSurface(surface *sdl.Surface) *Sprite {
	img := image.NewRGBA(image.Rect(0, 0, int(surface.W), int(surface.H)))
	b := img.Bounds()
	bpp := int(surface.Format.BytesPerPixel)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pixel := uintptr(unsafe.Pointer(surface.Pixels))
			pixel += uintptr(y*int(surface.Pitch) + x*bpp)
			p := (*color.RGBA)(unsafe.Pointer(pixel))
			img.SetRGBA(x, y, *p)
		}
	}
	tex := uploadTexture(img)
	return &Sprite{img, tex, float32(b.Max.X), float32(b.Max.Y)}
}

func (s *Sprite) Draw(x, y, angle, scale float32, blend bool) {
	gl.Enable(gl.TEXTURE_2D)
	gl.Disable(gl.COLOR_MATERIAL)
	if blend {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
		gl.BlendFunc(gl.ONE, gl.ZERO)
	}

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(x, y, 0)
	gl.Rotatef(angle*360/(2*math.Pi), 0, 0, 1)
	gl.Scalef(scale, scale, 1)
	s.tex.Bind(gl.TEXTURE_2D)
	gl.Begin(gl.QUADS)
	gl.Color3f(1, 1, 1)
	gl.TexCoord2d(0, 0)
	gl.Vertex3f(-0.5*s.width, -0.5*s.height, 0)
	gl.TexCoord2d(1, 0)
	gl.Vertex3f(0.5*s.width, -0.5*s.height, 0)
	gl.TexCoord2d(1, 1)
	gl.Vertex3f(0.5*s.width, 0.5*s.height, 0)
	gl.TexCoord2d(0, 1)
	gl.Vertex3f(-0.5*s.width, 0.5*s.height, 0)
	gl.End()
	gl.Disable(gl.TEXTURE_2D)
	gl.Disable(gl.BLEND)
}

func LoadImageRGBA(path string) (*image.RGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, nil
}

func LoadImageGray(path string) (*image.Gray, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if gray, ok := img.(*image.Gray); ok {
		return gray, nil
	}
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, image.Point{0, 0}, draw.Src)
	return gray, nil
}
