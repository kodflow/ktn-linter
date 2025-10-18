package func001

// goodMixedCaps utilise MixedCaps.
func goodMixedCaps() {}

// GoodExported est une fonction exportée avec MixedCaps.
func GoodExported() {}

// parseHTTPRequest utilise un initialisme.
func parseHTTPRequest() {}

// HTTPServer est un initialisme complet.
func HTTPServer() {}

// bad_snake_case utilise snake_case (incorrect).
func bad_snake_case() {} // want "KTN-FUNC-001.*MixedCaps"

// Bad_Snake_Case mélange majuscules et underscores (incorrect).
func Bad_Snake_Case() {} // want "KTN-FUNC-001.*MixedCaps"

// calculate_total utilise snake_case (incorrect).
func calculate_total() {} // want "KTN-FUNC-001.*MixedCaps"
