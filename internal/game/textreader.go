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

func (r *TextReader) ConsumeLine() error {
	_, err := r.readLine()

	return err
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

func (r *TextReader) ReadEntityID() (ecs.EntityID, error) {
	str, err := r.readLine()

	return ecs.EntityID(str), err
}

func (r *TextReader) ReadEntityType() (ecs.EntityType, error) {
	str, err := r.readLine()

	v, err := strconv.ParseInt(str, 10, 32)

	return ecs.EntityType(v), err
}

func (r *TextReader) ReadBlockType() (int, error) {
	str, err := r.readLine()

	v, err := strconv.ParseInt(str, 10, 32)

	return int(v), err
}

func (r *TextReader) ReadVersion() (int, error) {
	str, err := r.readLine()

	v, err := strconv.ParseInt(str, 10, 32)

	return int(v), err
}

func (r *TextReader) readValue() (string, bool, error) {
	var changed bool
	line, err := r.readLine()

	if err != nil {
		return "", changed, fmt.Errorf("invalid format")
	}

	if strings.HasPrefix(line, TokenChanged) {
		changed = true
		line = strings.TrimPrefix(line, TokenChanged)
	}

	return strings.TrimSpace(line), changed, nil
}

func (r *TextReader) readNumValues() (bool, string, string, error) {
	str, changed, err := r.readValue()

	if err != nil {
		return changed, "", "", fmt.Errorf("read value: %w", err)
	}

	nums := strings.Split(str, TokenSeparator)

	if len(nums) != 2 {
		return changed, "", "", fmt.Errorf("invalid format")
	}

	return changed, strings.TrimSpace(nums[0]), strings.TrimSpace(nums[1]), nil
}

func (r *TextReader) ReadBool() (bool, bool, error) {
	str, changed, err := r.readValue()

	if err != nil {
		return false, changed, fmt.Errorf("read value: %w", err)
	}

	v, err := strconv.ParseBool(str)

	return v, changed, err
}

func (r *TextReader) ReadEnum() (uint32, bool, error) {
	str, changed, err := r.readValue()

	v, err := strconv.ParseUint(str, 10, 32)

	return uint32(v), changed, err
}

func (r *TextReader) ReadFloat() (float32, bool, error) {
	changed, _, hex, err := r.readNumValues()

	u32, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return -1, false, err
	}

	v := math.Float32frombits(uint32(u32))

	return v, changed, err
}

func (r *TextReader) ReadInt() (int, bool, error) {
	changed, dec, _, err := r.readNumValues()

	v, err := strconv.ParseInt(dec, 10, 32)

	if err != nil {
		return -1, false, err
	}

	return int(v), changed, nil
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

func (r *TextReader) ReadString() (string, bool, error) {
	return r.readValue()
}
