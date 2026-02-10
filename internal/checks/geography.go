package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cultpedia/internal/models"
	"cultpedia/internal/utils"
)

func ValidateGeography() error {
	var allErrors []string
	var warnings []string

	if err := ValidateCountries(); err != nil {
		allErrors = append(allErrors, fmt.Sprintf("Countries validation failed:\n%v", err))
	}

	if err := ValidateContinents(); err != nil {
		allErrors = append(allErrors, fmt.Sprintf("Continents validation failed:\n%v", err))
	}

	if err := ValidateRegions(); err != nil {
		allErrors = append(allErrors, fmt.Sprintf("Regions validation failed:\n%v", err))
	}

	if err := ValidateFlags(); err != nil {
		warnings = append(warnings, fmt.Sprintf("Flags warning:\n%v", err))
	}

	if len(warnings) > 0 {
		fmt.Println(strings.Join(warnings, "\n\n"))
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(allErrors, "\n\n"))
	}

	return nil
}

func ValidateCountries() error {
	countries, err := utils.LoadCountries()
	if err != nil {
		return err
	}

	slugs := make(map[string]bool)
	isoAlpha2s := make(map[string]bool)
	isoAlpha3s := make(map[string]bool)
	var errors []string

	for i, c := range countries {
		if err := validateCountry(c); err != nil {
			errors = append(errors, fmt.Sprintf("line %d (slug: %s): %v", i+1, c.Slug, err))
			continue
		}

		if slugs[c.Slug] {
			errors = append(errors, fmt.Sprintf("duplicate slug '%s' at line %d", c.Slug, i+1))
		} else {
			slugs[c.Slug] = true
		}

		if isoAlpha2s[c.ISOAlpha2] {
			errors = append(errors, fmt.Sprintf("duplicate iso_alpha2 '%s' at line %d", c.ISOAlpha2, i+1))
		} else {
			isoAlpha2s[c.ISOAlpha2] = true
		}

		if isoAlpha3s[c.ISOAlpha3] {
			errors = append(errors, fmt.Sprintf("duplicate iso_alpha3 '%s' at line %d", c.ISOAlpha3, i+1))
		} else {
			isoAlpha3s[c.ISOAlpha3] = true
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func validateCountry(c models.Country) error {
	if c.Slug == "" {
		return fmt.Errorf("slug is required")
	}
	if c.ISOAlpha2 == "" {
		return fmt.Errorf("iso_alpha2 is required")
	}
	if len(c.ISOAlpha2) != 2 {
		return fmt.Errorf("iso_alpha2 must be 2 characters (got '%s')", c.ISOAlpha2)
	}
	if c.ISOAlpha3 == "" {
		return fmt.Errorf("iso_alpha3 is required")
	}
	if len(c.ISOAlpha3) != 3 {
		return fmt.Errorf("iso_alpha3 must be 3 characters (got '%s')", c.ISOAlpha3)
	}

	requiredLangs := []string{"en", "fr", "es"}
	for _, lang := range requiredLangs {
		if c.Name[lang] == "" {
			return fmt.Errorf("name.%s is required", lang)
		}
		if c.OfficialName[lang] == "" {
			return fmt.Errorf("official_name.%s is required", lang)
		}
	}

	if c.Continent == "" {
		return fmt.Errorf("continent is required")
	}

	if c.Coordinates.Lat < -90 || c.Coordinates.Lat > 90 {
		return fmt.Errorf("coordinates.lat must be between -90 and 90 (got %.6f)", c.Coordinates.Lat)
	}
	if c.Coordinates.Lng < -180 || c.Coordinates.Lng > 180 {
		return fmt.Errorf("coordinates.lng must be between -180 and 180 (got %.6f)", c.Coordinates.Lng)
	}

	if c.Currency.Code != "" {
		for _, lang := range requiredLangs {
			if c.Currency.Name[lang] == "" {
				return fmt.Errorf("currency.name.%s is required", lang)
			}
		}
	}

	validDrivingSides := []string{"left", "right"}
	if !contains(validDrivingSides, c.DrivingSide) {
		return fmt.Errorf("driving_side must be 'left' or 'right' (got '%s')", c.DrivingSide)
	}

	return nil
}

func ValidateContinents() error {
	continents, err := utils.LoadContinents()
	if err != nil {
		return err
	}

	ids := make(map[string]bool)
	var errors []string

	for i, c := range continents {
		if err := validateContinent(c); err != nil {
			errors = append(errors, fmt.Sprintf("line %d (slug: %s): %v", i+1, c.Slug, err))
			continue
		}

		if ids[c.Slug] {
			errors = append(errors, fmt.Sprintf("duplicate slug '%s' at line %d", c.Slug, i+1))
		} else {
			ids[c.Slug] = true
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func validateContinent(c models.Continent) error {
	if c.Slug == "" {
		return fmt.Errorf("slug is required")
	}

	requiredLangs := []string{"en", "fr", "es"}
	for _, lang := range requiredLangs {
		if c.Name[lang] == "" {
			return fmt.Errorf("name.%s is required", lang)
		}
	}

	if len(c.Countries) == 0 {
		return fmt.Errorf("countries list cannot be empty")
	}

	return nil
}

func ValidateRegions() error {
	regions, err := utils.LoadRegions()
	if err != nil {
		return err
	}

	ids := make(map[string]bool)
	var errors []string

	for i, r := range regions {
		if err := validateRegion(r); err != nil {
			errors = append(errors, fmt.Sprintf("line %d (slug: %s): %v", i+1, r.Slug, err))
			continue
		}

		if ids[r.Slug] {
			errors = append(errors, fmt.Sprintf("duplicate slug '%s' at line %d", r.Slug, i+1))
		} else {
			ids[r.Slug] = true
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func validateRegion(r models.Region) error {
	if r.Slug == "" {
		return fmt.Errorf("slug is required")
	}

	requiredLangs := []string{"en", "fr", "es"}
	for _, lang := range requiredLangs {
		if r.Name[lang] == "" {
			return fmt.Errorf("name.%s is required", lang)
		}
	}

	if r.Continent == "" {
		return fmt.Errorf("continent is required")
	}

	if len(r.Countries) == 0 {
		return fmt.Errorf("countries list cannot be empty")
	}

	return nil
}

func ValidateFlags() error {
	countries, err := utils.LoadCountries()
	if err != nil {
		return err
	}

	var errors []string

	for _, c := range countries {
		if c.Flag == "" {
			continue
		}

		flagPath := filepath.Join(utils.FlagsSVGDir, c.Flag+".svg")
		if _, err := os.Stat(flagPath); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("missing flag file for %s: %s", c.Slug, flagPath))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("missing flags:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func CheckGeographyDuplicates() string {
	var results []string

	countries, err := utils.LoadCountries()
	if err != nil {
		return fmt.Sprintf("✗ Error loading countries: %v", err)
	}

	slugs := make(map[string]int)
	for i, c := range countries {
		if firstLine, exists := slugs[c.Slug]; exists {
			results = append(results, fmt.Sprintf("Duplicate country slug '%s' found at lines %d and %d", c.Slug, firstLine, i+1))
		} else {
			slugs[c.Slug] = i + 1
		}
	}

	if len(results) == 0 {
		return "✔ No duplicates detected in geography dataset"
	}

	return fmt.Sprintf("✗ Duplicates detected:\n%s", strings.Join(results, "\n"))
}

func CheckGeographyTranslations() string {
	var missing []string
	requiredLangs := []string{"en", "fr", "es"}

	countries, err := utils.LoadCountries()
	if err != nil {
		return fmt.Sprintf("✗ Error loading countries: %v", err)
	}

	for i, c := range countries {
		for _, lang := range requiredLangs {
			if c.Name[lang] == "" {
				missing = append(missing, fmt.Sprintf("Country line %d (slug: %s): missing %s translation for 'name'", i+1, c.Slug, lang))
			}
			if c.OfficialName[lang] == "" {
				missing = append(missing, fmt.Sprintf("Country line %d (slug: %s): missing %s translation for 'official_name'", i+1, c.Slug, lang))
			}
			if c.Currency.Code != "" && c.Currency.Name[lang] == "" {
				missing = append(missing, fmt.Sprintf("Country line %d (slug: %s): missing %s translation for 'currency.name'", i+1, c.Slug, lang))
			}
		}
	}

	continents, err := utils.LoadContinents()
	if err != nil {
		return fmt.Sprintf("✗ Error loading continents: %v", err)
	}

	for i, c := range continents {
		for _, lang := range requiredLangs {
			if c.Name[lang] == "" {
				missing = append(missing, fmt.Sprintf("Continent line %d (slug: %s): missing %s translation for 'name'", i+1, c.Slug, lang))
			}
		}
	}

	regions, err := utils.LoadRegions()
	if err != nil {
		return fmt.Sprintf("✗ Error loading regions: %v", err)
	}

	for i, r := range regions {
		for _, lang := range requiredLangs {
			if r.Name[lang] == "" {
				missing = append(missing, fmt.Sprintf("Region line %d (slug: %s): missing %s translation for 'name'", i+1, r.Slug, lang))
			}
		}
	}

	if len(missing) == 0 {
		return "✔ All geography translations present"
	}

	return fmt.Sprintf("✗ Missing translations:\n%s", strings.Join(missing, "\n"))
}
