package rules_func

// ANTI-PATTERN: Fonction avec beaucoup trop de paramètres (> 5)
// Viole KTN-FUNC-001

// ProcessUserDataBad a 8 paramètres - TROP !
func ProcessUserDataBad(name string, email string, age int, country string, city string, zipcode string, phone string, newsletter bool) error {
	// Devrait utiliser un struct de configuration
	return nil
}

// CreateOrderBad a 10 paramètres - CATASTROPHIQUE !
func CreateOrderBad(userId int, productId int, quantity int, price float64, tax float64, shipping float64, discount float64, currency string, paymentMethod string, notes string) error {
	return nil
}

// SendEmailBad a 7 paramètres
func SendEmailBad(to string, from string, subject string, body string, cc string, bcc string, attachments []string) error {
	return nil
}

// GenerateReportBad a 9 paramètres
func GenerateReportBad(startDate string, endDate string, format string, includeCharts bool, includeTables bool, timezone string, locale string, currency string, decimals int) ([]byte, error) {
	return nil, nil
}
