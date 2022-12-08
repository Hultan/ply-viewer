package ply

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type FileFormat int

const (
	formatUnknown FileFormat = iota
	formatAscii
	formatBinaryLittleEndian
	formatBinaryBigEndian
)

var plyFormatBytes = []byte{112, 108, 121, 10}
var endHeaderBytes = []byte{101, 110, 100, 95, 104, 101, 97, 100, 101, 114}

type PLY struct {
	header           string
	data             []byte
	format           FileFormat
	numberOfVertexes int
	numberOfFaces    int
	Vertexes         []*Vertex
	Faces            []*Face
}

type Vertex struct {
	X, Y, Z float32
}

type Face struct {
	Indexes []int
}

func (p *PLY) Load(path string) error {
	bytes, err := loadFile(path)
	if err != nil {
		return err
	}

	i := p.findBytes(bytes, plyFormatBytes, 0)
	if i == -1 {
		return errors.New("this is not a correctly formatted PLY-file")
	}

	i = p.findBytes(bytes, endHeaderBytes, 0)
	if i == -1 {
		return errors.New("failed to find text 'end_header'")
	}

	h := p.findEndOfLine(bytes, i)
	if h == -1 {
		return errors.New("failed to find first new line character after 'end_header'")
	}

	p.header = string(bytes[:h])
	p.data = bytes[h+1:]

	err = p.parseHeader()
	if err != nil {
		return err
	}

	switch p.format {
	case formatAscii:
		err = p.parseDataAscii()
		if err != nil {
			return err
		}
	case formatBinaryLittleEndian:
		err = p.parseDataBinaryLittleEndian()
		if err != nil {
			return err
		}
	case formatBinaryBigEndian:
	}

	return nil
}

func loadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (p *PLY) findBytes(bytes, findBytes []byte, startAt int) int {
outer:
	for i := startAt; i < len(bytes); i++ {
		for e := range findBytes {
			if bytes[i+e] != findBytes[e] {
				continue outer
			}
		}
		return i
	}
	return -1
}

func (p *PLY) findEndOfLine(bytes []byte, start int) int {
	for i := start; i < len(bytes); i++ {
		if bytes[i] == 10 {
			return i
		}
	}
	return -1
}

func (p *PLY) parseHeader() error {
	r := strings.NewReader(p.header)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		err := p.parseHeaderLine(scanner.Text())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PLY) parseHeaderLine(text string) error {
	if strings.HasPrefix(text, "format") {
		return p.parseHeaderFormat(text)
	} else if strings.HasPrefix(text, "element vertex") {
		return p.parseHeaderVertex(text)
	} else if strings.HasPrefix(text, "element face") {
		return p.parseHeaderFace(text)
	}
	return nil
}

func (p *PLY) parseHeaderFormat(text string) error {
	if strings.HasPrefix(text, "format ascii") {
		p.format = formatAscii
	} else if strings.HasPrefix(text, "format binary_little_endian") {
		p.format = formatBinaryLittleEndian
	} else if strings.HasPrefix(text, "format binary_big_endian") {
		p.format = formatBinaryBigEndian
	} else {
		p.format = formatUnknown
		return errors.New("unknown PLY format")
	}
	return nil
}

func (p *PLY) parseHeaderVertex(text string) error {
	text = text[15:]
	text = p.getNumberPartOfString(text)
	v, err := strconv.Atoi(text)
	if err != nil {
		return errors.New("failed to parse number of vertexes")
	}
	p.numberOfVertexes = v
	return nil
}

func (p *PLY) parseHeaderFace(text string) error {
	text = text[13:]
	text = p.getNumberPartOfString(text)
	v, err := strconv.Atoi(text)
	if err != nil {
		return errors.New("failed to parse number of faces")
	}
	p.numberOfFaces = v
	return nil
}

func (p *PLY) getNumberPartOfString(text string) string {
	for i := range text {
		if text[i] >= 48 && text[i] <= 57 {
			continue
		}
		return text[:i]
	}
	return text
}

func (p *PLY) parseDataAscii() error {
	text := string(p.data)
	r := strings.NewReader(text)
	scanner := bufio.NewScanner(r)

	for i := 0; i < p.numberOfVertexes; i++ {
		scanner.Scan()
		v, err := p.parseDataVertexLine(scanner.Text())
		if err != nil {
			return err
		}
		if v != nil {
			p.Vertexes = append(p.Vertexes, v)
		}
	}

	for i := 0; i < p.numberOfFaces; i++ {
		scanner.Scan()
		f, err := p.parseDataFaceLine(scanner.Text())
		if err != nil {
			return err
		}
		if f != nil {
			p.Faces = append(p.Faces, f)
		}
	}

	return nil
}

func (p *PLY) parseDataVertexLine(text string) (*Vertex, error) {
	result := &Vertex{}

	f := strings.Fields(text)

	for i := 0; i < 3; i++ {
		s := f[i]
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, errors.New("failed to parse number of faces")
		}
		switch i {
		case 0:
			result.X = float32(v)
		case 1:
			result.Y = float32(v)
		case 2:
			result.Z = float32(v)
		}
	}

	return result, nil
}

func (p *PLY) parseDataFaceLine(text string) (*Face, error) {
	result := &Face{}

	f := strings.Fields(text)

	num, err := strconv.Atoi(f[0])
	if err != nil {
		return nil, errors.New("failed to parse face")
	}

	for i := 1; i <= num; i++ {
		s := f[i]

		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.New("failed to parse face")
		}
		result.Indexes = append(result.Indexes, v)
	}

	return result, nil
}

func (p *PLY) parseDataBinaryLittleEndian() error {
	delta := 0
	for i := 0; i < p.numberOfVertexes; i++ {
		f1 := p.float32FromBytesLE(delta + 0)
		f2 := p.float32FromBytesLE(delta + 4)
		f3 := p.float32FromBytesLE(delta + 8)
		v := &Vertex{f1, f2, f3}
		p.Vertexes = append(p.Vertexes, v)
		delta += 12
	}

	for i := 0; i < p.numberOfFaces; i++ {
		num := p.int8FromBytesLE(delta)
		delta += 1
		v := &Face{}
		for j := 0; j < int(num); j++ {
			index := p.int32FromBytesLE(delta)
			delta += 4
			v.Indexes = append(v.Indexes, index)
		}
		p.Faces = append(p.Faces, v)
	}

	return nil
}

func (p *PLY) float32FromBytesLE(id int) float32 {
	bytes := p.data[id : id+4]
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func (p *PLY) int32FromBytesLE(id int) int {
	bytes := p.data[id : id+4]
	bits := binary.LittleEndian.Uint32(bytes)
	return int(bits)
}

func (p *PLY) int8FromBytesLE(id int) int8 {
	return int8(p.data[id])
}
