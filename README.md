# go-fetch-weather
Weather CLI in Golang

## Usage

Sign up for a Meteostat API key in their developer portal: https://dev.meteostat.net/

The key can be used in a .env file or environment variable entered in the command:

Build a binary with go build:

```bash
go build go-fetch-weather.go
```

Then run the cli app:

```bash
METEOSTAT_API_KEY=1234 ./go-fetch-weather
```
