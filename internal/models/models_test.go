package models_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/models"
)

// --- Instantiation tests ---

func TestCurrency_Instantiation(t *testing.T) {
	c := models.Currency{
		Kod:    "EUR",
		Naziv:  "Euro",
		Simbol: "€",
		Drzava: "European Union",
	}
	if c.Kod != "EUR" {
		t.Errorf("Currency.Kod = %q, want %q", c.Kod, "EUR")
	}
	if c.Naziv != "Euro" {
		t.Errorf("Currency.Naziv = %q, want %q", c.Naziv, "Euro")
	}
	// Note: Aktivan default:true is a DB-level default, not a Go zero-value
}

func TestSifraDelatnosti_Instantiation(t *testing.T) {
	s := models.SifraDelatnosti{
		Sifra: "6419",
		Naziv: "Monetary intermediation",
	}
	if s.Sifra != "6419" {
		t.Errorf("SifraDelatnosti.Sifra = %q, want %q", s.Sifra, "6419")
	}
	if s.Naziv == "" {
		t.Error("SifraDelatnosti.Naziv should not be empty")
	}
}

func TestSifraPlacanja_Instantiation(t *testing.T) {
	s := models.SifraPlacanja{
		Sifra: "221",
		Naziv: "Plaćanje robe",
	}
	if s.Sifra != "221" {
		t.Errorf("SifraPlacanja.Sifra = %q, want %q", s.Sifra, "221")
	}
}

func TestFirma_Instantiation(t *testing.T) {
	f := models.Firma{
		Naziv:       "EXBanka d.o.o.",
		MaticniBroj: "12345678",
		PIB:         "987654321",
		Adresa:      "Bulevar Kralja Aleksandra 1",
		Telefon:     "011123456",
	}
	if f.Naziv == "" {
		t.Error("Firma.Naziv should not be empty")
	}
	if f.MaticniBroj == "" {
		t.Error("Firma.MaticniBroj should not be empty")
	}
	if f.PIB == "" {
		t.Error("Firma.PIB should not be empty")
	}
}

func TestAccount_Instantiation(t *testing.T) {
	clientID := uint(1)
	currencyID := uint(1)
	a := models.Account{
		BrojRacuna:        "123456789012345678",
		ClientID:          &clientID,
		CurrencyID:        currencyID,
		Tip:               "tekuci",
		Vrsta:             "licni",
		Stanje:            5000.0,
		RaspolozivoStanje: 5000.0,
		DnevniLimit:       100000,
		MesecniLimit:      1000000,
		Naziv:             "Moj tekući račun",
		Status:            "aktivan",
	}
	if a.BrojRacuna == "" {
		t.Error("Account.BrojRacuna should not be empty")
	}
	if a.Tip != "tekuci" && a.Tip != "devizni" {
		t.Errorf("Account.Tip = %q, want tekuci or devizni", a.Tip)
	}
	if a.Vrsta != "licni" && a.Vrsta != "poslovni" {
		t.Errorf("Account.Vrsta = %q, want licni or poslovni", a.Vrsta)
	}
}

func TestAccount_NullableClientID(t *testing.T) {
	// Bank-owned account has nil ClientID
	a := models.Account{
		BrojRacuna: "123456789012345678",
		CurrencyID: 1,
		Tip:        "tekuci",
		Vrsta:      "poslovni",
	}
	if a.ClientID != nil {
		t.Error("Account.ClientID should be nil for bank-owned accounts")
	}
}

func TestTransfer_Instantiation(t *testing.T) {
	now := time.Now()
	tr := models.Transfer{
		RacunPosiljaocaID: 1,
		RacunPrimaocaID:   2,
		Iznos:             1000.0,
		ValutaIznosa:      "RSD",
		Svrha:             "Uplata",
		Status:            "uspesno",
		VremeTransakcije:  now,
	}
	if tr.Iznos <= 0 {
		t.Error("Transfer.Iznos should be positive")
	}
	if tr.RacunPosiljaocaID == tr.RacunPrimaocaID {
		t.Error("Transfer sender and receiver should differ")
	}
}

func TestTransfer_StatusValues(t *testing.T) {
	validStatuses := []string{"uspesno", "neuspesno", "u_obradi"}
	for _, s := range validStatuses {
		tr := models.Transfer{Status: s}
		if tr.Status != s {
			t.Errorf("Transfer.Status = %q, want %q", tr.Status, s)
		}
	}
}

func TestPayment_Instantiation(t *testing.T) {
	now := time.Now()
	p := models.Payment{
		RacunPosiljaocaID: 1,
		RacunPrimaocaBroj: "987654321098765432",
		Iznos:             500.0,
		SifraPlacanja:     "221",
		PozivNaBroj:       "97-12345-678",
		Svrha:             "Kupovina",
		Status:            "u_obradi",
		VerifikacioniKod:  "123456",
		VremeTransakcije:  now,
	}
	if p.Iznos <= 0 {
		t.Error("Payment.Iznos should be positive")
	}
	if p.RacunPrimaocaBroj == "" {
		t.Error("Payment.RacunPrimaocaBroj should not be empty")
	}
}

func TestPayment_StatusValues(t *testing.T) {
	validStatuses := []string{"u_obradi", "uspesno", "neuspesno", "stornirano"}
	for _, s := range validStatuses {
		p := models.Payment{Status: s}
		if p.Status != s {
			t.Errorf("Payment.Status = %q, want %q", p.Status, s)
		}
	}
}

func TestPaymentRecipient_Instantiation(t *testing.T) {
	r := models.PaymentRecipient{
		ClientID:   1,
		Naziv:      "Elektrodistribucija",
		BrojRacuna: "111222333444555666",
	}
	if r.Naziv == "" {
		t.Error("PaymentRecipient.Naziv should not be empty")
	}
	if r.BrojRacuna == "" {
		t.Error("PaymentRecipient.BrojRacuna should not be empty")
	}
}

// --- GORM tag constraint tests ---

func getGormTag(t *testing.T, typ reflect.Type, fieldName string) string {
	t.Helper()
	f, ok := typ.FieldByName(fieldName)
	if !ok {
		t.Fatalf("field %q not found on type %s", fieldName, typ.Name())
	}
	return string(f.Tag.Get("gorm"))
}

func TestAccount_BrojRacuna_UniqueIndex(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Account{}), "BrojRacuna")
	if !strings.Contains(tag, "uniqueIndex") {
		t.Errorf("Account.BrojRacuna gorm tag = %q, want contains 'uniqueIndex'", tag)
	}
}

func TestAccount_BrojRacuna_Size18(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Account{}), "BrojRacuna")
	if !strings.Contains(tag, "size:18") {
		t.Errorf("Account.BrojRacuna gorm tag = %q, want contains 'size:18'", tag)
	}
}

func TestAccount_BrojRacuna_NotNull(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Account{}), "BrojRacuna")
	if !strings.Contains(tag, "not null") {
		t.Errorf("Account.BrojRacuna gorm tag = %q, want contains 'not null'", tag)
	}
}

func TestCurrency_Kod_UniqueIndex(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Currency{}), "Kod")
	if !strings.Contains(tag, "uniqueIndex") {
		t.Errorf("Currency.Kod gorm tag = %q, want contains 'uniqueIndex'", tag)
	}
}

func TestCurrency_Kod_Size3(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Currency{}), "Kod")
	if !strings.Contains(tag, "size:3") {
		t.Errorf("Currency.Kod gorm tag = %q, want contains 'size:3'", tag)
	}
}

func TestCurrency_Kod_NotNull(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Currency{}), "Kod")
	if !strings.Contains(tag, "not null") {
		t.Errorf("Currency.Kod gorm tag = %q, want contains 'not null'", tag)
	}
}

func TestFirma_MaticniBroj_UniqueIndex(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Firma{}), "MaticniBroj")
	if !strings.Contains(tag, "uniqueIndex") {
		t.Errorf("Firma.MaticniBroj gorm tag = %q, want contains 'uniqueIndex'", tag)
	}
}

func TestFirma_PIB_UniqueIndex(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Firma{}), "PIB")
	if !strings.Contains(tag, "uniqueIndex") {
		t.Errorf("Firma.PIB gorm tag = %q, want contains 'uniqueIndex'", tag)
	}
}

func TestTransfer_RacunPosiljaocaID_NotNull(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Transfer{}), "RacunPosiljaocaID")
	if !strings.Contains(tag, "not null") {
		t.Errorf("Transfer.RacunPosiljaocaID gorm tag = %q, want contains 'not null'", tag)
	}
}

func TestPayment_RacunPrimaocaBroj_NotNull(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Payment{}), "RacunPrimaocaBroj")
	if !strings.Contains(tag, "not null") {
		t.Errorf("Payment.RacunPrimaocaBroj gorm tag = %q, want contains 'not null'", tag)
	}
}

func TestPayment_Status_DefaultUObradi(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.Payment{}), "Status")
	if !strings.Contains(tag, "u_obradi") {
		t.Errorf("Payment.Status gorm tag = %q, want contains 'u_obradi' default", tag)
	}
}

func TestSifraDelatnosti_Sifra_UniqueIndex(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.SifraDelatnosti{}), "Sifra")
	if !strings.Contains(tag, "uniqueIndex") {
		t.Errorf("SifraDelatnosti.Sifra gorm tag = %q, want contains 'uniqueIndex'", tag)
	}
}

func TestSifraPlacanja_Sifra_UniqueIndex(t *testing.T) {
	tag := getGormTag(t, reflect.TypeOf(models.SifraPlacanja{}), "Sifra")
	if !strings.Contains(tag, "uniqueIndex") {
		t.Errorf("SifraPlacanja.Sifra gorm tag = %q, want contains 'uniqueIndex'", tag)
	}
}
