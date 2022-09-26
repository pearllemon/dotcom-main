package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DBInfo  string
	PgsqlDB *gorm.DB

	AppURL = "https://www.savethislife.com"
)

func OpenDB() (*gorm.DB, error) {
	if PgsqlDB != nil {
		return PgsqlDB, nil
	}
	var err error
	PgsqlDB, err = gorm.Open("postgres", DBInfo)
	if err != nil {
		return nil, err
	}
	return PgsqlDB, nil
}

func InitDB() error {
	_, err := OpenDB()
	if err != nil {
		return err
	}
	PgsqlDB.AutoMigrate(&Pet{}, &User{}, &TemplateEmail{}, &PlacePurchaseFrom{},
		&ConfigAutoPopulation{}, &Order{}, &OrderCustomer{}, &OrderPet{},
		&CompaniesKeys{}, &GoneURL{}, &GoogleMeta{}, &OauthToken{}, &OauthClient{}, &VetPartnership{})
	// Enable Logger, show detailed log
	// PgsqlDB.LogMode(true)
	return nil
}

type Pet struct {
	gorm.Model

	Name             string
	Microchip        string `gorm:"type:varchar(100);"`
	Breed            string
	Color            string
	Birthday         string
	Gender           string
	Needs            string
	Species          string
	SecContact       string
	SecPhone         string
	AddPhone1        string
	AddPhone2        string
	AddPhone3        string
	AddPhone4        string
	AddPhone5        string
	AddPhone6        string
	Veterinary       string
	PurchasedFrom    string
	DateVacc         string
	Status           uint
	WhySTL           string
	WhySTLLastUpdate *time.Time
	AcceptGoogle     bool
	ImageName        string
	StripeSubscID    string // stripe subscription ID

	Owner   User `gorm:"ForeignKey:ID;AssociationForeignKey:OwnerID"`
	OwnerID uint

	Registrer   User `gorm:"ForeignKey:ID;AssociationForeignKey:RegistrerID"`
	RegistrerID uint

	AccountManagers []User `gorm:"many2many:pet_users;"` // Many-To-Many relationship
}

type User struct {
	gorm.Model

	Login          string
	HashedPassword string
	Password       string
	DisplayName    string
	Email          string
	Status         uint
	Phone          string
	Address        string
	City           string
	State          string
	Country        string
	ZipCode        string

	ConfigAutoPopulated   ConfigAutoPopulation `gorm:"ForeignKey:ID;AssociationForeignKey:ConfigAutoPopulatedID"`
	ConfigAutoPopulatedID uint

	Role uint `sql:"DEFAULT:4"`
}
type TemplateEmail struct {
	gorm.Model

	Template string
}

type PlacePurchaseFrom struct {
	gorm.Model

	Name string
}

type ConfigAutoPopulation struct {
	gorm.Model

	Species        string
	OwnerName      string
	OwnerEmail     string
	City           string
	State          string
	ZipCode        string
	Country        string
	PlacePurchased string
}

type Order struct {
	gorm.Model

	Status  int
	Total   float64
	VetCode string

	Customer   OrderCustomer `gorm:"ForeignKey:ID;AssociationForeignKey:CustomerID"`
	CustomerID uint

	Pets []OrderPet
}

type OrderCustomer struct {
	gorm.Model

	Name           string
	Phone          string
	AlternatePhone string
	PhoneEmergency string
	Email          string
	Address        string
	Country        string
	City           string
	State          string
	ZipCode        string
	VetName        string
}

type OrderPet struct {
	gorm.Model

	Microchip  string
	Name       string
	Breed      string
	Color      string
	Species    string
	Gender     string
	Birthday   string
	Health     string
	VetContact string

	OrderID uint
}

type CompaniesKeys struct {
	gorm.Model

	CompanyName string
	Key         string
	Status      uint
}

type GoneURL struct {
	gorm.Model

	URL string
}

type GoogleMeta struct {
	gorm.Model

	Content string
	Email   string
}

type OauthToken struct {
	gorm.Model

	ExpiredAt int64
	Code      string `gorm:"type:varchar(512)"`
	Access    string `gorm:"type:varchar(512)"`
	Refresh   string `gorm:"type:varchar(512)"`
	Data      string `gorm:"type:text"`
}

type OauthClient struct {
	gorm.Model

	ClientID string `sql:"unique"`
	Secret   string `sql:"unique"`
	Domain   string

	Owner   User `gorm:"ForeignKey:ID;AssociationForeignKey:OwnerID"`
	OwnerID uint
}

type VetPartnership struct {
	gorm.Model

	Partner string
	Accept  bool

	Vet   User `gorm:"ForeignKey:ID;AssociationForeignKey:VetID"`
	VetID uint

	Timestamp int64
}
