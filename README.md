# go-openholidays

A golang client library for the [Open Holidays Api](https://www.openholidaysapi.org)


## How it's made ?

1. Install oapi-codegen : `go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen`
1. Generate types: `oapi-codegen -generate types -package openholidays -o types.go https://openholidaysapi.org/swagger/v1/swagger.json`
1. Generate client: `oapi-codegen -generate client -package openholidays -o client.go https://openholidaysapi.org/swagger/v1/swagger.json`

## Usage

Example: Get the public holidays for Canton de Vaud in 2025

```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"openholidays"

	"github.com/oapi-codegen/runtime/types"
)

func main() {

	ctx := context.Background()

	client, _ := openholidays.NewClient("https://openholidaysapi.org")

	from, _ := time.Parse("2006-01-02", "2025-01-01")
	to, _ := time.Parse("2006-01-02", "2025-12-31")

	validFrom := types.Date{
		Time: from,
	}

	validTo := types.Date{
		Time: to,
	}

	params := openholidays.GetPublicHolidaysParams{
		CountryIsoCode:  "CH",
		ValidFrom:       validFrom,
		ValidTo:         validTo,
		LanguageIsoCode: func(s string) *string { return &s }("FR"),
		SubdivisionCode: func(s string) *string { return &s }("CH-VD"),
	}

	// Example for GetHolidays function
	holidays, err := client.GetPublicHolidays(ctx, &params)
	if err != nil {
		log.Fatalf("Error getting holidays: %v", err)
	}

	var holidayResponses []openholidays.HolidayResponse
	bodyBytes, err := io.ReadAll(holidays.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	err = json.Unmarshal(bodyBytes, &holidayResponses)
	if err != nil {
		log.Fatalf("Error unmarshalling holidays: %v", err)
	}

	for _, holiday := range holidayResponses {
		name := "Unknown"
		for _, localizedText := range holiday.Name {
			if localizedText.Language == "FR" {
				name = localizedText.Text
				break
			}
		}
		fmt.Printf("Holiday: %s on %s\n", name, holiday.StartDate)
	}

}
```

Result:

```bash
❯ go run main/main.go
Holiday: Nouvel an on 2025-01-01
Holiday: Saint-Berchtold on 2025-01-02
Holiday: Vendredi saint on 2025-04-18
Holiday: Lundi de Pâques on 2025-04-21
Holiday: Ascension on 2025-05-29
Holiday: Lundi de Pentecôte on 2025-06-09
Holiday: Fête nationale suisse on 2025-08-01
Holiday: Jeûne fédéral on 2025-09-21
Holiday: Noël on 2025-12-25
```