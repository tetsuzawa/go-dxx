package dxx

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//type DataType string
//
//const (
//	DSA DataType = "DSA"
//	DFA DataType = "DFA"
//	DDA DataType = "DDA"
//	DSB DataType = "DSB"
//	DFB DataType = "DFB"
//	DDB DataType = "DDB"
//)

const (
	BitLenShort  = 16
	BitLenFloat  = 32
	BitLenDouble = 64
)

type DataType int

const (
	DSA DataType = iota + 1
	DFA
	DDA
	DSB
	DFB
	DDB
)

func (dt DataType) String() string {
	switch dt {
	case DSA:
		return "DSA"
	case DFA:
		return "DFA"
	case DDA:
		return "DDA"
	case DSB:
		return "DSB"
	case DFB:
		return "DFB"
	case DDB:
		return "DDB"
	default:
		return "unkwon"
	}
}

func StringToDataType(s string) (DataType, error) {
	switch s {
	case "DSA":
		return DSA, nil
	case "DFA":
		return DFA, nil
	case "DDA":
		return DDA, nil
	case "DSB":
		return DSB, nil
	case "DFB":
		return DFB, nil
	case "DDB":
		return DDB, nil
	default:
		return 0, errors.New("Unknown DataType")
	}
}

func (dt DataType) BitLen() int {
	switch dt {
	case DSA:
		return BitLenShort
	case DFA:
		return BitLenFloat
	case DDA:
		return BitLenDouble
	case DSB:
		return BitLenShort
	case DFB:
		return BitLenFloat
	case DDB:
		return BitLenDouble
	default:
		return -1 // will not be called
	}
}
func (dt DataType) ByteLen() int {
	return dt.BitLen() / 8
}

func Read(filename string) (data []float64, err error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	dt, err := StringToDataType(ext(filename))
	if err != nil {
		return nil, err
	}

	switch dt {
	case DSA:
		i16s, err := readDSA(f)
		if err != nil {
			return nil, err
		}
		return int16sToFloat64s(i16s), nil
	case DFA:
		f32s, err := readDFA(f)
		if err != nil {
			return nil, err
		}
		return float32sToFloat64s(f32s), nil
	case DDA:
		f64s, err := readDDA(f)
		if err != nil {
			return nil, err
		}
		return f64s, nil
	case DSB:
		i16s, err := readDSB(f)
		if err != nil {
			return nil, err
		}
		return int16sToFloat64s(i16s), nil
	case DFB:
		f32s, err := readDFB(f)
		if err != nil {
			return nil, err
		}
		return float32sToFloat64s(f32s), nil
	case DDB:
		f64s, err := readDDB(f)
		if err != nil {
			return nil, err
		}
		return f64s, nil
	default:
		return nil, errors.New("Unknown DataType")
	}
}

func readDSA(r io.Reader) ([]int16, error) {
	sc := bufio.NewScanner(r)
	var data []int16
	for sc.Scan() {
		v, err := strconv.ParseInt(sc.Text(), 10, 16)
		if err != nil {
			return nil, err
		}
		data = append(data, int16(v))
	}
	return data, nil
}

func readDFA(r io.Reader) ([]float32, error) {
	sc := bufio.NewScanner(r)
	var data []float32
	for sc.Scan() {
		v, err := strconv.ParseFloat(sc.Text(), 32)
		if err != nil {
			return nil, err
		}
		data = append(data, float32(v))
	}
	return data, nil
}

func readDDA(r io.Reader) ([]float64, error) {
	sc := bufio.NewScanner(r)
	var data []float64
	for sc.Scan() {
		v, err := strconv.ParseFloat(sc.Text(), 64)
		if err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

func readDSB(r io.Reader) ([]int16, error) {
	buf := make([]byte, BitLenShort)
	var data []int16
	for {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return data, err
		}
		v, err := bytesToInt16(buf)
		if err != nil {
			return data, err
		}
		data = append(data, v)
	}
}

func readDFB(r io.Reader) ([]float32, error) {
	buf := make([]byte, BitLenFloat)
	var data []float32
	for {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return data, err
		}
		v, err := bytesToFloat32(buf)
		if err != nil {
			return data, err
		}
		data = append(data, v)
	}
}

func readDDB(r io.Reader) ([]float64, error) {
	buf := make([]byte, BitLenDouble)
	var data []float64
	for {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return data, err
		}
		v, err := bytesToFloat64(buf)
		if err != nil {
			return data, err
		}
		data = append(data, v)
	}
}

func ext(path string) (string) {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}
