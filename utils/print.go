package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrintStruct(str any) error {
	jsonData, err := json.Marshal(str)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	json.Indent(&out, jsonData, "", "\t")
	fmt.Println(out.String())
	return err
}
