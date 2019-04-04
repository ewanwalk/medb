package file

import (
	"encoding/hex"
	"github.com/kalafut/imohash"
)

func Checksum(path string) (string, error) {
	sum, err := imohash.SumFile(path)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sum[:]), nil
}
