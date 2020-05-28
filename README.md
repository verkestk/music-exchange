# music-exchange
hacky script for a secret-santa style music exchange


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

Warning, if there's no combination that satisfies avoid=1 or avoid=2, then this will run in an infinite loop.
