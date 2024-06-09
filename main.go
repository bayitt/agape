package main;
import (
	"encoding/json"
	"os"
	"path/filepath"
	"fmt"
)


type Translation struct {
	Language string   `json:"language"` 
	Text	 string	  `json:"text"`
}

type TranslationData struct {
	Translations []Translation `json:"translations"`
}


func main() {
	directory, _ := os.Getwd();
	filePath := filepath.Join(directory, "translations.json");
	fileContents, _ := os.ReadFile(filePath);

	var translationData *TranslationData;
	json.Unmarshal(fileContents, &translationData);

	fmt.Println(translationData.Translations[12]);
}