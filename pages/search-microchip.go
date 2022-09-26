package pages

import (
	"net/http"
	"stl/model"
	"strings"
)

func GetSearchMicrochip(w http.ResponseWriter, r *http.Request) {
	qvalues := r.URL.Query()
	mic := qvalues.Get("search")
	pet, err := model.PetByMicrochip(strings.TrimSpace(mic))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pets := func() []interface{} {
		return []interface{}{
			map[string]string{
				"Microchip":     pet.Microchip,
				"PurchasedFrom": pet.PurchasedFrom,
			},
		}
	}

	viewArgs := map[string]interface{}{
		"Pets": pets(),
	}
	renderHTML(w, "search-results.html", viewArgs)
	return
}
