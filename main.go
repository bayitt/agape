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

	translation, _ := getCurrentTranslation();

	sendAgapeEmail(translation.Language, translation.Text);
}

func getCurrentTranslation() (Translation, int) {
	directory, _ := os.Getwd();
	translationsPath := filepath.Join(directory, "translations.json");
	translationsContents, _ := os.ReadFile(translationsPath);

	var translationData *TranslationData;
	json.Unmarshal(translationsContents, &translationData);

	recordPath := filepath.Join(directory, "record.json");
	recordContents, _ := os.ReadFile(recordPath);

	var recordData *Record;
	json.Unmarshal(recordContents, &recordData);

	index := recordData.Index + 1;
	return translationData.Translations[index], index;
}

func sendAgapeEmail(language string, translation string) error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"));
	params := &resend.SendEmailRequest{
        To:      []string{os.Getenv("RECIPIENT_EMAIL")},
        From:    os.Getenv("MAIL_FROM"),
        Text:    translation,
        Subject: fmt.Sprintf("Today's language is %s", language),
    }

	_, err := client.Emails.Send(params);

	return err;
}