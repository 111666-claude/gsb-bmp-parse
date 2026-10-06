// Command bmpcheck 解析 BMP 并打印头信息与校验和。
package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.com/bmp-parse/bmp"
)

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func build(width, height, bits, palette, offset int, pixels []byte) []byte {
	buf := make([]byte, offset)
	copy(buf, "BM")
	binary.LittleEndian.PutUint32(buf[10:14], uint32(offset))
	binary.LittleEndian.PutUint32(buf[18:22], uint32(width))
	binary.LittleEndian.PutUint32(buf[22:26], uint32(int32(height)))
	binary.LittleEndian.PutUint16(buf[28:30], uint16(bits))
	binary.LittleEndian.PutUint32(buf[46:50], uint32(palette))
	return append(buf, pixels...)
}

func stride(width, bits int) int {
	raw := width * bits / 8
	return (raw + 3) / 4 * 4
}

func show(data []byte) {
	got, err := bmp.Parse(data, &bmp.Counter{})
	if err != nil {
		panic(err)
	}
	emit(got)
}

func main() {
	sample := flag.String("sample", "", "padding|topdown|offset|work")
	flag.Parse()

	switch *sample {
	case "padding":
		rowA := make([]byte, stride(3, 24))
		for i := 0; i < 9; i++ {
			rowA[i] = byte(i + 1)
		}
		rowB := make([]byte, stride(3, 24))
		for i := 0; i < 9; i++ {
			rowB[i] = byte(21 + i)
		}
		show(build(3, 2, 24, 0, 54, append(append([]byte{}, rowA...), rowB...)))
	case "topdown":
		rowA := make([]byte, stride(4, 24))
		for i := range rowA {
			rowA[i] = byte(i + 1)
		}
		rowB := make([]byte, stride(4, 24))
		for i := range rowB {
			rowB[i] = byte(40 + i)
		}
		show(build(4, -2, 24, 0, 54, append(append([]byte{}, rowA...), rowB...)))
	case "offset":
		row := make([]byte, stride(4, 24))
		for i := range row {
			row[i] = byte(i + 1)
		}
		show(build(4, 1, 24, 0, 70, row))
	case "work":
		w := 2000
		row := make([]byte, stride(w, 8))
		for i := range row {
			row[i] = byte(i % 256)
		}
		pal := make([]byte, 256*4)
		for k := 0; k < 256; k++ {
			pal[k*4], pal[k*4+1], pal[k*4+2] = byte(k), byte(k), byte(k)
		}
		got, err := bmp.Parse(build(w, 1, 8, 256, 54, append(pal, row...)), &bmp.Counter{})
		if err != nil {
			panic(err)
		}
		emit(map[string]int{"width": w, "scanned": got.Scanned})
	default:
		fmt.Fprintln(os.Stderr, "未知场景："+*sample)
		os.Exit(2)
	}
}
