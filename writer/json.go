package writer

import (
	"encoding/json"
	"os"
)

func JSONFile(v interface{}, file *os.File) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = writeToFile(res, file)

	if err != nil {
		return err
	}

	return nil
}
