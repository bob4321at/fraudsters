package utils

import (
	"bytes"
	"io"
	"math"
	"sync"

	"compress/zlib"
)

var Mouse_X = 0.0
var Mouse_Y = 0.0

type Vec2 struct {
	X, Y float64
}

func DecodeBinary(compressedBytes []byte) []byte {
	reader := bytes.NewReader(compressedBytes)

	r, err := zlib.NewReader(reader)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return decompressed
}

func EncodeBinary(byte []byte) []byte {
	var buf bytes.Buffer

	w := zlib.NewWriter(&buf)

	_, err := w.Write(byte)
	if err != nil {
		panic(err)
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func GetDist(point_1, point_2 Vec2) float64 {
	offx := math.Abs(point_1.X - point_2.X)
	offy := math.Abs(point_1.Y - point_2.Y)

	return math.Sqrt((offx * offx) + (offy * offy))
}

func Deg2Rad(num float64) float64 {
	return num * (180 / 3.14159)
}
func Rad2Deg(num float64) float64 {
	return num * (3.14159 / 180)
}

func RemoveArrayElement[T any](index_to_remove int, slice *[]T) {
	*slice = append((*slice)[:index_to_remove], (*slice)[index_to_remove+1:]...)
}

func Collide(pos1, size1, pos2, size2 Vec2) bool {
	if pos1.X < pos2.X+size2.X && pos1.X+size1.X > pos2.X {
		if pos1.Y < pos2.Y+size2.Y && pos1.Y+size1.Y > pos2.Y {
			return true
		}
	}
	return false
}

func GetSyncLength(syncc *sync.Map) int {
	length := 0

	syncc.Range(func(key, value any) bool {
		length += 1
		return true
	})

	return length
}
