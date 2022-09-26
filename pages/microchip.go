package pages

import (
	"fmt"
	"net/http"
	"net/mail"
	"stl/model"
	"stl/postquery"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/ttacon/libphonenumber"
)

var (
	templateEmailAlert = `

Please contact: %s | %s ; 

Your pet %s a medical attention ;

Reported Location of your Lost Pet: %s `

	templateSMSAlert = `URGENT PET ALERT FROM SAVETHISLIFE. Please contact:%s|%s;Your pet %s a medical attention; Reported Location of your Lost Pet: %s`
)

func GetMicrochip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	microchip := vars["microchip"]
	// contentHTML, contentString, _, exist, err := model.PetInfo(microchip)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// if contentHTML == "Microchip not found" || !exist {
	// 	http.Error(w, "Not found", http.StatusNotFound)
	// 	return
	// }

	// viewArgs := map[string]interface{}{
	// 	"Title":       fmt.Sprintf("%v | Save This Life Microchip and Pet Recovery System", microchip),
	// 	"Microchip":   microchip,
	// 	"Content":     template.HTML(contentHTML),
	// 	"Description": contentString,
	// 	"SiteKey":     siteKey,
	// }

	// renderHTML(w, "microchipPost.html", viewArgs)

	contentHTML, contentString, _, exist, err := model.PetInfo(microchip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if contentHTML == "Microchip not found" || !exist {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	pet, err := model.PetByMicrochip(microchip)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			viewArgs := map[string]interface{}{
				"StepPetNotFound": true,
			}
			renderHTML(w, "found-pet.html", viewArgs)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewArgs := map[string]interface{}{
		"StepPetForm": true,
		"Pet": map[string]interface{}{
			"ID":        pet.ID,
			"Name":      pet.Name,
			"Microchip": pet.Microchip,
			"Breed":     pet.Breed,
			"Color":     pet.Color,
			"Species":   pet.Species,
			"Gender":    pet.Gender,
			"Image":     pet.ImageName,
		},
		"Description": contentString,
		"Title":       fmt.Sprintf("%v | Save This Life Microchip and Pet Recovery System", pet.Microchip),
	}
	renderHTML(w, "found-pet.html", viewArgs)
	return
}

func PostMicrochip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	microchip := vars["microchip"]
	_, _, petOwnerMetadata, exist, err := model.PetInfo(microchip)
	if err != nil {
		renderJsError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exist {
		renderJsError(w, "microchip not found", http.StatusNotFound)
		return
	}

	url_geoip, err := postquery.GetFirstValueByKey(r, "url_geoip")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}

	info_email, err := postquery.GetFirstValueByKey(r, "info_email")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if info_email != "" {
		_, err = mail.ParseAddress(info_email)
		if err != nil {
			renderJsError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if len(info_email) > 50 {
		renderJsError(w, "The maximum allowed characters is exceeded", http.StatusBadRequest)
		return
	}
	info_phone, err := postquery.GetFirstValueByKey(r, "info_phone")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}
	phone, err := libphonenumber.Parse(info_phone, "US")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !libphonenumber.IsValidNumber(phone) {
		renderJsError(w, "Invalid US phone number", http.StatusBadRequest)
		return
	}
	info_phone = libphonenumber.Format(phone, libphonenumber.NATIONAL)

	info_medical, err := postquery.GetFirstValueByKey(r, "info_medical")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// verify reCAPTCHA
	// isValid := re.Verify(*r)
	// if !isValid {
	// 	renderJsError(w, fmt.Sprintf("reCAPTCHA error: %v", re.LastError()), http.StatusBadRequest)
	// 	return
	// }

	emailBody := fmt.Sprintf(templateEmailAlert, info_email, info_phone, info_medical, url_geoip)
	smsBody := fmt.Sprintf(templateSMSAlert, info_email, info_phone, info_medical, url_geoip)

	owner_email := "alerts@savethislife.com"
	if petOwnerMetadata.Email != "" {
		owner_email = petOwnerMetadata.Email
	}

	err = sendMail(owner_email, "FOUND PET ALERT - "+microchip,
		[]string{},
		// []string{"gmailcomsavethislifealerts@microchip.freshdesk.com", "stl.said@gmail.com", "savethislifealerts@gmail.com", "alerts@savethislife.com"},
		[]string{"alerts@savethislife.com", "support@savethislife.com"},
		emailBody, "SaveThisLife<admin@savethislife.com>")
	if err != nil {
		renderJsError(w, err.Error(), http.StatusBadRequest)
		return
	}

	prefixMsg := "Email successfully sent but error when sending SMS: "

	if petOwnerMetadata.Phone != "" && petOwnerMetadata.Country != "" {
		phoneNumber, err := libphonenumber.Parse(petOwnerMetadata.Phone, petOwnerMetadata.Country)
		if err != nil {
			renderJsError(w, prefixMsg+err.Error(), http.StatusBadRequest)
			return
		}
		formattedPhoneNumber := fmt.Sprintf("+%d%v ", phoneNumber.GetCountryCode(), phoneNumber.GetNationalNumber())
		err = sendTwilioSMS(formattedPhoneNumber, smsBody)
		if err != nil {
			renderJsError(w, prefixMsg+err.Error(), http.StatusBadRequest)
			return
		}
	}

	renderText(w, "Email and text message are successfully sent")
	return
}
