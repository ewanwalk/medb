package file

import (
	"io"
	"os"
)

func Rename(src, dest string) error {

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	stat, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, stat.Mode()&os.ModePerm)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Close()
	if err != nil {
		return err
	}

	err = in.Close()
	if err != nil {
		return err
	}

	return os.Remove(src)
}
