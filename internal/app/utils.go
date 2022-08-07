package app

import (
	"encoding/json"
	"io/ioutil"
)

func (app *Application) writeJSONToFile(dest string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err

	}
	_ = ioutil.WriteFile(dest, file, 0644)
	return nil
}
