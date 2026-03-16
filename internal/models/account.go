package models

import "time"

type Account struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BrojRacuna        string    `gorm:"uniqueIndex;size:18;not null" json:"broj_racuna"`
	ClientID          *uint     `json:"client_id"`         // nullable for bank accounts
	FirmaID           *uint     `json:"firma_id"`          // nullable for personal accounts
	CurrencyID        uint      `gorm:"not null" json:"currency_id"`
	Tip               string    `gorm:"not null" json:"tip"`   // "tekuci" | "devizni"
	Vrsta             string    `gorm:"not null" json:"vrsta"` // "licni" | "poslovni"
	Stanje            float64   `gorm:"default:0" json:"stanje"`
	RaspolozivoStanje float64   `gorm:"default:0" json:"raspolozivo_stanje"`
	DnevniLimit       float64   `gorm:"default:100000" json:"dnevni_limit"`
	MesecniLimit      float64   `gorm:"default:1000000" json:"mesecni_limit"`
	Naziv             string    `json:"naziv"`
	Status            string    `gorm:"default:'aktivan'" json:"status"` // aktivan | blokiran | zatvoren
	Client            *Client   `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	Firma             *Firma    `gorm:"foreignKey:FirmaID" json:"firma,omitempty"`
	Currency          Currency  `gorm:"foreignKey:CurrencyID" json:"currency"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
