package game

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/tylland/dashland/internal/ecs"
)

type TextReader struct {
	reader *bufio.Reader
}

func NewTextReader(r io.Reader) *TextReader {
	return &TextReader{reader: bufio.NewReader(r)}
}

func (r *TextReader) readLine() (string, error) {
	data, prefix, err := r.reader.ReadLine()

	if err != nil {
		return "", err
	}

	if prefix {
		return "", fmt.Errorf("line to long")
	}

	line := string(data)

	if idx := strings.Index(line, TokenComment); idx != -1 {
		line = line[:idx]
	}

	return strings.TrimSpace(line), nil
}

func (r *TextReader) readString() (string, error) {
	return r.readLine()
}

func (r *TextReader) ReadEntityID() (ecs.EntityID, error) {
	str, err := r.readLine()

	return ecs.EntityID(str), err
}

func (r *TextReader) ReadEntityType() (ecs.EntityType, error) {
	str, err := r.readLine()

	v, err := strconv.ParseInt(str, 10, 32)

	return ecs.EntityType(v), err
}

func (r *TextReader) ReadVersion() (int, error) {
	str, err := r.readLine()

	v, err := strconv.ParseInt(str, 10, 32)

	return int(v), err
}

func (r *TextReader) readNumValues() (bool, string, string, error) {
	var changed bool
	line, err := r.readLine()

	if err != nil {
		return changed, "", "", fmt.Errorf("invalid format")
	}

	if strings.HasPrefix(line, TokenChanged) {
		changed = true
		line = strings.TrimPrefix(line, TokenChanged)
	}

	nums := strings.Split(line, TokenNumberDelimiter)

	if len(nums) != 2 {
		return changed, "", "", fmt.Errorf("invalid format")
	}

	return changed, strings.TrimSpace(nums[0]), strings.TrimSpace(nums[1]), nil
}

func (r *TextReader) ReadLabel() (string, error) {
	line, err := r.readLine()

	if err != nil {
		return "", fmt.Errorf("read line: %w", err)
	}

	if strings.HasPrefix(line, TokenLabel) {
		line = strings.TrimPrefix(line, TokenLabel)
	} else {
		return "", fmt.Errorf("invalid format, %q is missing for label", TokenLabel)
	}

	return strings.TrimSpace(line), nil
}

func (r *TextReader) ReadIntValue() (int, bool, error) {
	changed, dec, _, err := r.readNumValues()

	v, err := strconv.ParseInt(dec, 10, 32)

	if err != nil {
		return -1, false, err
	}

	return int(v), changed, nil
}

func (r *TextReader) ReadFloatValue() (float32, bool, error) {
	changed, _, hex, err := r.readNumValues()

	u32, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return -1, false, err
	}

	v := math.Float32frombits(uint32(u32))

	return v, changed, err

	// v, err := strconv.ParseFloat(hex, 32)

	// if err != nil {
	// 	return -1, false, err
	// }

	// return float32(v), changed, err
}
