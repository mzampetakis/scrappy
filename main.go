package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"scrappy/informer"
	"scrappy/scrapper"
	"scrappy/scraprompt"
	"strconv"
	"strings"

	"github.com/robfig/cron/v3"

	"github.com/gocolly/colly"
)

const (
	ScrapsFile     = "scraps.json"
	MailConfigFile = "mail.conf"
)

func main() {
	mode := flag.String("mode", "run", "Mode to start {run|add|remove};.")
	flag.Parse()
	fmt.Println("Staring at - " + *mode + " - mode")
	switch *mode {
	case "run":
		runScrappy()
	case "add":
		scraprompt.AddNew(ScrapsFile)
	case "remove":
		fmt.Println("remove")
	default:
		fmt.Println("Mode " + *mode + " is not supported")
	}
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, os.Interrupt)
	defer func() {
		signal.Stop(osSig)
	}()
	go func() {
		<-osSig
		os.Exit(0)
	}()
	fmt.Println("Exiting...")
	os.Exit(0)
}

func runScrappy() {
	scraps := scrapper.RetrieveScraps(ScrapsFile)
	mailConfig, err := retrieveMailCnf(MailConfigFile)
	if err != nil {
		panic(err)
		return
	}
	c := cron.New()
	fmt.Println("Srappy is starting...")
	for _, scrap := range scraps {
		scrapToProcess := scrap
		fmt.Println("Registering: " + scrap.Name + ", every " + "@every " + scrap.CheckPeriod.Duration.String())
		_, err := c.AddFunc(("@every " + scrap.CheckPeriod.Duration.String()), func() {
			handleScrap(&scrapToProcess, mailConfig)
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	c.Run()
}

func handleScrap(scrap *scrapper.Scrap, userInformer informer.InformUser) {
	if scrap.Status == scrapper.CompletedScrapStatus {
		return
	}
	c := colly.NewCollector()
	fmt.Print(scrap.Name + ": ")
	c.OnHTML("body", func(e *colly.HTMLElement) {
		searchResult, err := e.DOM.Children().Find(scrap.Attribute).First().Html()
		if err != nil {
			scrap.Status = scrapper.ErrorScrapStatus
			fmt.Println("Failed to retrieve content " + err.Error())
			return
		}
		if len(searchResult) <= scrap.TrimPrefixChars+scrap.TrimSuffixChars {
			fmt.Println("Failed to trim (" + fmt.Sprint(len(searchResult)) + ")")
			return
		}
		searchResult = searchResult[scrap.TrimPrefixChars : len(searchResult)-scrap.TrimSuffixChars]

		fmt.Print("(" + searchResult + ") ")
		if len(searchResult) > 0 && scrap.ComparatorType == scrapper.ExistsComparatorType {
			message := scrap.Name + ": value found."
			fmt.Println(message)
			err = userInformer.Inform("Scrappy: "+scrap.Name, message)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		if len(searchResult) == 0 && scrap.ComparatorType == scrapper.NotExistsComparatorType {
			message := scrap.Name + ": value not found."
			fmt.Println(message)
			err = userInformer.Inform("Scrappy: "+scrap.Name, message)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		switch scrap.ValueType {
		case scrapper.IntegerValueType:
			checkValue, err := strconv.ParseInt(scrap.CheckValue, 10, 64)
			if err != nil {
				scrap.Status = scrapper.ErrorScrapStatus
				fmt.Println("Failed to convert to integer " + err.Error())
				return
			}
			intResult, err := strconv.ParseInt(searchResult, 10, 64)
			if err != nil {
				scrap.Status = scrapper.ErrorScrapStatus
				fmt.Println("Failed to convert to integer " + err.Error())
				return
			}
			if scrap.ComparatorType == scrapper.LessThanNumberComparatorType && intResult < checkValue {
				message := scrap.Name + ": value " + searchResult + " became less than " + scrap.CheckValue + "."
				fmt.Println(message)
				err = userInformer.Inform("Scrappy: "+scrap.Name, message)
				if err != nil {
					fmt.Println(err.Error())
				}
				return
			}
			if scrap.ComparatorType == scrapper.GreaterThanNumberComparatorType && intResult > checkValue {
				message := scrap.Name + ": value " + searchResult + " became greater than " + scrap.CheckValue + "."
				fmt.Println(message)
				err = userInformer.Inform("Scrappy: "+scrap.Name, message)
				if err != nil {
					fmt.Println(err.Error())
				}
				return
			}
		case scrapper.FloatValueType:
			floatValue := strings.Replace(scrap.CheckValue, ",", ".", 1)
			checkValue, err := strconv.ParseFloat(floatValue, 64)
			if err != nil {
				scrap.Status = scrapper.ErrorScrapStatus
				fmt.Println("Failed to convert to float " + err.Error())
				return
			}
			floatValue = strings.Replace(searchResult, ",", ".", 1)
			floatResult, err := strconv.ParseFloat(floatValue, 64)
			if err != nil {
				scrap.Status = scrapper.ErrorScrapStatus
				fmt.Println("Failed to convert to float " + err.Error())
				return
			}
			if scrap.ComparatorType == scrapper.LessThanNumberComparatorType && floatResult < checkValue {
				message := scrap.Name + ": value " + searchResult + " became less than " + scrap.CheckValue + "."
				fmt.Println(message)
				err = userInformer.Inform("Scrappy: "+scrap.Name, message)
				if err != nil {
					fmt.Println(err.Error())
				}
				return
			}
			if scrap.ComparatorType == scrapper.GreaterThanNumberComparatorType && floatResult > checkValue {
				message := scrap.Name + ": value " + searchResult + " became greater than " + scrap.CheckValue + "."
				fmt.Println(message)
				err = userInformer.Inform("Scrappy: "+scrap.Name, message)
				if err != nil {
					fmt.Println(err.Error())
				}
				return
			}
		case scrapper.StringValueType:
			if scrap.ComparatorType == scrapper.IsSameStringComparatorType &&
				searchResult == scrap.CheckValue {
				message := scrap.Name + ": value " + searchResult + " is the same with " + scrap.CheckValue + "."
				fmt.Println(message)
				err = userInformer.Inform("Scrappy: "+scrap.Name, message)
				if err != nil {
					fmt.Println(err.Error())
				}
				return
			}
			if scrap.ComparatorType == scrapper.IsNotSameStringComparatorType &&
				searchResult != scrap.CheckValue {
				message := scrap.Name + ": value " + searchResult + " is not the same with " + scrap.CheckValue + "."
				fmt.Println(message)
				err = userInformer.Inform("Scrappy: "+scrap.Name, message)
				if err != nil {
					fmt.Println(err.Error())
				}
				return
			}
			if scrap.ComparatorType == scrapper.ContainsStringComparatorType &&
				strings.Contains(searchResult, scrap.CheckValue) {
				message := scrap.Name + ": value " + searchResult + " is contained in " + scrap.CheckValue + "."
				fmt.Println(message)
				err = userInformer.Inform("Scrappy: "+scrap.Name, message)
				if err != nil {
					fmt.Println(err.Error())
				}
				return
			}
			if scrap.ComparatorType == scrapper.LongerThanStringComparatorType {
				checkValue, err := strconv.ParseInt(scrap.CheckValue, 10, 64)
				if err != nil {
					fmt.Println("Failed to convert to integer " + err.Error())
					return
				}
				if len(searchResult) > int(checkValue) {
					message := scrap.Name + ": value " + searchResult + " larger then " + scrap.CheckValue + "."
					fmt.Println(message)
					err = userInformer.Inform("Scrappy: "+scrap.Name, message)
					if err != nil {
						fmt.Println(err.Error())
					}
					return
				}
			}
			if scrap.ComparatorType == scrapper.ShorterThanStringComparatorType {
				checkValue, err := strconv.ParseInt(scrap.CheckValue, 10, 64)
				if err != nil {
					scrap.Status = scrapper.ErrorScrapStatus
					fmt.Println("Failed to convert to integer " + err.Error())
					return
				}
				if len(searchResult) < int(checkValue) {
					message := scrap.Name + ": value " + searchResult + " is shorter than " + scrap.CheckValue + "."
					fmt.Println(message)
					err = userInformer.Inform("Scrappy: "+scrap.Name, message)
					if err != nil {
						fmt.Println(err.Error())
					}
					return
				}
			}
		default:
			fmt.Println("Not Supported Type")
			scrap.Status = scrapper.ErrorScrapStatus
		}
	})
	c.Visit(scrap.URL)
	fmt.Println("")
}

func retrieveMailCnf(mailConfigFile string) (*informer.MailConfig, error) {
	jsonFile, err := os.Open(mailConfigFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	body, _ := ioutil.ReadAll(jsonFile)
	mailCfg := informer.MailConfig{}
	err = json.Unmarshal(body, &mailCfg)
	if err != nil {
		return nil, err
	}
	return &mailCfg, nil
}
