package input

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/machinebox/graphql"
)

// AdItem represents the structure for storing ad information
type AdItem struct {
	ID           string
	Title        string
	Description  string
	Link         string
	ImageLink    string
	Brand        string
	Price        string
	Availability string
	CodeNumber   json.Number // Handle GTIN as json.Number
}

// AdAttributes represents the structure of attributes for each ad
type AdAttributes struct {
	StepsData []struct {
		Name string `json:"name"`
		Data struct {
			ID struct {
				ID    string `json:"id"`
				Value string `json:"value"`
			} `json:"id"`
			InputSearchValue struct {
				Value string `json:"value"`
			} `json:"inputSearchValue"`
			Values struct {
				Brand  string `json:"brand"`
				Price  string `json:"price"`
				Images []struct {
					Src string `json:"src"`
				} `json:"images"`
				AdType string `json:"ad_type"`
			} `json:"values"`
			PaymentMethods struct {
				Data []struct {
					Value string `json:"value"`
				} `json:"data"`
			} `json:"paymentMethods"`
		} `json:"data"`
	} `json:"stepsData"`
}

func FetchAds(endpoint, adminSecret string) ([]AdItem, error) {
	client := graphql.NewClient(endpoint)

	// GraphQL query with status, category, and payment method filter
	req := graphql.NewRequest(`
	query ($last24Hours: timestamptz!) {
		ads(where: {
			status: {_eq: "Published"},
			category_id: {_eq: "e87e7959-03ef-4bd1-930d-4a96c5743108"},
			updated_at: { _gte: $last24Hours }
		}) {
			id
			draft_id
			description
			attributes
			code_number
		}
	}
`)

	// Calculate the timestamp for the last 24 hours
	last24Hours := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	req.Var("last24Hours", last24Hours)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hasura-Admin-Secret", adminSecret)

	var response struct {
		Ads []struct {
			ID          string          `json:"id"`
			DraftID     string          `json:"draft_id"`
			Description string          `json:"description"`
			CodeNumber  json.Number     `json:"code_number"`
			Attributes  json.RawMessage `json:"attributes"`
		} `json:"ads"`
	}

	err := client.Run(context.Background(), req, &response)
	if err != nil {
		return nil, err
	}

	var items []AdItem
	auctionCount := 0
	otherCount := 0

	// Process the attributes of each ad
	for _, ad := range response.Ads {
		var attrs AdAttributes
		err := json.Unmarshal(ad.Attributes, &attrs)
		if err != nil {
			log.Printf("Error unmarshalling attributes for ad ID %s: %v", ad.ID, err)
			continue
		}

		// Check for specific subcategories
		shouldInclude := false
		allowedSubcategories := map[string]bool{
			"212818c2-5ae3-4a95-88c9-370b3b906df0": true,
			"456ceaaa-de4d-449f-8621-3af7253fe452": true,
			"5feb2aa4-3361-401b-ab05-d0623bab291b": true,
			"7685d106-a4dd-48ed-876b-4dd8116f114c": true,
			"34991934-f9ef-457c-9824-c82dad366889": true,
			"1c4df47a-e94a-49b4-aeea-1d77dc4f5458": true,
			"e84fd5e8-c303-46db-b1c6-e493781aef40": true,
			"63d47c2b-a5eb-4439-b45d-ccbaa4ca671a": true,
			"73a17eb3-1686-40d6-bcce-edc5c69b5540": true,
		}

		for _, step := range attrs.StepsData {
			if step.Name == "search_product" {
				if allowedSubcategories[step.Data.ID.ID] {
					shouldInclude = true
					break
				}
			}
		}

		if !shouldInclude {
			continue
		}

		adType := ""
		price := ""
		hasOnlinePayment := false
		for _, step := range attrs.StepsData {
			if step.Name == "delivery_and_payment_methods" {
				for _, payment := range step.Data.PaymentMethods.Data {
					if payment.Value == "Online Payment" {
						hasOnlinePayment = true
					}
				}
			} else if step.Name == "product_detail" {
				adType = step.Data.Values.AdType
				price = step.Data.Values.Price
			}
		}

		// Count ad types
		if adType == "auction" {
			auctionCount++
		} else {
			otherCount++
		}

		// If "Online Payment" is found and no "Auctions" in attributes, process the ad
		if hasOnlinePayment {
			// Extract title, brand, and image src from attributes
			title, brand, imageSrc := "", "", ""
			for _, step := range attrs.StepsData {
				if step.Name == "search_product" {
					title = step.Data.InputSearchValue.Value
				} else if step.Name == "product_detail" {
					brand = step.Data.Values.Brand
					if len(step.Data.Values.Images) > 0 {
						imageSrc = step.Data.Values.Images[0].Src
					}
				}
			}

			// Ensure that `imageSrc` is properly formatted without encoding issues
			if imageSrc != "" {
				imageSrc = fmt.Sprintf(
					"https://ayshei.com/_next/image?url=https://storage.ayshei.com/prod/public/drafts/%s/web/%s&amp;w=3840&amp;q=75",
					ad.DraftID, imageSrc)
			}

			// Skip items with empty CodeNumber
			if ad.CodeNumber == "" {
				log.Printf("Skipping ad %s due to missing code_number", ad.ID)
				continue
			}

			// Clean up description by removing U+200E character
			description := strings.ReplaceAll(ad.Description, "\u200E", "")

			// Clean up title by removing '&' symbol
			title = strings.ReplaceAll(title, "&", "")

			// Build the AdItem
			items = append(items, AdItem{
				ID:           ad.ID,
				Title:        title,
				Description:  description,
				Link:         fmt.Sprintf("https://ayshei.com/product/%s", ad.ID),
				ImageLink:    imageSrc,
				Brand:        brand,
				Price:        price + " AED",
				Availability: "in stock",
				CodeNumber:   ad.CodeNumber,
			})
		}
	}

	// Log counts of "auction" and other ad types
	log.Printf("Total ads with ad_type 'auction': %d", auctionCount)
	log.Printf("Total ads with other ad types: %d", otherCount)

	return items, nil
}
