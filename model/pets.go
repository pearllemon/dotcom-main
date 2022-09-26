package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	errNotAMicrochip = fmt.Errorf("Microchip doesn't match the correct format: 9, 10, 12 or 15 character microchip number, with no punctuation or spaces.")
)

func PetByMicrochip(microchip string) (*Pet, error) {
	var pet Pet
	err := PgsqlDB.Preload("Owner").Preload("Registrer").
		Preload("AccountManagers").Limit(1).
		Order("id desc").Where("lower(microchip) = ?", strings.ToLower(microchip)).
		Find(&pet).Error
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

type PetOwnerMetadata struct {
	Email   string
	Phone   string
	Country string
}

func PetInfo(mic string) (contentHTML, contentString string, petOwnerData PetOwnerMetadata, exist bool, err error) {
	pet, err := PetByMicrochip(mic)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			contentHTML = "Microchip not found"
			err = nil
			return
		}
		return
	}

	country := pet.Owner.Country
	if country == "" {
		// default to US
		country = "US"
	}
	petOwnerData = PetOwnerMetadata{
		Email:   pet.Owner.Email,
		Phone:   pet.Owner.Phone,
		Country: country,
	}

	templateImage := `<p><b>Image</b>:<br/> 
	<a href="#" data-toggle="modal" data-target="#modal-image">
	<img src="/image-pets/%s" width="%v" title="Click to enlarge image of %s"/></p>
	</a>
	<div id="modal-image" class="modal fade" role="dialog" style="top:100px">
      <div class="modal-dialog">

        <!-- Modal content-->
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h7 class="modal-title">Image</h7>
          </div>
          <div class="modal-body">
            <div class="col-lg-12">
                <img  class="image-pet-modal" src="/image-pets/%s" width="%v" />
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
          </div>
        </div>

      </div>
    </div> `
	templatePurchasedFrom := `<p><b>Microchip Registration Purchased From</b>: %s </p>`
	// templateWhySTL := `<p><b>Why I Saved This Life?</b> <br/> %v </p>`
	templateContent := `<p><b>Microchip #</b>: %s </p>

<p><b>Pet's Name</b>: %s </p>

<p><b>Species</b>: %s </p>

<p><b>Birthdate</b>: %s </p>

<p><b>Breed</b>: %s </p>

<p><b>Color</b>: %s </p>

<p><b>Gender</b>: %s </p>

<p><b>Microchip Company</b>: Save This Life </p>

%s

%s
<br/>
%v
`

	contentString = fmt.Sprintf(`Microchip #: %s  Pet's Name: %s  Species: %s  Birthdate: %s
Breed: %s  Color:  %s  Gender: %s  Microchip Registration Purchased From: %s`,
		pet.Microchip, pet.Name, strings.Title(pet.Species), pet.Birthday, pet.Breed, pet.Color, pet.Gender,
		pet.PurchasedFrom)

	imageContent := ""
	if pet.ImageName != "" {
		imageContent = fmt.Sprintf(templateImage, pet.ImageName, "50%", pet.Microchip, pet.ImageName, "100%")
	}
	purchasedFromContent := ""
	if pet.PurchasedFrom != "" {
		purchasedFromContent = fmt.Sprintf(templatePurchasedFrom, pet.PurchasedFrom)
	}
	whySTLContent := ""
	// if pet.WhySTL != "" {
	// 	whySTLContent = fmt.Sprintf(templateWhySTL, pet.WhySTL)
	// }
	contentHTML = fmt.Sprintf(templateContent, pet.Microchip, pet.Name, strings.Title(pet.Species), pet.Birthday,
		pet.Breed, pet.Color, pet.Gender, purchasedFromContent, imageContent, whySTLContent)

	exist = true
	return
}

func RecentMicrochips(offset int64) (microchips []string, err error) {
	var pets []Pet
	err = PgsqlDB.Order("updated_at desc").
		Limit(20).Offset(offset).
		Find(&pets).Error
	if err != nil {
		return nil, err
	}

	microchips = make([]string, len(pets))
	for i, pet := range pets {
		microchips[i] = pet.Microchip
	}
	return
}

func CheckValidMicrochip(microchip string) error {
	if len(microchip) != 9 && len(microchip) != 10 && len(microchip) != 12 && len(microchip) != 15 {
		return errNotAMicrochip
	}
	matched, err := regexp.MatchString("^[a-zA-Z0-9]*$", microchip)
	if err != nil {
		return err
	}
	if !matched {
		return errNotAMicrochip
	}
	return nil
}
