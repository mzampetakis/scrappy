# Scrappy - A multi-type web scrapper with alerting

Scrappy is a cli tool that allows multiple web scrappers to monitor periodically 
for a basic ruleset coverage and inform users when the criteria have been met.

## About
Scrappy can accept and manage multiple rule sets to scrap each one on its own 
period. Each scrap rule is consisted of a user-friendly name and all the 
required data used for scrapping and alerting.
These fields are:
```
name - a friendly name of the scrapper (string)
url - a URL to scrap from (string)
attribute - the attribute used to scrap from (string)
trim_prefix_chars - prefix characters to trim (integer)
trim_suffix_chars - suffix characters to trim (integer)
value_type - value's type (string|integer|float)
check_value - value to check against (string parsed to value_type)
comparator_type - comparison type between value_type and check_value
check_period - Period to check for change (duration of type 1h30m10s)
status current status of the scrapped (active|error|complete)
```

The `comparator_type` field can have one of the following values depending on the
`value_type`.
For `string` value_type the applicable `comparator_type` are:
```
longer_than"
shorter_than"
contains"
is_same"
is_not_same"
exists"
not_exists"
```
For `integer` and `float` value_type the applicable `comparator_type` are:
```
less_than
greater_than
exists
not_exists
```

When a scrap rule is met an email is sent to the set-up email account.

## Installation

Clone the repo with
```
$ go get -u github.com/mzampetakis/scrappy
```
> Go modules are required.

In order to build the project use:
```
go build -o scrappy main.go
```

## Usage

### Configuration

A valid configuration is required for accessing an email account.
The configuration for email account used to send the email can be placed at the
`email.conf` with the following format:
```
{
    "email": "some@email.com",
    "password": "mails_password"
}
```

> The email is used through SMTP protocol.

> Don't use your personal email password. Issue a third party account access.

The `scraps.json` files contains all the rules for the available scraps. 
It's a json formatted text file. It is highly recommended using the 
available CLI tools to manage this file.

### Adding a scrapper

Adding a new scrap rule-set can be done through the available CLI command.
The command to add a new scrap is
```
./scrappy --mode add
```
The CLI will prompt for the required fields and validate the given values.

### Starting the scrap
In order to start the scrap to run use
```
./scrappy
```
or in order to start it and let it on the background to run
```
./scrappy &
```

## Contribute
You can contribute to this project by just opening a PR or open first an issue. 
Please describe thoroughly what are your PR solves or adds.

Some ideas for contribution:

* Add other types of informers
* Use separate email per scrapper
* Improve CLIs in terms of suggestions
* Your idea here...