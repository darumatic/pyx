package cmd

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func DecompressTar(in io.Reader, dir string) error {
	tarReader := tar.NewReader(in)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("DecompressTar: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(dir, header.Name), 0755); err != nil {
				fmt.Printf("DecompressTar: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			file := filepath.Join(dir, header.Name)
			_ = os.MkdirAll(filepath.Dir(file), os.ModePerm)
			outFile, err := os.Create(file)
			if err != nil {
				fmt.Printf("DecompressTar: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				fmt.Printf("DecompressTar: Copy() failed: %s", err.Error())
			}
			_ = outFile.Close()
		default:
		}
	}
	return nil
}
