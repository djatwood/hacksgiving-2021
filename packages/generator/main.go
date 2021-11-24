package generator

import "github.com/ojrac/opensimplex-go"

func GenerateRow(gen opensimplex.Noise32, simple bool, width, height int, scale, cursor float32) []uint8 {
	row := make([]uint8, width*4)
	y := cursor / scale
	for col := 0; col < width; col++ {
		x := float32(col%width) / scale
		value := gen.Eval2(x, y)
		color := GetColor(value, simple)
		i := col * 4
		row[i] = color[0]
		row[i+1] = color[1]
		row[i+2] = color[2]
		row[i+3] = color[3]
	}
	return row
}
