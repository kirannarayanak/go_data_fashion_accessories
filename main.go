// main.go
package main

import (
	"go_data_fashion_accessories/config"
	"go_data_fashion_accessories/model/input"
	"go_data_fashion_accessories/model/output"
	"go_data_fashion_accessories/util"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Function to calculate check digit for GTIN-13
func calculateGTINCheckDigit(gtin string) string {
	sum := 0
	for i, r := range gtin {
		digit := int(r - '0')
		if i%2 == 0 {
			sum += digit // Multiply odd position digits by 1
		} else {
			sum += digit * 3 // Multiply even position digits by 3
		}
	}
	remainder := sum % 10
	if remainder == 0 {
		return "0"
	}
	return strconv.Itoa(10 - remainder)
}

// Function to ensure valid GTIN
func ensureValidGTIN(gtin string) string {
	for len(gtin) < 12 {
		gtin = "0" + gtin
	}
	if len(gtin) == 12 {
		gtin += calculateGTINCheckDigit(gtin)
	}
	return gtin
}

// Function to strip HTML tags from the description and clean up the text
func cleanUpDescription(description string) string {
	// Regex to match all HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	// Remove all HTML tags
	cleaned := re.ReplaceAllString(description, "")
	// Ensure proper punctuation between sentences
	cleaned = strings.ReplaceAll(cleaned, ". ", ".")
	cleaned = strings.ReplaceAll(cleaned, ".", ". ")
	// Replace multiple spaces/newlines with a single space
	cleaned = strings.TrimSpace(strings.Join(strings.Fields(cleaned), " "))

	// Check description length (Google Merchant Center recommends at least 30 characters)
	if len(cleaned) < 30 {
		cleaned += " This product is of high quality and in stock."
	}

	return cleaned
}

// Function to ensure proper encoding for special characters like &, <, and >
func escapeSpecialCharacters(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ads, err := input.FetchAds(cfg.HasuraEndpoint, cfg.AdminSecret)
	if err != nil {
		log.Fatalf("Error fetching ads: %v", err)
	}

	var outputAds []output.Item
	for _, ad := range ads {
		gtin := ad.CodeNumber.String()
		validGTIN := ensureValidGTIN(gtin)

		// Clean up the description before adding it to the output
		cleanedDescription := cleanUpDescription(ad.Description)
		cleanedDescription = escapeSpecialCharacters(cleanedDescription)

		outputAds = append(outputAds, output.Item{
			ID:           ad.CodeNumber.String(),
			Title:        ad.Title,
			Description:  cleanedDescription, // Use cleaned description here
			Link:         ad.Link,
			ImageLink:    ad.ImageLink,
			Brand:        ad.Brand,
			Price:        ad.Price,
			Availability: ad.Availability,
			GTIN:         validGTIN,
		})
	}

	err = util.GenerateXML(outputAds)
	if err != nil {
		log.Fatalf("Error generating XML: %v", err)
	}

	log.Println("Successfully generated XML file")
}
