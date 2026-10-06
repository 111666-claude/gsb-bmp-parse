package bmp

import (
	"encoding/binary"
	"testing"
)

func header(width, height, bits, palette, offset int) []byte {
	buf := make([]byte, 54)
	copy(buf, "BM")
	binary.LittleEndian.PutUint32(buf[10:14], uint32(offset))
	binary.LittleEndian.PutUint32(buf[18:22], uint32(width))
	binary.LittleEndian.PutUint32(buf[22:26], uint32(int32(height)))
	binary.LittleEndian.PutUint16(buf[28:30], uint16(bits))
	binary.LittleEndian.PutUint32(buf[46:50], uint32(palette))
	return buf
}

func TestSingleRow24Bit(t *testing.T) {
	data := append(header(4, 1, 24, 0, 54), 1, 2, 3, 0, 4, 5, 6, 0, 7, 8, 9, 0)
	got, err := Parse(data, &Counter{})
	if err != nil {
		t.Fatalf("解析失败：%v", err)
	}
	if got.Width != 4 || got.Height != 1 || got.Bits != 24 {
		t.Fatalf("头字段解析不对：%+v", got)
	}
	if got.Checksum != 1*1+2*2+3*3+0*4+4*5+5*6+6*7+0*8+7*9+8*10+9*11+0*12 {
		t.Fatalf("校验和不对：%d", got.Checksum)
	}
}

func TestBadMagic(t *testing.T) {
	if _, err := Parse([]byte("XX............"), &Counter{}); err == nil {
		t.Fatal("非 BMP 数据应报错")
	}
}

func TestTruncatedData(t *testing.T) {
	if _, err := Parse([]byte("BM"), &Counter{}); err == nil {
		t.Fatal("太短的数据应报错")
	}
}

func TestPaletteCountReported(t *testing.T) {
	data := append(header(4, 1, 8, 256, 54), 0, 0, 0, 0)
	got, err := Parse(data, &Counter{})
	if err != nil {
		t.Fatalf("解析失败：%v", err)
	}
	if got.Palette != 256 {
		t.Fatalf("调色板条目数应报 256，得到 %d", got.Palette)
	}
}

func TestAbs(t *testing.T) {
	if Abs(-3) != 3 || Abs(3) != 3 {
		t.Fatal("Abs 不对")
	}
}
