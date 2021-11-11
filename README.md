# music-exchange
hacky script for a secret-santa style music exchange

This script takes a JSON file as input and outputs markdown files containing instructions for each participant in an exchange.

## parsing the survey results

```
go run cli/main.go collect-survey-results --survey=survey-results.csv --previous-participants=people.json --email-address=1 --platforms=7 --ignore=0 > people-new.json
```

This takes the rows from `survey-results.csv` and ands some previous participant information from `poeple.json` and writes a new JSON file to `people-new.json`

Required Parameters:
* `--survey`: filepath of a CSV file containing the survey results. This file must contain a header row and at least 2 participants
* `--email-address`: index of the CSV column containing the participant email address
* `--platforms`: index of the CSV column containing the platform preferences


Optional Parameters:
* `--ignore`: comma-separated list of indexes of CSV columns to ignore
* `--previous-participants`: filepath of a JSON file used for a the most recent previous exchange
* `--separator`: the character that separates the list of platforms in the CSV, defaults to `;`

Columns other than the email address column, the platforms column, and the ignored columns will be treated as generic questions and answers and will ultimately be included in the instructions sent out to participants. There's an example survey CSV file in the root of this repo.

**Pro tip**: You can download a CSV of results from a Google Form.

## generating pairs for the exchange

There are 2 algorithms for generating the pairs. These mostly share the same parameters in common.

```
go run cli/main.go pair [bfScored|bfRandom] --participants=./people-test.json --instructions=./instructions-file-template.md
```

Required Parameters
* `--participants`: filepath of a JSON file containing the participant information
* `--instructions`: filepath of a golang template file containing an instructions template - either HTML (for email) or a format that can be written to local file.

Optional Parameters
* `--update`: pass as `true` to update the input JSON file with the latest pairing information

Parameters for sending instructions via email
* `--email` (**required**): set this flag to true to send instructions by email, otherwise defaults to writing instructions to local files
* `--smtp-host`: set this to set what OS Env Var this tool will reference for the SMTP host (default `MUSIC_EXCHANGE_SMTP_HOST`)
* `--smtp-port`: set this to set what OS Env Var this tool will reference for the SMTP port (default `MUSIC_EXCHANGE_SMTP_PORT`)
* `--smtp-username`: set this to set what OS Env Var this tool will reference for the SMTP username, which should be a full email address (default `MUSIC_EXCHANGE_SMTP_USERNAME`)
* `--smtp-password`: set this to set what OS Env Var this tool will reference for the SMTP password (default `MUSIC_EXCHANGE_SMTP_PASSWORD`)
* `--subject`: what email subject to use, defaults to "Music Exchange Assignment"
* `--recipient`: will send all instructions to this email address rather than the recipients' email addresses, useful for testing

Parameters for writing instructions to local files
* `--extension`: what file extension to use when writing local instructions files, defaults to "md"


#### brute force scored

This approach generates all possible compatible pairings and scores them based on how recently that participant has been pair previously with the same recipient.

Pairing sets with the best sum scores are preferred. Amongst those set, the best individual scores are preferred. Amongst those, the longest cycles are preferred. Amongst those, a random set is selected.

This approach is smarter and does not have a risk of running in an infinite loop.

#### brute force random

This approach shuffles the participants until a set of compatible pairs are found.

This basic run avoids someone being paired with themself.

If you want to avoid repeat recipients from previous runs, include the `--avoid` parameter. It's value will equal the number of subsequent previous pairings to avoid.

 _**Warning**: if there's no combination that satisfies all the requirements, then this can run in an infinite loop. So be prepared to kill the process._

## the participants JSON file

```
[{
    "EmailAddress": "jane123@gmail.com",
    "Skip": true,
    "LatestRecipients": ["jane456@gmail.com", "daisuke@watanabedaisuke.me"],
    "Platforms": [
      "Spotify", "Pandora"
    ],
    "Responses": [{
      "Question": "How would you describe your taste in music?",
      "Answer": "It like the classics - any from 40s big band hits to some funky disco."
    },{
      "Question": "Name 3-5 artists you enjoy.",
      "Answer": "Glen Miller, ABBA, Kool & The Gang"
    },{
      "Question": "What do you NOT enjoy?",
      "Answer": "Opera, Heavy Metal"
    },{
      "Question": "How close to your current tastes would you like this mix to be?",
      "Answer": "Pretty close, but I'm feeling a little adventurous."
    },{
      "Question": "How do you want this music to make you feel?",
      "Answer": "Happy, Up-Beat"
    }]
  }
  ...
]
```

* **EmailAddress**: the participant's email address - must be unique.
* **Skip** (optional): set to `true` to exclude this participant - this is the equivalent of "commenting out" this part of the JSON.
* **LatestRecipients** (optional): contains a list of previous participants that this person has created a playlist for, from most recent to least recent. This can be used to avoid pairing the same people over and over. These values should refer to `ID`s of other participants.
* **Platforms**: what music platforms this participant uses. This tool will only pair people with at least one platform in common (list must contain at least one identical string).
* **Responses**: list of general questions and answers. These are included in the instructions. Can be used for things like capturing someone's musical preferences.

There's also an example included in the root of this repo.

## the instructions template file

You have 2 options for distributing instructions to participants:
* You can have the CLI write local files containing instructions that you email manually. (This option can use any file format you want.)
* You can have the CLI send the emails containing instructions for you. (This can be tricky, depending on what level of access you have to an SMTP server. If you don't already know enough about SMPT to do this, I recommend the previous option.) (This option requires an HTML template.)

Either way you need an instructions template. Two example templates are provided in root of this repo. One MarkDown template for writing local files, and one HTML template for emailing. Feel free to use one of these or make your own. The only limitation is that it must be a valid golang template.