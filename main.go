package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Translation struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

type TranslationData struct {
	Translations []Translation `json:"translations"`
}

type Record struct {
	Index    int     `json:"index"`
	Date     *string `json:"date"`
	Language *string `json:"language"`
}

func main() {
	godotenv.Load()

	router := gin.Default()
	router.POST("/send", sendRouteHandler)
	router.Run("localhost:8080")
}

func sendRouteHandler(context *gin.Context) {
	directory, _ := os.Getwd()
	templatePath := filepath.Join(directory, "templates", "agape.html")
	templateString, _ := os.ReadFile(templatePath)
	templateObject, _ := template.New("Agape Template").Parse(string(templateString))

	var parsedTemplateString bytes.Buffer

	translation, index, err := getCurrentTranslation()

	if err != nil && err.Error() == "DAY_TRANSLATION_ALREADY_SENT" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Translation for today has been sent already"})
		return
	}

	if err != nil && err.Error() == "TRANSLATIONS_EXHAUSTED" {
		sendEmail(os.Getenv("ADMIN_EMAIL"), "All translations have been sent", "There are no more translations to be sent. Terminate the program.")
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "All translations have been sent!"})
		return
	}

	templateObject.Execute(&parsedTemplateString, map[string]string{"Translation": translation.Text, "Language": translation.Language})

	mailErr := sendEmail(os.Getenv("RECIPIENT_EMAIL"), fmt.Sprintf("Today's language is %s!", translation.Language), parsedTemplateString.String())

	if mailErr == nil {
		sendEmail(os.Getenv("ADMIN_EMAIL"), fmt.Sprintf("Today's Language - %s", translation.Language), fmt.Sprintf("Today's Email has been sent - %s", translation.Language))
		updateRecord(translation.Language, index)
		context.IndentedJSON(http.StatusOK, gin.H{"message": "Agape email sent successfully"})
		return
	}

	sendEmail(os.Getenv("ADMIN_EMAIL"), fmt.Sprintf("Today's Language - %s", translation.Language), fmt.Sprintf("Today's Email (%s) was not sent. There was an issue", translation.Language))
	context.IndentedJSON(http.StatusOK, gin.H{"message": "There was a problem sending the Agape email"})
}

func getCurrentTranslation() (Translation, int, error) {
	directory, _ := os.Getwd()
	translationsPath := filepath.Join(directory, "translations.json")
	translationsContents, _ := os.ReadFile(translationsPath)

	var translationData *TranslationData
	json.Unmarshal(translationsContents, &translationData)

	recordPath := filepath.Join(directory, "record.json")
	recordContents, _ := os.ReadFile(recordPath)

	var recordData *Record
	json.Unmarshal(recordContents, &recordData)

	index := recordData.Index + 1

	date := (time.Now()).Format("2006-01-02")

	if recordData.Date != nil && *recordData.Date == date {
		return translationData.Translations[0], index, errors.New("DAY_TRANSLATION_ALREADY_SENT")
	}

	if index > len(translationData.Translations)-1 {
		return translationData.Translations[0], index, errors.New("TRANSLATIONS_EXHAUSTED")
	}

	return translationData.Translations[index], index, nil
}

func sendEmail(recipient string, subject string, html string) error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	params := &resend.SendEmailRequest{
		To:      []string{recipient},
		From:    os.Getenv("MAIL_FROM"),
		Html:    html,
		Subject: subject,
	}

	_, err := client.Emails.Send(params)

	return err
}

func updateRecord(language string, index int) {
	date := (time.Now()).Format("2006-01-02")
	record := &Record{Index: index, Language: &language, Date: &date}
	recordJson, _ := json.Marshal(record)

	directory, _ := os.Getwd()
	recordPath := filepath.Join(directory, "record.json")
	os.WriteFile(recordPath, recordJson, 0777)
}
