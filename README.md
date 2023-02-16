# Content Generator

## Getting started

To start this CLI in it's first iteration:

```sh
go run main.go
```

## Input your OpenAI credentials

On your first time using this application, you will be prompted for you API key which you can obtain from [Open AI's website](https://platform.openai.com/account/api-keys) if you do not already have it to hand.

This will be saved to your own local `.env` file for the application to use each time it is started.

## Add your files

It is possible to use any 3 column csv, just provide the relative path when prompted. An example is provided for your use:

```sh
storage/example.csv
```

## Promts generated

Currently, the prompting format is a little hard-coded for demonstration purposes. Understanding it will allow you to make more effective input csv files.

The prompt base sentence is:

> "Write me a 150 word story about a <\column-#1> <\column-#2> that is <\column-#3>"

Each `<\column-#>` corresponds to the numbered column of the csv. So based on this, it is best to have rows that go something like `"adjective","noun","verb OR desciption of an action"`. Take a look at [./storage/example.csv](./storage/example.csv) to see examples.

>! Be aware that there is currently no heading row, the first row is just the first row.

## Files generated

The output from each individual prompt will be written to a separate file inside a directory created and timestamped when the csv file is provided. The filename for each row's output will be `story-#` where # = the row number, so you may refer back to the supplied csv to see what prompt resulted in which story.

***

## Authors and acknowledgment

Sam and Madalin

## License

Built by Vlah on behalf of Intelligiants

## Project status

Proof of concept
