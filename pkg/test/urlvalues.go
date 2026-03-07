package test

import (
	"fmt"
	"net/url"
	"time"
)

func TestStructURLValues() url.Values {
	values := url.Values{}
	values.Set("Flags", "4")
	values.Set("PrimaryEmail", "primary@example.com")
	values.Set("EmailSecondary", "secondary@example.com")
	values.Set("FirstName", "John")
	values.Set("LastName", "Smith")
	values.Set("Age", "37")
	values.Set("Price", "444")
	values.Set("PostCode", "00-000")
	values.Set("PostCode2", "11-111")
	values.Set("Password", "yyy")
	values.Set("Key", fmt.Sprintf("12345679012345678901234567890%d", time.Now().UnixNano()))
	return values
}
