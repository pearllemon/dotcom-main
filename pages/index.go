package pages

import (
	"net/http"
	"stl/model"
	"strconv"
)

var (
	PublicFolder string
)

func Index(w http.ResponseWriter, r *http.Request) {
	metas, err := model.AllGoogleMetas()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	viewArgs := map[string]interface{}{
		"GoogleMeta": metas,
	}
	renderHTML(w, "index.html", viewArgs)
	return
}

func HttpGone(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGone)
	return
}

func Maintenance(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "maintenance.html", nil)
	return
}

func AboutUs(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "about.html", nil)
	return
}

func RobotsText(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.Dir(PublicFolder))
	fileServer.ServeHTTP(w, r)
	return
}

func PackageInfo(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "packageInfo.html", nil)
	return
}

func Microchip(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "microchip.html", nil)
	return
}

func PetHealthInsurance(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "petHealthInsurance.html", nil)
	return
}

func Faq(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://desk.zoho.com/portal/savethislifeinc/home", http.StatusMovedPermanently)
	// renderHTML(w, "faq.html", nil)
	return
}

func OrderReceived(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "orderReceived.html", nil)
	return
}

func RegistrationTypes(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "registration-types.html", nil)
	return
}

func Geoip(w http.ResponseWriter, r *http.Request) {
	qvalues := r.URL.Query()
	lat, long := qvalues.Get("lat"), qvalues.Get("long")
	viewArgs := map[string]interface{}{
		"NoCoordinates": lat == "" || long == "",
	}
	renderHTML(w, "geoip.html", viewArgs)
	return
}

func Recents(w http.ResponseWriter, r *http.Request) {
	qvalues := r.URL.Query()
	var offset int64

	if qvalues.Get("offset") != "" {
		var err error
		offset, err = strconv.ParseInt(qvalues.Get("offset"), 0, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	microchips, err := model.RecentMicrochips(offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newer := func() int64 {
		if offset < 20 {
			return -1
		} else {
			return offset - 20
		}
	}

	viewArgs := map[string]interface{}{
		"Microchips": microchips,
		"Older":      offset + 20,
		"Newer":      newer(),
		"AppURL":     model.AppURL,
	}

	renderHTML(w, "recents.html", viewArgs)
	return
}

func LostFoundPet(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "lostFoundPet.html", nil)
	return
}

func PrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "privacyPolicy.html", nil)
	return
}

func WhyYouSaveThisLife(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "whystl.html", nil)
	return
}

func Help(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://desk.zoho.com/portal/savethislifeinc/en/kb", http.StatusTemporaryRedirect)
	// renderHTML(w, "help.html", nil)
	return
}
