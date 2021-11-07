package scrapper

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

type ScrapStatus string

const (
	ActiveScrapStatus    ScrapStatus = "active"
	CompletedScrapStatus ScrapStatus = "completed"
	ErrorScrapStatus     ScrapStatus = "error"
)

type ValueType string

func (v ValueType) Stringer() string {
	return string(v)
}

const (
	StringValueType  ValueType = "string"
	IntegerValueType ValueType = "integer"
	FloatValueType   ValueType = "float"
)

type ComparatorType string

func (c ComparatorType) Stringer() string {
	return string(c)
}

const (
	LessThanNumberComparatorType    ComparatorType = "less_than"
	GreaterThanNumberComparatorType ComparatorType = "greater_than"
	LongerThanStringComparatorType  ComparatorType = "longer_than"
	ShorterThanStringComparatorType ComparatorType = "shorter_than"
	ContainsStringComparatorType    ComparatorType = "contains"
	IsSameStringComparatorType      ComparatorType = "is_same"
	IsNotSameStringComparatorType   ComparatorType = "is_not_same"
	ExistsComparatorType            ComparatorType = "exists"
	NotExistsComparatorType         ComparatorType = "not_exists"
)

type Scrap struct {
	Name            string         `json:"name"`
	URL             string         `json:"url"`
	Attribute       string         `json:"attribute"`
	TrimPrefixChars int            `json:"trim_prefix_chars"`
	TrimSuffixChars int            `json:"trim_suffix_chars"`
	ValueType       ValueType      `json:"value_type"`
	CheckValue      string         `json:"check_value"`
	ComparatorType  ComparatorType `json:"comparator_type"`
	CheckPeriod     Duration       `json:"check_period"`
	Status          ScrapStatus    `json:"status"`
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func RetrieveScraps(scrapsFile string) []Scrap {
	scraps := []Scrap{}
	jsonFile, err := os.Open(scrapsFile)
	if nil == err {
		defer jsonFile.Close()
		body, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(body, &scraps)
	}
	return scraps
}
func WriteScraps(scrapsFile string, data []byte) error {
	jsonFile, err := os.OpenFile(scrapsFile, os.O_CREATE|os.O_WRONLY, 644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	_, err = jsonFile.Write(data)
	return err
}
