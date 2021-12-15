package record

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

const (
	kindByteSize   = 1
	crcLen         = 4
	keyLenByteSize = 4
	valLenByteSize = 4
	metaLength     = kindByteSize + crcLen + keyLenByteSize + valLenByteSize
)

const (
	valueKind = iota
	tombstoneKind
)

var ErrInsufficientData = errors.New("could not parse bytes")
var ErrCorruptData = errors.New("the record has been corrupted")

type Record struct {
	kind  byte
	key   string
	value []byte
}

//NewValue returns a new record of value kind
func NewValue(key string, value []byte) *Record {
	return &Record{
		kind:  valueKind,
		key:   key,
		value: value,
	}
}

func NewTombstone(key string) *Record {
	return &Record{
		kind:  tombstoneKind,
		key:   key,
		value: []byte{},
	}
}

func (r *Record) Key() string {
	return r.key
}

func (r *Record) Value() []byte {
	return r.value
}

func (r *Record) IsTombstone() bool {
	return r.kind == tombstoneKind
}

func (r *Record) ToBytes() []byte {
	keyBytes := []byte(r.key)

	keyLen := make([]byte, keyLenByteSize)
	binary.BigEndian.PutUint32(keyLen, uint32(len(keyBytes)))

	valLen := make([]byte, valLenByteSize)
	binary.BigEndian.PutUint32(valLen, uint32(len(r.value)))

	var data []byte
	crc := crc32.NewIEEE()
	for _, v := range [][]byte{{r.kind}, keyLen, valLen, []byte(r.key), r.value} {
		data = append(data, v...)
		_, _ = crc.Write(v)
	}

	crcData := make([]byte, crcLen)
	binary.BigEndian.PutUint32(crcData, crc.Sum32())
	return append(crcData, data...)
}

// FromBytes deserialize []byte into a record. If the data cannot be
// deserialized a wrapped ErrParse error will be returned.
func FromBytes(data []byte) (*Record, error) {
	if len(data) < metaLength {
		return nil, ErrInsufficientData
	}

	keyLenStart := crcLen + kindByteSize
	klb := data[keyLenStart : keyLenStart+keyLenByteSize]
	vlb := data[keyLenStart+keyLenByteSize : keyLenStart+keyLenByteSize+valLenByteSize]

	crc := binary.BigEndian.Uint32(data[:4])
	keyLen := int(binary.BigEndian.Uint32(klb))
	valLen := int(binary.BigEndian.Uint32(vlb))

	if len(data) < metaLength+keyLen+valLen {
		return nil, ErrInsufficientData
	}

	keyStartIdx := metaLength
	valStartIdx := keyStartIdx + keyLen

	kind := data[crcLen]
	key := make([]byte, keyLen)
	val := make([]byte, valLen)
	copy(key, data[keyStartIdx:valStartIdx])
	copy(val, data[valStartIdx:valStartIdx+valLen])

	check := crc32.NewIEEE()
	_, _ = check.Write(data[4 : metaLength+keyLen+valLen])
	if check.Sum32() != crc {
		return nil, ErrCorruptData
	}

	return &Record{kind: kind, key: string(key), value: val}, nil
}

// Size returns the serialized byte size
func (r *Record) Size() int {
	return crcLen + kindByteSize + keyLenByteSize + valLenByteSize + len(r.key) + len(r.value)
}

func (r *Record) Write(w io.Writer) (int, error) {
	data := r.ToBytes()
	return w.Write(data)
}
