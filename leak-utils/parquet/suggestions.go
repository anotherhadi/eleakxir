package parquet

import (
	"regexp"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	knownColumnNames = []string{
		"date",
		"creation_date",
		"bio",
		"phone",
		"username",
		"iban",
		"address",
		"email",
		"postal_code",
		"city",
		"country",
		"state",
		"age",
		"gender",
		"password",
		"password_hash",
		"full_name",
		"last_name",
		"name", // Will be renamed to full_name later
		"first_name",
		"birth_date",
		"url",
		"ip",
	}

	suggestions = map[string]string{
		"user":            "username",
		"pseudo":          "username",
		"login":           "username",
		"description":     "bio",
		"sex":             "gender",
		"civilite":        "gender",
		"genre":           "gender",
		"ipaddress":       "ip",
		"firstname":       "first_name",
		"prenom":          "first_name",
		"lastname":        "last_name",
		"nom":             "last_name",
		"fullname":        "full_name",
		"nomcomplet":      "full_name",
		"adresse":         "address",
		"streetaddress":   "address",
		"ville":           "city",
		"pays":            "country",
		"mail":            "email",
		"zip":             "postal_code",
		"postalcode":      "postal_code",
		"zipcode":         "postal_code",
		"postal":          "postal_code",
		"codepostal":      "postal_code",
		"hash":            "password_hash",
		"hashedpassword":  "password_hash",
		"hashpassword":    "password_hash",
		"passwordhashed":  "password_hash",
		"passwd":          "password",
		"birthdate":       "birth_date",
		"dob":             "birth_date",
		"dateofbirth":     "birth_date",
		"datenaissance":   "birth_date",
		"birthday":        "birth_date",
		"datecreation":    "creation_date",
		"datedecreation":  "creation_date",
		"createdat":       "creation_date",
		"phonenumber":     "phone",
		"numero":          "phone",
		"numerotelephone": "phone",
		"numeromobile":    "phone",
		"mobilephone":     "phone",
		"mobile":          "phone",
	}
)

func getSuggestion(col string) string {
	colFormated := formatColumnName(col)
	if slices.Contains(knownColumnNames, colFormated) {
		return colFormated
	}

	col = cleanString(col)

	if val, ok := suggestions[col]; ok {
		return val
	}

	return ""
}

// HINTS:
// date: _date
// url: _url, link
// address: _address
//

func cleanString(input string) string {
	t := norm.NFD.String(input)
	var sb strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		sb.WriteRune(r)
	}
	s := strings.ToLower(sb.String())
	reg, _ := regexp.Compile("[^a-z]+")
	s = reg.ReplaceAllString(s, "")
	return s
}
