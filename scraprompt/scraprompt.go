package scraprompt

import (
	"encoding/json"
	"errors"
	"fmt"
	"scrappy/scrapper"
	"strconv"
	"time"

	"github.com/c-bata/go-prompt"
)

func AddNew(scrapsFile string) {
	newScrapper := scrapper.Scrap{}
	err := errors.New("")
	fmt.Println("Please give a name to your scrapper.")
	newScrapper.Name = prompt.Input("> ", freeText)
	newScrapper.URL = prompt.Input(" URL> ", freeText)
	newScrapper.Attribute = prompt.Input("Attribute to search for> ", freeText)
	trimPrefixChars := prompt.Input("# of prefix to trim (number)> ", freeText)
	newScrapper.TrimPrefixChars, err = strconv.Atoi(trimPrefixChars)
	if err != nil {
		fmt.Println("# of prefix to trim should be a number")
		return
	}
	trimSuffixChars := prompt.Input("# of suffix to trim (number)> ", freeText)
	newScrapper.TrimSuffixChars, err = strconv.Atoi(trimSuffixChars)
	if err != nil {
		fmt.Println("# of suffix to trim should be a number")
		return
	}
	newScrapper.ValueType = scrapper.ValueType(prompt.Input("Type of data (string|integer|float))> ", valueTypeSelector))

	newScrapper.ComparatorType = scrapper.ComparatorType(prompt.Input("Type of comparison)> ", valueTypeComparator))
	if newScrapper.ComparatorType != scrapper.ExistsComparatorType && newScrapper.ComparatorType != scrapper.NotExistsComparatorType {
		newScrapper.CheckValue = prompt.Input("Value to check against ("+string(newScrapper.ComparatorType)+")> ", freeText)
	}
	checkPeriod := prompt.Input("Period to check (12h20m30s)> ", freeText)
	checkPeriodParsed, err := time.ParseDuration(checkPeriod)
	if err != nil {
		fmt.Println("Period to check should be in duration format")
		return
	}
	newScrapper.CheckPeriod = scrapper.Duration{checkPeriodParsed}
	newScrapper.Status = scrapper.ActiveScrapStatus
	err = appendScrapper(scrapsFile, newScrapper)
	if err != nil {
		fmt.Println("Could not save the new scrapper")
		return
	}
	return

}
func appendScrapper(scrapsFile string, newScrapper scrapper.Scrap) error {
	scraps := scrapper.RetrieveScraps(scrapsFile)
	scraps = append(scraps, newScrapper)
	data, _ := json.Marshal(scraps)
	fmt.Println(string(data))
	err := scrapper.WriteScraps(scrapsFile, data)
	return err
}

func freeText(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func valueTypeComparator(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: string(scrapper.LessThanNumberComparatorType), Description: "Less than"},
		{Text: string(scrapper.GreaterThanNumberComparatorType), Description: "Greater than"},
		{Text: string(scrapper.LongerThanStringComparatorType), Description: "Longer than"},
		{Text: string(scrapper.ShorterThanStringComparatorType), Description: "Shorter than"},
		{Text: string(scrapper.ContainsStringComparatorType), Description: "Contains"},
		{Text: string(scrapper.IsSameStringComparatorType), Description: "Is same"},
		{Text: string(scrapper.IsNotSameStringComparatorType), Description: "Is not same"},
		{Text: string(scrapper.ExistsComparatorType), Description: "Exists"},
		{Text: string(scrapper.NotExistsComparatorType), Description: "Not exist"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func valueTypeSelector(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: string(scrapper.StringValueType), Description: "String Type"},
		{Text: string(scrapper.IntegerValueType), Description: "Integer Type"},
		{Text: string(scrapper.FloatValueType), Description: "Float Type"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
