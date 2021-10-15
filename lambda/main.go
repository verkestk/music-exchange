package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/verkestk/music-exchange/src/lambda_exchange"
)

// MusicExchangeEvent lambda event
type MusicExchangeEvent struct {
	SurveyS3ObjectKey string
}

// HandleRequest required by AWS lambda
func HandleRequest(ctx context.Context, event MusicExchangeEvent) error {
	if event.SurveyS3ObjectKey == "" {
		return fmt.Errorf("SurveyS3ObjectKey required")
	}

	config := &lambda_exchange.Config{
		EmailSender:                 os.Getenv("EMAIL_SENDER"),
		EmailSubject:                os.Getenv("EMAIL_SUBJECT"),
		EmailTemplateS3Key:          os.Getenv("EMAIL_TEMPLATE_S3_KEY"),
		EmailTestRecipient:          os.Getenv("EMAIL_TEST_RECIPIENT"),
		S3Bucket:                    os.Getenv("S3_BUCKET"),
		SurveyEmailAddressColumnStr: os.Getenv("SURVEY_EMAIL_ADDRESS_COLUMN"),
		SurveyIgnoreColumnsStr:      os.Getenv("SURVEY_IGNORE_COLUMNS"),
		SurveyPlatformsColumnStr:    os.Getenv("SURVEY_PLATFORMS_COLUMN"),
		SurveyPlatformsSeparator:    os.Getenv("SURVEY_PLATFORMS_SEPARATOR"),
		SurveyS3ObjectKey:           event.SurveyS3ObjectKey,
	}

	return lambda_exchange.Do(config)
}

func main() {
	lambda.Start(HandleRequest)
}
