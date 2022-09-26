package pages

import (
	"fmt"

	"github.com/ttacon/libphonenumber"
)

func ExamplePhoneNumber() {
	phone, err := libphonenumber.Parse("eds23355jhgvj", "US")
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	format := libphonenumber.Format(phone, libphonenumber.NATIONAL)

	info_phone := fmt.Sprintf("%v", phone)
	fmt.Println(info_phone)
	info_format := fmt.Sprintf("%v", format)
	fmt.Println(info_format)
	// Output:
	// test
}
