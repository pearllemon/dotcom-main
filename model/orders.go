package model

import "github.com/jinzhu/gorm"

type PetPostRequest struct {
	Mic     string `json:"mic"`
	Name    string `json:"name"`
	Breed   string `json:"breed"`
	Color   string `json:"color"`
	Gender  string `json:"gender"`
	Health  string `json:"health"`
	Vet     string `json:"vet"`
	Species string `json:"species"`
}

type OwnerPostRequest struct {
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Carrier        string `json:"carrier"`
	AlternatePhone string `json:"alternatephone"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	Country        string `json:"country"`
	City           string `json:"city"`
	State          string `json:"state"`
	Zipcode        string `json:"zipcode"`
	Contacts       string `json:"contacts"`
	IsVet          bool   `json:"isvet"`
	NameVet        string `json:"namevet"`
}

type OrderPostRequest struct {
	DateAdded string           `json:"dateadded"`
	Total     float64          `json:"total"`
	Status    bool             `json:"status"`
	VetCode   string           `json:"vetcode"`
	OwnerId   OwnerPostRequest `json:"owner"`
	PetInfo   []PetPostRequest `json:"pet"`
}

func InsertOrder(order OrderPostRequest) (idorder uint, err error) {
	tx := PgsqlDB.Begin()

	lastCustomer := OrderCustomer{}
	err = tx.Unscoped().Model(&lastCustomer).Order("id desc").First(&lastCustomer).Error
	if err != nil {
		tx.Rollback()
		return
	}

	lastOrder := Order{}
	err = tx.Unscoped().Model(&lastOrder).Order("id desc").First(&lastOrder).Error
	if err != nil {
		tx.Rollback()
		return
	}

	c := &OrderCustomer{
		Address:        order.OwnerId.Address,
		AlternatePhone: order.OwnerId.AlternatePhone,
		City:           order.OwnerId.City,
		Country:        order.OwnerId.Country,
		Email:          order.OwnerId.Email,
		Name:           order.OwnerId.Name,
		Phone:          order.OwnerId.Phone,
		PhoneEmergency: order.OwnerId.Contacts,
		State:          order.OwnerId.State,
		VetName:        order.OwnerId.NameVet,
		ZipCode:        order.OwnerId.Zipcode,

		Model: gorm.Model{ID: lastCustomer.ID + 1},
	}

	err = tx.Create(c).Error
	if err != nil {
		tx.Rollback()
		return
	}

	o := &Order{
		Status:     0,
		Total:      order.Total,
		VetCode:    order.VetCode,
		CustomerID: c.ID,

		Model: gorm.Model{ID: lastOrder.ID + 1},
	}

	err = tx.Create(o).Error
	if err != nil {
		tx.Rollback()
		return
	}

	for _, petInfo := range order.PetInfo {
		opet := &OrderPet{
			Breed:      petInfo.Breed,
			Color:      petInfo.Color,
			Gender:     petInfo.Gender,
			Health:     petInfo.Health,
			Microchip:  petInfo.Mic,
			Name:       petInfo.Name,
			Species:    petInfo.Species,
			VetContact: petInfo.Vet,
			OrderID:    o.ID,
		}

		lastPet := OrderPet{}
		err = tx.Unscoped().Model(&lastPet).Order("id desc").First(&lastPet).Error
		if err != nil {
			tx.Rollback()
			return
		}

		opet.Model = gorm.Model{ID: lastPet.ID + 1}

		err = tx.Create(&opet).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return
	}

	return o.ID, nil
}
