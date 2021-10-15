package lambda_exchange

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/ses"

	"github.com/verkestk/music-exchange/operation"
	"github.com/verkestk/music-exchange/src/email"
)

// Config contains all the input necessary to run the music exchange lambda
type Config struct {
	EmailSender                 string
	EmailSubject                string
	EmailTemplateS3Key          string
	EmailTestRecipient          string
	S3Bucket                    string
	SurveyEmailAddressColumnStr string
	SurveyIgnoreColumnsStr      string
	SurveyPlatformsColumnStr    string
	SurveyPlatformsSeparator    string
	SurveyS3ObjectKey           string
}

// Do does the music exchange
func Do(exchangeConfig *Config) error {

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	sesClient := ses.NewFromConfig(cfg)

	// step 1: get previous participants JSON
	previousParticipantsObject, err := getMostRecentParticipantsObject(ctx, s3Client, exchangeConfig.S3Bucket)
	if err != nil {
		return fmt.Errorf("error getting most recent participants from s3: %w", err)
	}

	previousParticipantsJSON := ""
	if previousParticipantsObject != nil {
		previousParticipantsJSON, err = getObjectContent(ctx, s3Client, exchangeConfig.S3Bucket, *previousParticipantsObject.Key)
		if err != nil {
			return fmt.Errorf("error reading previous participant JSON object")
		}
	}

	// step 2: get survey results CSV
	surveyResultsCSV, err := getObjectContent(ctx, s3Client, exchangeConfig.S3Bucket, exchangeConfig.SurveyS3ObjectKey)
	if err != nil {
		return fmt.Errorf("error loading the survey results from object %s: %w", exchangeConfig.SurveyS3ObjectKey, err)
	}

	// step 3: get email instructions template
	emailTemplateStr, err := getObjectContent(ctx, s3Client, exchangeConfig.S3Bucket, exchangeConfig.EmailTemplateS3Key)
	if err != nil {
		return fmt.Errorf("error loading the email template from object %s: %w", exchangeConfig.EmailTemplateS3Key, err)
	}

	// step 4: do the collect operation
	emailAddressColumn, err := strconv.Atoi(exchangeConfig.SurveyEmailAddressColumnStr)
	if err != nil {
		return fmt.Errorf("invalid SURVEY_EMAIL_ADDRESS_COLUMN: %s", exchangeConfig.SurveyEmailAddressColumnStr)
	}
	platformsColumn, err := strconv.Atoi(exchangeConfig.SurveyPlatformsColumnStr)
	if err != nil {
		return fmt.Errorf("invalid SURVEY_PLATFORMS_COLUMN: %s", exchangeConfig.SurveyPlatformsColumnStr)
	}
	collectConfig := &operation.CollectConfig{
		SurveyCSV:                surveyResultsCSV,
		PreviousParticipantsJSON: previousParticipantsJSON,
		EmailAddressColumn:       emailAddressColumn,
		PlatformsColumn:          platformsColumn,
		IgnoreColumnsStr:         exchangeConfig.SurveyIgnoreColumnsStr,
		PlatformsSeparator:       exchangeConfig.SurveyPlatformsSeparator,
	}
	err = collectConfig.Prepare()
	if err != nil {
		return err
	}
	newJSON, err := operation.DoCollect(collectConfig)
	if err != nil {
		return err
	}

	// step 5: write the new JSON
	newJSONObjectKey := fmt.Sprintf("participants_%d", time.Now().Unix())
	err = putObjectContent(ctx, s3Client, exchangeConfig.S3Bucket, newJSONObjectKey, newJSON)
	if err != nil {
		return fmt.Errorf("error writing new participants JSON file")
	}

	// step 6: do the pair operation
	pairConfig := &operation.PairConfig{
		ParticipantsJSON:        newJSON,
		InstructionsTemplateStr: emailTemplateStr,
		UpdateParticipantsFile:  true,
		Algorithm:               operation.BFScored,
		EmailInstructions:       true,
		EmailSubject:            exchangeConfig.EmailSubject,
		EmailTestRecipient:      exchangeConfig.EmailTestRecipient,
		EmailSender:             email.GetSESSender(ctx, sesClient, exchangeConfig.EmailSender),
	}

	err = pairConfig.Prepare()
	if err != nil {
		return err
	}

	// step 4: run the pairing algorithm, send email, and update JSON in s3
	return operation.DoPair(pairConfig)
}

// gets the latest of executable for this build's OS/Architecture
func getMostRecentParticipantsObject(ctx context.Context, s3Client *s3.Client, bucket string) (object *types.Object, err error) {
	params := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String("particiants_"),
		MaxKeys: 9999,
	}

	output, err := s3Client.ListObjectsV2(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(output.Contents) == 0 {
		return nil, nil
	}

	sort.Slice(output.Contents, func(i, j int) bool {
		return !(*output.Contents[i].Key < *output.Contents[j].Key)
	})

	return &output.Contents[0], nil
}

func getObjectContent(ctx context.Context, s3Client *s3.Client, bucket, key string) (string, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	output, err := s3Client.GetObject(ctx, params)
	if err != nil {
		return "", err
	}

	// output, err := client.GetObject(ctx, params)
	// if err != nil {
	// 	return "", err
	// }
	defer output.Body.Close()
	bytes := make([]byte, output.ContentLength)
	output.Body.Read(bytes)

	return strings.Trim(string(bytes), " \n"), nil

	// if output.ContentLength == 0 {
	// 	return "", nil
	// }
	//
	// defer output.Body.Close()
	// bytes := make([]byte, output.ContentLength)
	// _, err = output.Body.Read(bytes)
	// if err != nil {
	// 	return "", err
	// }
	//
	// return string(bytes), nil
}

func putObjectContent(ctx context.Context, s3Client *s3.Client, bucket, key, content string) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	}

	_, err := s3Client.PutObject(ctx, params)
	return err
}
