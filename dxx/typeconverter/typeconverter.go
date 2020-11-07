package typeconverter

import (
	"bytes"
	"encoding/binary"
)

func BytesToFloat64(b []byte) (float64, error) {
	var v float64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &v)
	return v, err
}

func Float64ToBytes(v float64) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, v)
	return buf.Bytes(), err
}

func BytesToFloat32(b []byte) (float32, error) {
	var v float32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &v)
	return v, err
}

func Float32ToBytes(v float32) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, v)
	return buf.Bytes(), err
}

func BytesToInt16(b []byte) (int16, error) {
	var v int16
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &v)
	return v, err
}

func Int16ToBytes(v int16) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, v)
	return buf.Bytes(), err
}

func Float32sToInt16s(data []float32) []int16 {
	const amp = 1<<(16-1) - 1
	absData := AbsFloat32s(data)
	max := MaxFloat32s(absData)
	min := MinFloat32s(absData)

	ret := make([]int16, len(data))
	for i, v := range data {
		vv := int16((AbsFloat32(v) - min) / (max - min) * amp)
		if v < 0 {
			vv = -vv
		}
		ret[i] = vv
	}
	return ret
}

func Float64sToInt16s(data []float64) []int16 {
	const amp = 1<<(16-1) - 1 // default amp for .DSX
	absData := AbsFloat64s(data)
	max := MaxFloat64s(absData)
	min := MinFloat64s(absData)

	ret := make([]int16, len(data))
	for i, v := range data {
		vv := int16((AbsFloat64(v) - min) / (max - min) * amp)
		if v < 0 {
			vv = -vv
		}
		ret[i] = vv
	}
	return ret
}

func Int16sToFloat32s(data []int16) []float32 {
	const amp = 10000.0 // default amp for .DFX
	absData := AbsInt16s(data)
	max := MaxInt16s(absData)
	min := MinInt16s(absData)

	ret := make([]float32, len(data))
	for i, v := range data {
		vv := float32(AbsInt16(v)-min) / float32(max-min) * amp
		if v < 0 {
			vv = -vv
		}
		ret[i] = vv
	}
	return ret
}

func Int16sToFloat64s(data []int16) []float64 {
	const amp = 10000.0 // default amp for .DDX
	absData := AbsInt16s(data)
	max := MaxInt16s(absData)
	min := MinInt16s(absData)

	ret := make([]float64, len(data))
	for i, v := range data {
		vv := float64(AbsInt16(v)-min) / float64(max-min) * amp
		ret[i] = vv
		if v < 0 {
			vv = -vv
		}
	}
	return ret
}

func Float32sToFloat64s(data []float32) []float64 {
	const amp = 10000.0 // default amp for .DDX
	absData := AbsFloat32s(data)
	max := MaxFloat32s(absData)
	min := MinFloat32s(absData)

	ret := make([]float64, len(data))
	for i, v := range data {
		vv := float64(AbsFloat32(v)-min) / float64(max-min) * amp
		ret[i] = vv
		if v < 0 {
			vv = -vv
		}
	}
	return ret
}

func Float64sToFloat32s(data []float64) []float32 {
	const amp = 10000.0 // default amp for .DDX
	absData := AbsFloat64s(data)
	max := MaxFloat64s(absData)
	min := MinFloat64s(absData)

	ret := make([]float32, len(data))
	for i, v := range data {
		vv := float32((AbsFloat64(v) - min) / (max - min) * amp)
		ret[i] = vv
		if v < 0 {
			vv = -vv
		}
	}
	return ret
}

func AbsInt16(x int16) int16 {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func AbsFloat32(x float32) float32 {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func AbsFloat64(x float64) float64 {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func AbsInt16s(data []int16) []int16 {
	ret := make([]int16, len(data))
	for i, v := range data {
		ret[i] = AbsInt16(v)
	}
	return ret
}

func AbsFloat32s(data []float32) []float32 {
	ret := make([]float32, len(data))
	for i, v := range data {
		ret[i] = AbsFloat32(v)
	}
	return ret
}

func AbsFloat64s(data []float64) []float64 {
	ret := make([]float64, len(data))
	for i, v := range data {
		ret[i] = AbsFloat64(v)
	}
	return ret
}

func MaxInt16s(data []int16) int16 {
	var max int16
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func MinInt16s(data []int16) int16 {
	var min int16
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

func MaxFloat32s(data []float32) float32 {
	var max float32
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func MinFloat32s(data []float32) float32 {
	var min float32
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

func MaxFloat64s(data []float64) float64 {
	var max float64
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func MinFloat64s(data []float64) float64 {
	var min float64
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}
