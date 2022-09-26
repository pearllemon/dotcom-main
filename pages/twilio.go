package pages

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-acme/lego/v3/platform/config/env"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	// Set initial variables FOR Twilio
	accountSid    = "AC0898f0a47a642a4928b89f98c93df3da"
	urlStr        = "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	twiolioNumber = "+18305496062"
)

func authToken() string {
	return env.GetOrDefaultString("TWILIO_AUTH_TOKEN", "")
}

func sendTwilioSMS(phoneNumber, bodymsg string) (err error) {
	// Build out the data for our message
	v := url.Values{}
	v.Set("To", phoneNumber)
	v.Set("From", twiolioNumber)
	v.Set("Body", bodymsg)
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, err := http.NewRequest("POST", urlStr, &rb)
	if err != nil {
		return err
	}
	req.SetBasicAuth(accountSid, authToken())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("SMS cannot be sent: error=%s", resp.Status)
	}
	return nil
}

func MakeTestCall() error {
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username:   accountSid,
		Password:   authToken(),
		AccountSid: accountSid,
	})

	params := &openapi.CreateCallParams{}
	params.SetTo("+16692816656")
	params.SetFrom("+18557772447")
	// params.SetUrl("http://demo.twilio.com/docs/voice.xml")
	params.SetTwiml(`<Response>
    <Say>Hi, Let me call the owner!</Say>
	<Dial action="https://stl-8745.twil.io/async-callback">669-333-1241</Dial>
</Response>`)

	resp, err := client.ApiV2010.CreateCall(params)
	if err != nil {
		fmt.Println(err.Error())
		err = nil
	} else {
		fmt.Println("Call Status: " + *resp.Status)
		fmt.Println("Call Sid: " + *resp.Sid)
		fmt.Println("Call Direction: " + *resp.Direction)
	}
	return nil
}
