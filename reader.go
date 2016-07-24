package unity

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
)

// DataReader Unityファイル用バイナリリーダー
type DataReader struct {
	raw    *[]byte
	buffer *bytes.Reader
}

// NewDataReader new DataReader instance
func NewDataReader(b []byte) (*DataReader, error) {
	buf := bytes.NewReader(b)
	return &DataReader{
		&b,
		buf,
	}, nil
}

// NewDataReaderFromFilePath new DataReader instance
func NewDataReaderFromFilePath(path string) (*DataReader, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewDataReader(b)
}

// ReadByte read 1 byte
func (data *DataReader) ReadByte() (byte, error) {
	return data.buffer.ReadByte()
}

func selectByteOrder(isLittleEndian bool) binary.ByteOrder {
	var endian binary.ByteOrder
	if isLittleEndian {
		endian = binary.LittleEndian
	} else {
		endian = binary.BigEndian
	}
	return endian
}

// ReadChar read short
func (data *DataReader) ReadChar(isLittleEndian bool) (int8, error) {
	var i int8
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// ReadUchar read short
func (data *DataReader) ReadUchar(isLittleEndian bool) (uint8, error) {
	var i uint8
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// ReadShort read short
func (data *DataReader) ReadShort(isLittleEndian bool) (int16, error) {
	var i int16
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// ReadUshort read unsigned short
func (data *DataReader) ReadUshort(isLittleEndian bool) (uint16, error) {
	var i uint16
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// ReadInt read int
func (data *DataReader) ReadInt(isLittleEndian bool) (int32, error) {
	var i int32
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// ReadUint read unsigned int
func (data *DataReader) ReadUint(isLittleEndian bool) (uint32, error) {
	var i uint32
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// ReadLong read long
func (data *DataReader) ReadLong(isLittleEndian bool) (int64, error) {
	var i int64
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// ReadUlong read unsigned long
func (data *DataReader) ReadUlong(isLittleEndian bool) (uint64, error) {
	var i uint64
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (data *DataReader) ReadBytes(size int, isLittleEndian bool) ([]byte, error) {
	b := make([]byte, size)
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// ReNew
func (data *DataReader) ReNew(size int, isLittleEndian bool) (*DataReader, error) {
	b := make([]byte, size)
	err := binary.Read(data.buffer, selectByteOrder(isLittleEndian), &b)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(b)
	return &DataReader{
		&b,
		buf,
	}, nil
}

// ReadStringNull 0x00に到達するまで文字列を読み込む
func (data *DataReader) ReadStringNull(limit int) (string, error) {
	b := []byte{}
	for i := 0; i < limit; i++ {
		c, err := data.ReadByte()
		if err != nil {
			return "", err
		}
		if c == 0 {
			break
		}
		b = append(b, c)
	}
	return string(b), nil
}

// Seek implements the io.Seeker interface.
func (data *DataReader) Seek(offset int64, whence int) (int64, error) {
	return data.buffer.Seek(offset, whence)
}

// Len Return the current stream position.
func (data *DataReader) Len() int {
	return data.buffer.Len()
}

func (data *DataReader) Align() error {
	size := len(*data.raw)
	oldPos := size - data.buffer.Len()
	newPos := (oldPos + 3) & -4
	if newPos > oldPos {
		_, err := data.buffer.Seek(int64(newPos-oldPos), os.SEEK_CUR)
		if err != nil {
			return err
		}
	}
	return nil
}
