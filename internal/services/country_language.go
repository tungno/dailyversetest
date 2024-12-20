/**
 *  CountryLanguageMap provides a mapping of country names to their ISO 3166-1 alpha-2
 *  country codes and primary language codes. It enables efficient retrieval of these
 *  codes for use in applications like news APIs and localization services.
 *
 *  @map       CountryLanguageMap
 *  @methods
 *  - GetCountryAndLanguageCode(countryName)  - Retrieves the country code and primary language code for a given country.
 *
 *  @dependencies
 *  - strings.Title: Used to normalize country names for case-insensitive matching.
 *  - fmt.Errorf: Provides formatted error messages for unmatched countries.
 *
 *  @file      country_language.go
 *  @project   DailyVerse
 *  @purpose   Country and language mapping utility for localization and API integration.
 *  @framework Go Standard Library
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package services

import (
	"fmt"
	"strings"
)

// CountryLanguageMap maps country names to their two-letter ISO country codes and primary two-letter language codes.
var CountryLanguageMap = map[string]struct {
	CountryCode  string
	LanguageCode string
}{
	"Afghanistan":                      {"AF", "fa"},
	"Albania":                          {"AL", "sq"},
	"Algeria":                          {"DZ", "ar"},
	"Andorra":                          {"AD", "ca"},
	"Angola":                           {"AO", "pt"},
	"Argentina":                        {"AR", "es"},
	"Armenia":                          {"AM", "hy"},
	"Australia":                        {"AU", "en"},
	"Austria":                          {"AT", "de"},
	"Azerbaijan":                       {"AZ", "az"},
	"Bahamas":                          {"BS", "en"},
	"Bahrain":                          {"BH", "ar"},
	"Bangladesh":                       {"BD", "bn"},
	"Belarus":                          {"BY", "be"},
	"Belgium":                          {"BE", "nl"}, // Also "fr" and "de"
	"Belize":                           {"BZ", "en"},
	"Benin":                            {"BJ", "fr"},
	"Bhutan":                           {"BT", "dz"},
	"Bolivia":                          {"BO", "es"},
	"Bosnia and Herzegovina":           {"BA", "bs"},
	"Botswana":                         {"BW", "en"},
	"Brazil":                           {"BR", "pt"},
	"Brunei":                           {"BN", "ms"},
	"Bulgaria":                         {"BG", "bg"},
	"Burkina Faso":                     {"BF", "fr"},
	"Burundi":                          {"BI", "fr"},
	"Cambodia":                         {"KH", "km"},
	"Cameroon":                         {"CM", "fr"},
	"Canada":                           {"CA", "en"}, // Also "fr"
	"Cape Verde":                       {"CV", "pt"},
	"Central African Republic":         {"CF", "fr"},
	"Chad":                             {"TD", "fr"},
	"Chile":                            {"CL", "es"},
	"China":                            {"CN", "zh"},
	"Colombia":                         {"CO", "es"},
	"Comoros":                          {"KM", "ar"},
	"Congo (Congo-Brazzaville)":        {"CG", "fr"},
	"Congo (Democratic Republic)":      {"CD", "fr"},
	"Costa Rica":                       {"CR", "es"},
	"Croatia":                          {"HR", "hr"},
	"Cuba":                             {"CU", "es"},
	"Cyprus":                           {"CY", "el"},
	"Czech Republic":                   {"CZ", "cs"},
	"Denmark":                          {"DK", "da"},
	"Djibouti":                         {"DJ", "fr"},
	"Dominica":                         {"DM", "en"},
	"Dominican Republic":               {"DO", "es"},
	"Ecuador":                          {"EC", "es"},
	"Egypt":                            {"EG", "ar"},
	"El Salvador":                      {"SV", "es"},
	"Equatorial Guinea":                {"GQ", "es"}, // Also "fr" and "pt"
	"Eritrea":                          {"ER", "ti"},
	"Estonia":                          {"EE", "et"},
	"Eswatini":                         {"SZ", "en"}, // Also "ss"
	"Ethiopia":                         {"ET", "am"},
	"Fiji":                             {"FJ", "en"},
	"Finland":                          {"FI", "fi"},
	"France":                           {"FR", "fr"},
	"Gabon":                            {"GA", "fr"},
	"Gambia":                           {"GM", "en"},
	"Georgia":                          {"GE", "ka"},
	"Germany":                          {"DE", "de"},
	"Ghana":                            {"GH", "en"},
	"Greece":                           {"GR", "el"},
	"Grenada":                          {"GD", "en"},
	"Guatemala":                        {"GT", "es"},
	"Guinea":                           {"GN", "fr"},
	"Guinea-Bissau":                    {"GW", "pt"},
	"Guyana":                           {"GY", "en"},
	"Haiti":                            {"HT", "fr"}, // Also "ht"
	"Honduras":                         {"HN", "es"},
	"Hungary":                          {"HU", "hu"},
	"Iceland":                          {"IS", "is"},
	"India":                            {"IN", "hi"}, // Also "en"
	"Indonesia":                        {"ID", "id"},
	"Iran":                             {"IR", "fa"},
	"Iraq":                             {"IQ", "ar"},
	"Ireland":                          {"IE", "en"},
	"Italy":                            {"IT", "it"},
	"Jamaica":                          {"JM", "en"},
	"Japan":                            {"JP", "ja"},
	"Jordan":                           {"JO", "ar"},
	"Kazakhstan":                       {"KZ", "kk"},
	"Kenya":                            {"KE", "en"}, // Also "sw"
	"Kiribati":                         {"KI", "en"},
	"Kuwait":                           {"KW", "ar"},
	"Kyrgyzstan":                       {"KG", "ky"},
	"Laos":                             {"LA", "lo"},
	"Latvia":                           {"LV", "lv"},
	"Lebanon":                          {"LB", "ar"},
	"Lesotho":                          {"LS", "en"},
	"Liberia":                          {"LR", "en"},
	"Libya":                            {"LY", "ar"},
	"Liechtenstein":                    {"LI", "de"},
	"Lithuania":                        {"LT", "lt"},
	"Luxembourg":                       {"LU", "fr"}, // Also "de" and "lb"
	"Madagascar":                       {"MG", "fr"},
	"Malawi":                           {"MW", "en"},
	"Malaysia":                         {"MY", "ms"},
	"Maldives":                         {"MV", "dv"},
	"Mali":                             {"ML", "fr"},
	"Malta":                            {"MT", "mt"},
	"Marshall Islands":                 {"MH", "en"},
	"Mauritania":                       {"MR", "ar"},
	"Mauritius":                        {"MU", "en"},
	"Mexico":                           {"MX", "es"},
	"Micronesia":                       {"FM", "en"},
	"Moldova":                          {"MD", "ro"},
	"Monaco":                           {"MC", "fr"},
	"Mongolia":                         {"MN", "mn"},
	"Montenegro":                       {"ME", "sr"},
	"Morocco":                          {"MA", "ar"},
	"Mozambique":                       {"MZ", "pt"},
	"Myanmar":                          {"MM", "my"},
	"Namibia":                          {"NA", "en"},
	"Nauru":                            {"NR", "en"},
	"Nepal":                            {"NP", "ne"},
	"Netherlands":                      {"NL", "nl"},
	"New Zealand":                      {"NZ", "en"},
	"Nicaragua":                        {"NI", "es"},
	"Niger":                            {"NE", "fr"},
	"Nigeria":                          {"NG", "en"},
	"North Korea":                      {"KP", "ko"},
	"North Macedonia":                  {"MK", "mk"},
	"Norway":                           {"NO", "no"},
	"Oman":                             {"OM", "ar"},
	"Pakistan":                         {"PK", "ur"},
	"Palau":                            {"PW", "en"},
	"Palestine":                        {"PS", "ar"},
	"Panama":                           {"PA", "es"},
	"Papua New Guinea":                 {"PG", "en"},
	"Paraguay":                         {"PY", "es"},
	"Peru":                             {"PE", "es"},
	"Philippines":                      {"PH", "en"},
	"Poland":                           {"PL", "pl"},
	"Portugal":                         {"PT", "pt"},
	"Qatar":                            {"QA", "ar"},
	"Romania":                          {"RO", "ro"},
	"Russia":                           {"RU", "ru"},
	"Rwanda":                           {"RW", "rw"},
	"Saint Kitts and Nevis":            {"KN", "en"},
	"Saint Lucia":                      {"LC", "en"},
	"Saint Vincent and the Grenadines": {"VC", "en"},
	"Samoa":                            {"WS", "sm"},
	"San Marino":                       {"SM", "it"},
	"Saudi Arabia":                     {"SA", "ar"},
	"Senegal":                          {"SN", "fr"},
	"Serbia":                           {"RS", "sr"},
	"Seychelles":                       {"SC", "fr"},
	"Sierra Leone":                     {"SL", "en"},
	"Singapore":                        {"SG", "en"},
	"Slovakia":                         {"SK", "sk"},
	"Slovenia":                         {"SI", "sl"},
	"Solomon Islands":                  {"SB", "en"},
	"Somalia":                          {"SO", "so"},
	"South Africa":                     {"ZA", "en"}, // Also "af"
	"South Korea":                      {"KR", "ko"},
	"South Sudan":                      {"SS", "en"},
	"Spain":                            {"ES", "es"},
	"Sri Lanka":                        {"LK", "si"},
	"Sudan":                            {"SD", "ar"},
	"Suriname":                         {"SR", "nl"},
	"Sweden":                           {"SE", "sv"},
	"Switzerland":                      {"CH", "de"}, // Also "fr", "it", "rm"
	"Syria":                            {"SY", "ar"},
	"Taiwan":                           {"TW", "zh"},
	"Tajikistan":                       {"TJ", "tg"},
	"Tanzania":                         {"TZ", "sw"},
	"Thailand":                         {"TH", "th"},
	"Togo":                             {"TG", "fr"},
	"Tonga":                            {"TO", "to"},
	"Trinidad and Tobago":              {"TT", "en"},
	"Tunisia":                          {"TN", "ar"},
	"Turkey":                           {"TR", "tr"},
	"Turkmenistan":                     {"TM", "tk"},
	"Tuvalu":                           {"TV", "en"},
	"Uganda":                           {"UG", "en"},
	"Ukraine":                          {"UA", "uk"},
	"United Arab Emirates":             {"AE", "ar"},
	"United Kingdom":                   {"GB", "en"},
	"United States":                    {"US", "en"},
	"Uruguay":                          {"UY", "es"},
	"Uzbekistan":                       {"UZ", "uz"},
	"Vanuatu":                          {"VU", "en"},
	"Vatican City":                     {"VA", "it"},
	"Venezuela":                        {"VE", "es"},
	"Vietnam":                          {"VN", "vi"},
	"Yemen":                            {"YE", "ar"},
	"Zambia":                           {"ZM", "en"},
	"Zimbabwe":                         {"ZW", "en"},
}

// GetCountryAndLanguageCode retrieves the country code and primary language code for a given country name.
// Parameters:
//   - countryName (string): The name of the country (case-insensitive).
//
// Returns:
//   - string: ISO country code (e.g., "US" for the United States).
//   - string: Primary language code (e.g., "en" for English).
//   - error: Returns an error if the country is not found in the map.
func GetCountryAndLanguageCode(countryName string) (string, string, error) {
	// Normalize the country name for case-insensitive matching.
	normalizedCountryName := strings.Title(strings.ToLower(countryName))

	// Retrieve the country and language codes from the map.
	if entry, exists := CountryLanguageMap[normalizedCountryName]; exists {
		return strings.ToLower(entry.CountryCode), strings.ToLower(entry.LanguageCode), nil
	}

	// Return an error if the country is not found.
	return "", "", fmt.Errorf("country not found in map: %s", countryName)
}
