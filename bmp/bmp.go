// Package bmp 解析 BMP 位图并算出像素校验和。
package bmp

import (
	"encoding/binary"
	"errors"
)

// Counter 统计本轮检查过的调色板条目字节数。
type Counter struct {
	Scanned int
}

// Result 是一次解析的输出。
type Result struct {
	Width    int `json:"width"`
	Height   int `json:"height"`
	Bits     int `json:"bits"`
	Palette  int `json:"palette"`
	Checksum int `json:"checksum"`
	Scanned  int `json:"scanned"`
}

// Abs 返回绝对值。
func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// Parse 解析 BMP 数据。校验和按逐行逐字节的顺序，用「字节 × 序号」累加。
func Parse(data []byte, c *Counter) (Result, error) {
	if len(data) < 54 || data[0] != 'B' || data[1] != 'M' {
		return Result{}, errors.New("不是 BMP 数据")
	}
	width := int(binary.LittleEndian.Uint32(data[18:22]))
	height := int(int32(binary.LittleEndian.Uint32(data[22:26])))
	bits := int(binary.LittleEndian.Uint16(data[28:30]))
	palette := int(binary.LittleEndian.Uint32(data[46:50]))
	start := int(binary.LittleEndian.Uint32(data[10:14]))
	stride := (width*bits/8 + 3) / 4 * 4
	rows := Abs(height)

	pixel := func(b byte) int {
		c.Scanned++
		return int(b)
	}
	if bits == 8 {
		pixel = func(b byte) int {
			k := int(b)
			if k >= palette {
				return 0
			}
			off := 54 + k*4
			c.Scanned += 4
			if off+2 >= len(data) {
				return 0
			}
			return int(data[off]) + int(data[off+1]) + int(data[off+2])
		}
	}

	sum := 0
	n := 0
	for y := 0; y < rows; y++ {
		srcRow := y
		if height > 0 {
			srcRow = rows - 1 - y
		}
		src := start + srcRow*stride
		for x := 0; x < stride; x++ {
			p := src + x
			if p >= len(data) {
				break
			}
			n++
			v := int(data[p])
			if bits == 8 && x < width {
				v = pixel(data[p])
			}
			sum += v * n
		}
	}
	return Result{
		Width: width, Height: height, Bits: bits,
		Palette: palette, Checksum: sum, Scanned: c.Scanned,
	}, nil
}
