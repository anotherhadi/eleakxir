package parquet

import (
	"slices"
)

func getSuggestion(col string) string {
	col = formatColumnName(col)
	knownNames := []string{
		"date",
		"phone",
		"username",
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
	if slices.Contains(knownNames, col) {
		return col
	}
	if col == "user" {
		return "username"
	}
	if col == "login" {
		return "username"
	}
	if col == "sex" {
		return "gender"
	}
	if col == "ip_address" {
		return "ip"
	}
	if col == "password_hashed" {
		return "password_hash"
	}
	if col == "firstname" {
		return "first_name"
	}
	if col == "lastname" {
		return "last_name"
	}
	if col == "fullname" {
		return "full_name"
	}
	if col == "mail" {
		return "email"
	}
	if col == "zip" || col == "postalcode" || col == "zipcode" || col == "postal" || col == "zip_code" {
		return "postal_code"
	}
	if col == "street_address" {
		return "address"
	}
	if col == "hash" || col == "hashed_password" || col == "hash_password" {
		return "password_hash"
	}
	if col == "birthdate" || col == "dob" || col == "date_of_birth" {
		return "birth_date"
	}

	return ""
}

// HINTS:
// date: _date
// url: _url, link
// address: _address
//
