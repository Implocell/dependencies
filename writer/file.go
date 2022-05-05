package writer

import (
	"fmt"
	"os"
)

func writeToFile(data []byte, file *os.File) error {
	n, err := file.Write(data)

	if n == 0 {
		return fmt.Errorf("no bytes were written to file %s", file.Name())
	}

	if err != nil {
		return err
	}

	return nil
}
