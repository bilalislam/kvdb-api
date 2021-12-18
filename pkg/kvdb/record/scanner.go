package record

import (
	"bufio"
	"errors"
	"io"
)

type Scanner struct {
	*bufio.Scanner
}

func NewScanner(r io.Reader, maxScanTokenSize int) (*Scanner, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 4096)
	scanner.Buffer(buf, maxScanTokenSize+metaLength)
	scanner.Split(split)
	return &Scanner{scanner}, nil
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	r, err := FromBytes(data)
	if errors.Is(err, ErrInsufficientData) {
		return 0, nil, nil
	}

	if err != nil {
		return 0, nil, err
	}

	adv := r.Size()

	return adv, data[:adv], nil
}

func (r *Scanner) Record() *Record {
	data := r.Bytes()
	record, _ := FromBytes(data)
	return record
}
