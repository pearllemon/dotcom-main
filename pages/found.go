package pages

import (
	"fmt"
	"net/http"
	"stl/model"
	"strings"

	"github.com/jinzhu/gorm"
)

func HandleFoundPetGET(w http.ResponseWriter, r *http.Request) {
	qvalues := r.URL.Query()
	mic := qvalues.Get("microchip")
	mic = strings.TrimSpace(mic)
	if mic == "" {
		viewArgs := map[string]interface{}{
			"StepSearch": true,
		}
		renderHTML(w, "found-pet.html", viewArgs)
		return
	}

	pet, err := model.PetByMicrochip(mic)
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
		"Description": pet.Microchip,
		"Title":       fmt.Sprintf("%v | Save This Life Microchip and Pet Recovery System", pet.Microchip),
	}
	renderHTML(w, "found-pet.html", viewArgs)
	return
}
