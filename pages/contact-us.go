package pages

import (
	"fmt"
	"net/http"
	"stl/postquery"
)

func GetContactUs(w http.ResponseWriter, r *http.Request) {
	viewArgs := map[string]interface{}{
		"PostURL": "/contact-us",
		"SiteKey": siteKey,
	}
	renderHTML(w, "contactUs.html", viewArgs)
	return
}

func PostContactUs(w http.ResponseWriter, r *http.Request) {

	// test here
	// MakeTestCall()
	// remove until here
	values := make(map[string]string)
	for _, key := range []string{"name", "phone", "email", "subject",
		"microchip", "message"} {
		value, err := postquery.GetFirstValueByKey(r, key)
		if err != nil {
			renderJsError(w, err.Error(), http.StatusBadRequest)
			return
		}
		values[key] = value
	}

	isValid := re.Verify(*r)
	if !isValid {
		renderJsError(w, fmt.Sprintf("reCAPTCHA error: %v", re.LastError()), http.StatusBadRequest)
		return
	}

	// emailReceiver := "gmailcomsavethislifealerts@microchip.freshdesk.com"
	emailReceiver := "support@savethislifeinc.zohodesk.com"
	subject := fmt.Sprintf("CONTACT US FORM - %v", values["subject"])

	body := fmt.Sprintf(`#original_sender {%s}
Name: %v

Phone: %v

Email : %v

Microchip(s) # : %v

Message : 

%v`, values["email"], values["name"], values["phone"], values["email"], values["microchip"], values["message"])

	err := sendMail(emailReceiver, subject,
		// []string{"savethislifeinfo@gmail.com", "savethislifealerts@gmail.com"},
		[]string{},
		[]string{}, body,
		"SaveThisLife<admin@savethislife.com>")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderText(w, "Your message is successfully sent. Thank you for contacting us!")
	return
}
