package utils

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"

	"github.com/Tnze/go-mc/nbt"
)

func readCompressedDatFile(dat_path string) (content []byte, err error) {
	file, err := os.Open(dat_path)
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	content, err = io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func TryExtractDatNBT(dat_path string) (NBT map[string]interface{}, err error) {
	content, err := readCompressedDatFile(dat_path)
	if err != nil {
		return nil, err
	}
	err = nbt.Unmarshal(content, &NBT)
	if err != nil {
		return nil, err
	}
	return NBT, err
}

func TrySaveNBTToDat(NBT map[string]interface{}, dat_path string) bool {
	content, err := nbt.Marshal(NBT)
	if err != nil {
		return false
	}
	file, err := os.OpenFile(dat_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false
	}
	writer := gzip.NewWriter(file)
	writer.Write(content)
	writer.Close()
	return true
}
