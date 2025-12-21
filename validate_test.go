package crud

import (
	"fmt"
	"testing"
	"time"
)

type ValidationTestStruct struct {
	Id             int64
	Flags          int64
	PrimaryEmail   string `crud:"req email"`
	EmailSecondary string `crud:"req email"`
	FirstName      string `crud:"req lenmin:2 lenmax:30"`
	LastName       string `crud:"req lenmin:0 lenmax:255"`
	Age            int    `crud:"valmin:1 valmax:120"`
	Price          int    `crud:"valmin:0 valmax:999"`
	PostCode       string `crud:"req lenmin:6 regexp:^[0-9]{2}\\-[0-9]{3}$"`
	PostCode2      string `crud:"lenmin:6" crud_regexp:"^[0-9]{2}\\-[0-9]{3}$"`
	Password       string
	CreatedBy      int64
	Key            string `crud:"req uniq lenmin:30 lenmax:255"`
}

func validationTestStructWithData() *ValidationTestStruct {
	ts := &ValidationTestStruct{}
	ts.Flags = 4
	ts.PrimaryEmail = "primary@example.com"
	ts.EmailSecondary = "secondary@example.com"
	ts.FirstName = "John"
	ts.LastName = "Smith"
	ts.Age = 37
	ts.Price = 444
	ts.PostCode = "00-000"
	ts.PostCode2 = "11-111"
	ts.Password = "yyy"
	ts.CreatedBy = 4
	ts.Key = fmt.Sprintf("12345679012345678901234567890%d", time.Now().UnixNano())

	return ts
}

// TestValidateWithValidStruct tests if Validate successfully validates object with valid values
func TestValidateWithValidStruct(t *testing.T) {
	ts := validationTestStructWithData()
	ok, violations, err := Validate(ts, nil, "crud")
	if !ok {
		t.Fatal("Validate failed validate valid struct")
	}
	if len(violations) > 0 {
		t.Fatal("Validate return non-empty failed field list when validating a valid struct")
	}
	if err != nil {
		t.Fatalf("Validate failed to validate valid struct: %s", err.Error())
	}
}

// TestValidateWithValidStructAndListOfFields tests if Validate successfully validates object with valid values
func TestValidateWithValidStructAndListOfFields(t *testing.T) {
	ts := validationTestStructWithData()
	ts.Age = 0
	ts.FirstName = "x"
	ts.Key = "tooshort"
	ts.PrimaryEmail = "thisis@valid.email.com"
	ok, violations, err := Validate(ts, map[string]bool{
		"PrimaryEmail": true,
		"Price":        true,
	}, "crud")
	if !ok {
		t.Fatal("Validate failed to validate listed fields")
	}
	if len(violations) > 0 {
		t.Fatal("Validate return non-empty failed field list when validating listed fields")
	}
	if err != nil {
		t.Fatalf("Validate failed to validate listed fields: %s", err.Error())
	}
}

// TestValidateWithInvalidStruct tests if Validate invalidates object with invalid values
func TestValidateWithInvalidStruct(t *testing.T) {
	ts := validationTestStructWithData()
	ts.PrimaryEmail = "invalidemail"
	ts.EmailSecondary = "invalidemail"
	ts.FirstName = "x"
	ts.LastName = "aFbdsZFYxMpUNKCkBrHhhODrMBEHtmRAJjoqSSfUotvsfMXcJGPrCRaDOsyuyrXYfACjsJEMUoxNvTwRMUaWYruOxgzTXJRzobmxaFbdsZFYxMpUNKCkBrHhhODrMBEHtmRAJjoqSSfUotvsfMXcJGPrCRaDOsyuyrXYfACjsJEMUoxNvTwRMUaWYruOxgzTXJRzobmxaFbdsZFYxMpUNKCkBrHhhODrMBEHtmRAJjoqSSfUotvsfMXcJGPrCRaDOsyuyrXYfACjsJEMUoxNvTwRMUaWYruOxgzTXJRzobmxaFbdsZFYxMpUNKCkBrHhhODrMBEHtmRAJjoqSSfUotvsfMXcJGPrCRaDOsyuyrXYfACjsJEMUoxNvTwRMUaWYruOxgzTXJRzobmx"
	ts.Age = 0
	ts.Price = 1000
	ts.PostCode = "inv"
	ts.PostCode2 = "inv"
	ts.Key = "tooshort"
	ok, violations, err := Validate(ts, nil, "crud")
	if err != nil {
		t.Fatal("Validate failed with an err")
	}
	if ok {
		t.Fatal("Validate failed to return false for struct with invalid field values")
	}
	for _, f := range []string{"PrimaryEmail", "EmailSecondary", "FirstName", "LastName", "Age", "Price", "PostCode", "PostCode2", "Key"} {
		if violations[f] == 0 {
			t.Fatalf("Validate failed to return field %s in failed fields", f)
		}
	}
}
