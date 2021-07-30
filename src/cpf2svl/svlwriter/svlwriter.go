package svlwriter

import (
	"bufio"
	"encoding/binary"
	"io"
)

// SVLWriter is moe binary file writer type
type SVLWriter struct {
	writer *bufio.Writer
}

// NewSVLWriter creates new SVLWriter type
func NewSVLWriter(writer io.Writer) SVLWriter {
	return SVLWriter{writer: bufio.NewWriter(writer)}
}

// Flush is flushing buffer
func (w *SVLWriter) Flush() error {
	return w.writer.Flush()
}

func (w *SVLWriter) writeSize(v uint32) error {
	return binary.Write(w.writer, binary.BigEndian, v)
}

// WriteToken writes MOE token
func (w *SVLWriter) WriteToken(tokens []string) error {
	if err := w.writer.WriteByte(4); err != nil {
		return err
	}
	if err := w.writeSize(uint32(len(tokens))); err != nil {
		return err
	}
	for _, token := range tokens {
		if err := w.writeSize(uint32(len(token))); err != nil {
			return err
		}
		if _, err := w.writer.WriteString(token); err != nil {
			return err
		}
	}
	return nil
}

// WriteInt writes MOE integer
func (w *SVLWriter) WriteInt(vals []int) error {
	if err := w.writer.WriteByte(2); err != nil {
		return err
	}
	if err := w.writeSize(uint32(len(vals))); err != nil {
		return err
	}
	for _, val := range vals {
		if err := binary.Write(w.writer, binary.BigEndian, int32(val)); err != nil {
			return err
		}
	}
	return nil
}

// WriteFloat writes MOE float
func (w *SVLWriter) WriteFloat(vals []float64) error {
	if err := w.writer.WriteByte(3); err != nil {
		return err
	}
	if err := w.writeSize(uint32(len(vals))); err != nil {
		return err
	}
	for _, val := range vals {
		if err := binary.Write(w.writer, binary.BigEndian, val); err != nil {
			return err
		}
	}
	return nil
}
