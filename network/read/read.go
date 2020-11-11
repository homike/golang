package connread

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

// proto
//  ----------------
// | size | content |

const (
	HEADER_SIZE = 2 // 2 byte head size
)

func Read_Bufio(bufReader *bufio.Reader) ([]byte, error) {
	var (
		err    error
		buffer []byte
	)

	buffer, err = bufReader.Peek(HEADER_SIZE)
	if err != nil {
		return nil, err
	}

	size, err := strconv.Atoi(string(buffer))
	if err != nil {
		return nil, err
	}

	totalSize := HEADER_SIZE + size
	buffer, err = bufReader.Peek(totalSize)
	if err != nil {
		return nil, err
	}

	data := make([]byte, totalSize-HEADER_SIZE)
	copy(data, buffer[HEADER_SIZE:])

	_, err = bufReader.Discard(totalSize)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type ByteBuf struct {
	buf  *bytes.Buffer
	size int // last packet length
}

func (b *ByteBuf) forward() error {
	header := b.buf.Next(HEADER_SIZE)
	size, err := strconv.Atoi(string(header))
	if err != nil {
		return err
	}
	b.size = size
	return nil
}

func (b *ByteBuf) Read_ByteBuf(data []byte) ([][]byte, error) {
	b.buf.Write(data)

	var err error
	// check length
	if b.buf.Len() < HEADER_SIZE {
		return nil, err
	}

	// first time
	if b.size < 0 {
		if err = b.forward(); err != nil {
			return nil, err
		}
	}

	arrData := [][]byte{}
	for b.size <= b.buf.Len() {
		arrData = append(arrData, b.buf.Next(b.size))

		// more packet
		if b.buf.Len() < HEADER_SIZE {
			b.size = -1
			break
		}

		if err = b.forward(); err != nil {
			return nil, err
		}
	}

	return arrData, nil
}

func Read_ReadFull(r io.Reader) ([]byte, error) {
	header := make([]byte, HEADER_SIZE)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(string(header), 10, 16)
	if err != nil {
		return nil, err
	}

	data := make([]byte, size)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
