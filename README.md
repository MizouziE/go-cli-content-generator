# Content Generator

## Getting started

To start this CLI in it's first iteration:

```sh
go run main.go <optional-path-to-config.yaml> <optional-path-to-data.csv>
```

## Input your OpenAI credentials

On your first time using this application, you will be prompted for you API key which you can obtain from [Open AI's website](https://platform.openai.com/account/api-keys) if you do not already have it to hand.

This will be saved to your own local `.env` file for the application to use each time it is started.

## Add your files

Optionally, you may add a config yaml file path as a first argument on initiation. This config yaml can also already include a field `data:` with the relative path to the csv data file.

The path to the csv data file may also be passed as the second argument to the initiation.

It is possible to use a csv with any number of columns. An example is provided for your use:

```sh
storage/example.csv
```

## Promts generated

The prompt structure is to be provided in the configuration under a `prompts:` field.

An example prompt base sentence is:

>prompt:
>- Write me a 150 word story about a {{ .mood }} {{ .animal }} that is {{ .action }}
>- Write a poem about a {{ .animal }} that you met while {{ .action }}

Each `{{ .curly-braced-word }}` corresponds to the column of the csv. So based on this, it is best to have a heading row and csv layout that go something like:

```csv
mood,animal,action
happy,duck,swimming
upset,bird,flying
sombre,dog,eating
```

Take a look at [./storage/example.csv](./storage/example.csv) to see examples.

## Files generated

The prompts will be run in the order they are provided for each row.

Each prompts output is concatenated into a single file **for each row**.

The output from each individual row will be written to a separate file inside a directory created and timestamped when the csv file is provided. The filename for each row's output will be `row-output-#` where # = the row number, so you may refer back to the supplied csv to see what prompt resulted in which section.

***

## Authors and acknowledgment

Sam and Madalin

## License

Built by Vlah on behalf of Intelligiants

## Project status

Proof of concept
