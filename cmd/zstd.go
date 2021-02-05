package cmd

import (
	"io"
)
import "github.com/klauspost/compress/zstd"

func CompressZstd(in io.Reader, out io.Writer) error {
	enc, err := zstd.NewWriter(out)
	if err != nil {
		return err
	}
	_, err = io.Copy(enc, in)
	if err != nil {
		_ = enc.Close()
		return err
	}
	return enc.Close()
}

func DecompressZstd(in io.Reader, out io.Writer) error {
	d, err := zstd.NewReader(in)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, d)
	if err != nil {
		return err
	}
	d.Close()
	return nil
}
