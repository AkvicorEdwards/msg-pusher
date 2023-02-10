package mod

import "fmt"

func ExampleMessageModel_String() {
	msg := MessageModel{
		Title:      "t_title",
		Content:    "t_content",
		Sender:     "t_sender",
		Urgency:    1,
		TimeCreate: 2,
		TimeSend:   3,
		Extra:      make(map[string]map[string]string),
	}
	msg.Extra["t_mod"] = make(map[string]string)
	msg.Extra["t_mod"]["t_key"] = "t_value"
	fmt.Println(msg.String())

	// Output:
	// {"title":"t_title","content":"t_content","sender":"t_sender","urgency":1,"time_create":2,"time_send":3,"extra":{"t_mod":{"t_key":"t_value"}}}
}
