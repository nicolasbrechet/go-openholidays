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
