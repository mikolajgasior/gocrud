package test

import (
	"fmt"
	"time"
)

// Test struct for all the tests
type TestStruct struct {
	ID             uint64
	Flags          uint64
	PrimaryEmail   string `crud:"req"`
	EmailSecondary string `crud:"req email"`
	FirstName      string `crud:"req len:2,30"`
	LastName       string `crud:"req len:0,255"`
	Age            int    `crud:"val:1,120"`
	Price          int    `crud:"val:0,999"`
	PostCode       string `crud:"req len:6 regexp:^[0-9]{2}\\-[0-9]{3}$"`
	PostCode2      string `crud:"len:6" 2db_regexp:"^[0-9]{2}\\-[0-9]{3}$"`
	Password       string `json:"password"`
	Key            string `crud:"req uniq len:30,255"`
	CreatedAt      int64
	CreatedBy      uint64
	ModifiedAt     int64
	ModifiedBy     uint64
}

func TestStructWithData() *TestStruct {
	ts := &TestStruct{}
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
	ts.Key = fmt.Sprintf("12345679012345678901234567890%d", time.Now().UnixNano())
	return ts
}

func AreTestStructObjectsSame(ts1 *TestStruct, ts2 *TestStruct) bool {
	if ts1.Flags != ts2.Flags {
		return false
	}
	if ts1.PrimaryEmail != ts2.PrimaryEmail {
		return false
	}
	if ts1.EmailSecondary != ts2.EmailSecondary {
		return false
	}
	if ts1.FirstName != ts2.FirstName {
		return false
	}
	if ts1.LastName != ts2.LastName {
		return false
	}
	if ts1.Age != ts2.Age {
		return false
	}
	if ts1.Price != ts2.Price {
		return false
	}
	if ts1.PostCode != ts2.PostCode {
		return false
	}
	if ts1.PostCode2 != ts2.PostCode2 {
		return false
	}
	if ts1.Password != ts2.Password {
		return false
	}
	if ts1.CreatedBy != ts2.CreatedBy {
		return false
	}
	if ts1.Key != ts2.Key {
		return false
	}
	return true
}
