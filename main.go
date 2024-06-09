package main;
import (
	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
	"encoding/json"
	"os"
	"path/filepath"
	"fmt"
	"time"
)


type Translation struct {
	Language string   `json:"language"` 
	Text	 string	  `json:"text"`
}

type TranslationData struct {
	Translations []Translation `json:"translations"`
}

type Record struct {
	Index  		int  		`json:"index"`
	Timestamp   *time.Time  `json:"timestamp"`
	Language	*string		`json:"language"` 
}


func main() {
	godotenv.Load();

	directory, _ := os.Getwd();
	filePath := filepath.Join(directory, "translations.json");
	fileContents, _ := os.ReadFile(filePath);

	var translationData *TranslationData;
	json.Unmarshal(fileContents, &translationData);

	fmt.Println(translationData.Translations[12]);
}