package tcp

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
)

func encode(msg *Message) (*bytes.Buffer, error) {

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(msg); err != nil {
		return nil, err
	}

	return &buf, nil
}

func decode(data []byte, msg *Message) error {
	var buffer = bytes.NewBuffer(data)

	dec := gob.NewDecoder(buffer)
	if err := dec.Decode(msg); err != nil {
		return err
	}
	return nil
}

func pack(msg *Message) (*bytes.Buffer, error) {

	buf, err := encode(msg)
	if err != nil {
		return nil, err
	}
	n := int64(buf.Len())
	if n > maxMessageLength {
		return nil, errDataLenOvertakeMaxLen
	}

	buffer := bytes.NewBuffer(nil)
	binary.Write(buffer, binary.BigEndian, n)

	buffer.Write(buf.Bytes())
	return buffer, nil
}

func read(r io.Reader, msg *Message) error {
	var length int64
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	if length < 0 {
		return fmt.Errorf("read error, data length < 0, length=%d", length)
	}
	if length > maxMessageLength {
		return errDataLenOvertakeMaxLen
	}

	buffer := make([]byte, length)
	n, err := io.ReadFull(r, buffer)
	if err != nil {
		return fmt.Errorf("read error=%s", err.Error())
	}
	if int64(n) != length {
		return fmt.Errorf("read error")
	}

	return decode(buffer, msg)
}
