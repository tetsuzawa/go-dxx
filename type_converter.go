package dxx

import (
	"bytes"
	"encoding/binary"
)

func bytesToFloat64(b []byte) (float64, error) {
	var v float64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &v)
	return v, err
}

func float64ToBytes(v float64) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, v)
	return buf.Bytes(), err
}

func bytesToFloat32(b []byte) (float32, error) {
	var v float32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &v)
	return v, err
}

func float32ToBytes(v float32) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, v)
	return buf.Bytes(), err
}

func bytesToInt16(b []byte) (int16, error) {
	var v int16
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &v)
	return v, err
}

func int16ToBytes(v int16) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, v)
	return buf.Bytes(), err
}

func float32sToInt16s(data []float32) []int16 {
	const amp = 1<<(BitLenShort-1) - 1
	absData := absFloat32s(data)
	max := maxFloat32s(absData)
	min := minFloat32s(absData)

	ret := make([]int16, 0, len(data))
	for _, v := range data {
		vv := int16((v - min) / (max - min) * amp)
		ret = append(ret, vv)
	}
	return ret
}

func float64sToInt16s(data []float64) []int16 {
	const amp = 1<<(BitLenShort-1) - 1 // default amp for .DSX
	absData := absFloat64s(data)
	max := maxFloat64s(absData)
	min := minFloat64s(absData)

	ret := make([]int16, 0, len(data))
	for _, v := range data {
		vv := int16((v - min) / (max - min) * amp)
		ret = append(ret, vv)
	}
	return ret
}

func int16sToFloat32s(data []int16) []float32 {
	const amp = 10000.0 // default amp for .DFX
	absData := absInt16s(data)
	max := maxInt16s(absData)
	min := minInt16s(absData)

	ret := make([]float32, 0, len(data))
	for _, v := range data {
		vv := float32(v-min) / float32(max-min) * amp
		ret = append(ret, vv)
	}
	return ret
}

func int16sToFloat64s(data []int16) []float64 {
	const amp = 10000.0 // default amp for .DDX
	absData := absInt16s(data)
	max := maxInt16s(absData)
	min := minInt16s(absData)

	ret := make([]float64, 0, len(data))
	for _, v := range data {
		vv := float64(v-min) / float64(max-min) * amp
		ret = append(ret, vv)
	}
	return ret
}

func float32sToFloat64s(data []float32) []float64 {
	const amp = 10000.0 // default amp for .DDX
	absData := absFloat32s(data)
	max := maxFloat32s(absData)
	min := minFloat32s(absData)

	ret := make([]float64, 0, len(data))
	for _, v := range data {
		vv := float64(v-min) / float64(max-min) * amp
		ret = append(ret, vv)
	}
	return ret
}

func float64sToFloat32s(data []float64) []float32 {
	const amp = 10000.0 // default amp for .DDX
	absData := absFloat64s(data)
	max := maxFloat64s(absData)
	min := minFloat64s(absData)

	ret := make([]float32, 0, len(data))
	for _, v := range data {
		vv := float32((v - min) / (max - min) * amp)
		ret = append(ret, vv)
	}
	return ret
}

func absInt16s(data []int16) []int16 {
	ret := make([]int16, 0, cap(data))
	for _, v := range data {
		if v < 0 {
			ret = append(ret, -v)
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}

func absFloat32s(data []float32) []float32 {
	ret := make([]float32, 0, len(data))
	for _, v := range data {
		if v < 0 {
			ret = append(ret, -v)
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}

func absFloat64s(data []float64) []float64 {
	ret := make([]float64, 0, cap(data))
	for _, v := range data {
		if v < 0 {
			ret = append(ret, -v)
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}

func maxInt16s(data []int16) int16 {
	var max int16
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func minInt16s(data []int16) int16 {
	var min int16
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

func maxFloat32s(data []float32) float32 {
	var max float32
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func minFloat32s(data []float32) float32 {
	var min float32
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

func maxFloat64s(data []float64) float64 {
	var max float64
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func minFloat64s(data []float64) float64 {
	var min float64
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}
