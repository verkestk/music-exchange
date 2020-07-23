# music-exchange
hacky script for a secret-santa style music exchange

This script takes a JSON file as input and outputs markdown files containing instructions for each participant in an exchange.

## very basic instructions

This basic run avoids someone being paired with themself.
```
go run main.go --avoid=0 --filepath=./people.json
```

To avoid getting the same recipient you go last time, run thus:
```
go run main.go --avoid=1 --filepath=./people.json
```

And to avoid the last 2 recipients:
```
go run main.go --avoid=2 --filepath=./people.json
```

And so on. Warning, if there's no combination that satisfies all the requirements, then this can run in an infinite loop. So be prepared to kill the process.

## the JSON file

```
[{
    "Name": "Jane",
    "ID": "jane123",
    "Skip": true,
    "LatestRecipients": ["jane456", "daisuke"],
    "Platforms": [
      "Spotify", "Pandora"
    ],
    "Responses": [{
      "Question": "How would you describe your taste in music?",
      "Answer": ""
    },{
      "Question": "Name 3-5 artists you enjoy.",
      "Answer": ""
    },{
      "Question": "What do you NOT enjoy?",
      "Answer": ""
    },{
      "Question": "How close to your current tastes would you like this mix to be?",
      "Answer": ""
    },{
      "Question": "How do you want this music to make you feel?",
      "Answer": ""
    }]
  }
  ...
]
```

* **Name**: the participant's name - does _not_ need to be unique.
* **ID**: any unique string that identifies the participant - like a username. Should contain only filename-legal characters.
* **Skip** (optional): set to `true` to exclude this participant - this is the equivalent of "commenting out" this part of the JSON.
* **LatestRecipients** (optional): contains a list of previous participants that this person has created a playlist for, from most recent to least recent. This can be used to avoid pairing the same people over and over. These values should refer to `ID`s of other participants.
* **Platforms**: what music platforms this participant uses. This tool will only pair people with at least one platform in common (list must contain at least one identical string).
* **Responses**: list of general questions and answers. These are included in the instructions. Can be used for things like capturing someone's musical preferences.

## the markdown file Instructions

Running the script results in the creation of individual markdown files containing instructions for each participant. Don't peak if you want it to be a surprise!

These instructions are generated with a golang template. Modify the template if you want to change the instructions.

Share the instructions however you want - like in an email! And have fun!
