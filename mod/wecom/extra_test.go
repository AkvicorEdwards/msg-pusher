package wecom

import "fmt"

func ExampleExtraModel_String() {
	extra := ExtraModel{
		Type:    "textcard",
		ToUser:  "Akvicor",
		ToParty: "",
		ToTag:   "",
		Url:     "url",
		Btn:     "BUTTON",
	}
	fmt.Println(extra.String())

	// Output:
	// {"type":"textcard","touser":"Akvicor","toparty":"","totag":"","url":"url","btn":"BUTTON"}
}
