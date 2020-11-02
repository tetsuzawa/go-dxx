package dxx

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	BitLenShort  = 16
	BitLenFloat  = 32
	BitLenDouble = 64
)

var (
	ErrUnknownDataType = errors.New("unknown data type")
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
		return "unknown data type" // unreachable code
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
		return 0, ErrUnknownDataType
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
		return -1 // unreachable code
	}
}
func (dt DataType) ByteLen() int {
	return dt.BitLen() / 8
}

func Read(r io.Reader, dt DataType) ([]float64, error) {
	switch dt {
	case DSA:
		i16s, err := readDSA(r)
		if err != nil {
			return nil, err
		}
		return int16sToFloat64s(i16s), nil
	case DFA:
		f32s, err := readDFA(r)
		if err != nil {
			return nil, err
		}
		return float32sToFloat64s(f32s), nil
	case DDA:
		f64s, err := readDDA(r)
		if err != nil {
			return nil, err
		}
		return f64s, nil
	case DSB:
		i16s, err := readDSB(r)
		if err != nil {
			return nil, err
		}
		return int16sToFloat64s(i16s), nil
	case DFB:
		f32s, err := readDFB(r)
		if err != nil {
			return nil, err
		}
		return float32sToFloat64s(f32s), nil
	case DDB:
		f64s, err := readDDB(r)
		if err != nil {
			return nil, err
		}
		return f64s, nil
	default:
		return nil, errors.New(" ")
	}
}

// ReadFromFile reads .DXX file.
// This func determines the data type from the filename extension and reads that data.
// The return type is []float64 to make the data easier to handle.
func ReadFromFile(filename string) ([]float64, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	dt, err := StringToDataType(ext(filename))
	if err != nil {
		return nil, err
	}
	return Read(f, dt)
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

func Write(r io.Writer, dt DataType, data []float64) error {
	switch dt {
	case DSA:
		return writeDSA(r, float64sToInt16s(data))
	case DFA:
		return writeDFA(r, float64sToFloat32s(data))
	case DDA:
		return writeDDA(r, data)
	case DSB:
		return writeDSB(r, float64sToInt16s(data))
	case DFB:
		return writeDFB(r, float64sToFloat32s(data))
	case DDB:
		return writeDDB(r, data)
	default:
		return ErrUnknownDataType
	}
}

// Writes writes data to .DXX file.
// This func determines the data type from the filename extension and writes the data to the file.
// The return type is []float64 to make the data easier to handle.
func WriteToFile(filename string, data []float64) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	dt, err := StringToDataType(ext(filename))
	if err != nil {
		return err
	}
	return Write(f, dt, data)
}

func writeDSA(r io.Writer, data []int16) error {
	for _, v := range data {
		if _, err := fmt.Fprintf(r, "%d\n", v); err != nil {
			return err
		}
	}
	return nil
}

func writeDFA(r io.Writer, data []float32) error {
	for _, v := range data {
		if _, err := fmt.Fprintf(r, "%e\n", v); err != nil {
			return err
		}
	}
	return nil
}

func writeDDA(r io.Writer, data []float64) error {
	for _, v := range data {
		if _, err := fmt.Fprintf(r, "%e\n", v); err != nil {
			return err
		}
	}
	return nil
}

func writeDSB(r io.Writer, data []int16) error {
	for _, v := range data {
		buf, err := int16ToBytes(v)
		if err != nil {
			return err
		}
		_, err = r.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDFB(r io.Writer, data []float32) error {
	for _, v := range data {
		buf, err := float32ToBytes(v)
		if err != nil {
			return err
		}
		_, err = r.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDDB(r io.Writer, data []float64) error {
	for _, v := range data {
		buf, err := float64ToBytes(v)
		if err != nil {
			return err
		}
		_, err = r.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// ext returns the path of extension *without* dot.
// eg: ext(/path/to/file.aaa) -> aaa
func ext(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}
