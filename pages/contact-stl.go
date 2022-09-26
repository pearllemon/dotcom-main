package pages

import (
	"fmt"
	"net/http"
	"stl/postquery"
)

func GetContactStl(w http.ResponseWriter, r *http.Request) {
	qvalues := r.URL.Query()
	id := qvalues.Get("id")
	viewArgs := map[string]interface{}{
		"IsFields1": id == "1" || id == "3",
		"IsFields2": id == "2",
		"IsFields3": id == "4",
	}
	renderHTML(w, "contactSaveThisLife.html", viewArgs)
	return
}

func PostContactStl(w http.ResponseWriter, r *http.Request) {
	qvalues := r.URL.Query()
	id := qvalues.Get("id")

	firstName, err := postquery.GetFirstValueByKey(r, "name1")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}
	lastName, err := postquery.GetFirstValueByKey(r, "name2")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}
	mic, err := postquery.GetFirstValueByKey(r, "mic")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}
	phone, err := postquery.GetFirstValueByKey(r, "phone")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}
	email, err := postquery.GetFirstValueByKey(r, "email")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}
	desc, err := postquery.GetFirstValueByKey(r, "desc")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// emailReceiver := "savethislifeinfo@gmail.com"
	emailReceiver := "support@savethislifeinc.zohodesk.com"
	msgSuccess := "Thank you for contacting us! A Save This Life specialist will respond promptly."

	switch id {
	case "1":
		if firstName == "" || lastName == "" || mic == "" || phone == "" || email == "" {
			renderJsError(w, "A required field is empty!", http.StatusBadRequest)
			return
		}

		body := fmt.Sprintf(`#original_sender {%s}
First name: %v
Last name : %v
Microchip : %v
Phone : %v
Email : %v`, email, firstName, lastName, mic, phone, email)
		err = sendMail(emailReceiver, "CHECK IF REGISTERED - "+mic,
			// []string{"stl.said@gmail.com", "savethislifealerts@gmail.com"},
			[]string{},
			// []string{"gmailcomsavethislifealerts@microchip.freshdesk.com"},
			[]string{},
			body, "SaveThisLife<admin@savethislife.com>")
		if err != nil {
			renderJsError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderText(w, msgSuccess)
		return

	case "2":
		if firstName == "" || lastName == "" || mic == "" || phone == "" || email == "" || desc == "" {
			renderJsError(w, "A required field is empty!", http.StatusBadRequest)
			return
		}
		body := fmt.Sprintf(`#original_sender {%s}
First name: %v
Last name : %v
Description: %v
Microchip : %v
Phone : %v
Email : %v`, email, firstName, lastName, desc, mic, phone, email)
		err = sendMail(emailReceiver, "REPORT LOST PET - "+mic,
			// []string{"stl.said@gmail.com", "savethislifealerts@gmail.com"},
			[]string{},
			// []string{"gmailcomsavethislifealerts@microchip.freshdesk.com"},
			[]string{},
			body, "SaveThisLife<admin@savethislife.com>")
		if err != nil {
			renderJsError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderText(w, msgSuccess)
		return

	case "3":
		if firstName == "" || lastName == "" || mic == "" || phone == "" || email == "" {
			renderJsError(w, "A required field is empty!", http.StatusBadRequest)
			return
		}
		body := fmt.Sprintf(`#original_sender {%s}
First name: %v
Last name : %v
Microchip : %v
Phone : %v
Email : %v`, email, firstName, lastName, mic, phone, email)
		err = sendMail(emailReceiver, "ISSUES LOGGING IN - "+mic,
			// []string{"stl.said@gmail.com", "savethislifealerts@gmail.com"},
			[]string{},
			// []string{"gmailcomsavethislifealerts@microchip.freshdesk.com"},
			[]string{},
			body, "SaveThisLife<admin@savethislife.com>")
		if err != nil {
			renderJsError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderText(w, msgSuccess)
		return

	case "4":
		if firstName == "" || lastName == "" || mic == "" || phone == "" || email == "" {
			renderJsError(w, "A required field is empty!", http.StatusBadRequest)
			return
		}

		body := fmt.Sprintf(`#original_sender {%s}
First name: %v
Last name : %v
Microchips : %v
Phone : %v
Email : %v`, email, firstName, lastName, mic, phone, email)
		err = sendMail(emailReceiver, "COMBINE MY ACCOUNTS - "+email,
			// []string{"stl.said@gmail.com", "savethislifealerts@gmail.com"},
			[]string{},
			// []string{"gmailcomsavethislifealerts@microchip.freshdesk.com"},
			[]string{},
			body, "SaveThisLife<admin@savethislife.com>")
		if err != nil {
			renderJsError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderText(w, msgSuccess)
		return

	default:
		renderJsError(w, fmt.Sprintf("Unknow form id '%v'", id), http.StatusInternalServerError)
		return
	}
}
