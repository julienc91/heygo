package tools

// From https://github.com/oz/osdb

import (
	"encoding/binary"
	"os"
	"strconv"
)

func Hash(filename string) (string, uint64, error) {

	file, err := os.Open(filename)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", 0, err
	}

	var size = uint64(stat.Size())
	var hash = size
	var b = make([]byte, 131072)

	_, err = file.Read(b[:65536])
	if err != nil {
		return "", 0, err
	}

	_, err = file.Seek(-65536, 2)
	if err != nil {
		return "", 0, err
	}

	_, err = file.Read(b[65536:])
	if err != nil {
		return "", 0, err
	}

	for i := 0; i < 16384; i++ {
		hash += binary.LittleEndian.Uint64(b[i*8 : i*8+8])
	}

	return strconv.FormatUint(hash, 16), size, nil
}
