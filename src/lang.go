package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)
var i18n map[string]string

func loadLanguage(lang string ) error {
	filePath := filepath.Join("lang",lang + ".json")
	data , err := os.ReadFile(filePath)
	if err != nil{
		return fmt.Errorf("failed to parse language file: %v", err);
	}
	err = json.Unmarshal(data, &i18n)
	if err != nil {
		return fmt.Errorf("failed to parse language file: %v", err)
	}
	return nil
}

