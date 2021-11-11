# How to deploy to AWS lambda!

## Before you even get to the lambda:

### Create S3 Bucket

You'll need an s3 bucket used solely for this tool. Call it whatever you want.

### Create IAM Policy

You'll also need an IAM role that has access to that bucket as well as SES for sending email. Something like so:

**Trusted Entity**: AWS Service: lambda
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "S3Access",
            "Effect": "Allow",
            "Action": [
                "s3:*"
            ],
            "Resource": [
                "arn:aws:s3:::music-exchange",
                "arn:aws:s3:::music-exchange/*"
            ]
        },
        {
            "Sid": "SESAccess",
            "Effect": "Allow",
            "Action": [
                "ses:SendEmail"
            ],
            "Resource": [
                "arn:aws:ses:*:0123456789:identity/*"
            ]
        }
    ]
}
```

This policy will NOT work for you as is, unless your AWS account number is magically `0123456789` and your bucket name is `music-exchange`. Change as necessary.

## Now, Lambda Things

### Create lambda

When creating the lambda, use the IAM policy you created above.

### Configuring the lambda

The following environment variables are required:

* **EMAIL_SENDER**: This sender must be pre-validated in SES.
* **EMAIL_SUBJECT**: Subject of the instruction emails sent by the tool.
* **EMAIL_TEMPLATE_S3_KEY**: S3 object key of your instruction email body template.
* **SURVEY_EMAIL_ADDRESS_COLUMN**: Index of email address column in survey CSV (zero-indexed)
* **SURVEY_PLATFORMS_COLUMN**: Index of column containing the participant's music platforms (zero-indexed)
* **SURVEY_PLATFORMS_SEPARATOR**: Character that separates the platforms.

The following environment variables are optional:

* **EMAIL_TEST_RECIPIENT**: If you set this, all emails will be sent to this address rather than the participants' email addresses. Good for testing.
* **SURVEY_IGNORE_COLUMNS**: comma-separate list of column indexes to ignore (zero-indexed)

### Building

You should build like so:
```
GOOS=linux GOARCH=amd64 go build -o main lambda/main.go
zip lambda-function.zip main
rm main
```

If you run this from the root of this repo (not in the lambda directory), this will build the lambda as `lambda-function.zip`.

### Deploying

You should deploy like so, assuming you have the AWS CLI:
```
aws lambda update-function-code --function-name music-exchange --zip-file fileb://lambda-function.zip
```

This builds to an existing `music-exchange` lambda, but you can (of course) name it whatever you want.

### Before Running The First Time

#### Upload email template

You'll need to place your email template file in your s3 bucket. Do not prefix the key (filename) with `participants_`.

You can use the example provided in the root of this repo or make up your own.

#### Capture previous participants (optional)

Before running the first time, if you have already run the music exchange locally, you may have a previous participants file to upload. Upload that to you s3 bucket with the name `participants_1000000000.json`.

### Running

You'll need to place a CSV of your survey results in your S3 bucket. Do not prefix the key (filename) with `participants_`.

You can run this lambda using a test event directly in the browser. Here's the message to send:

```
{
  "SurveyS3ObjectKey": "Survey_2021-08-09.csv",
  "AllowRepeatRecipients": false
}
```

`SurveyS3ObjectKey` should be the S3 object key for your survey results. Do _not_ include the s3 bucket prefix. This CSV file must have a header row.

`AllowRepeatRecipients` is optional. You can leave it out. If this value is missing or `false`, this exchange will allow there to be repeat recipients from the immediately previous exchange. Set this value to `true` to avoid this.