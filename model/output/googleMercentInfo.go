package output

import (
	"encoding/json"
	"encoding/xml"
	"os"
)

// Item represents a single product in the Google Merchant format
type Item struct {
	XMLName      xml.Name `xml:"item"`
	ID           string   `xml:"g:id"`
	Title        string   `xml:"g:title"`
	Description  string   `xml:"g:description"`
	Link         string   `xml:"g:link"`
	ImageLink    string   `xml:"g:image_link"`
	Brand        string   `xml:"g:brand"`
	Price        string   `xml:"g:price"`
	Availability string   `xml:"g:availability"`
	GTIN         string   `xml:"g:gtin"` // GTIN is for product identification
}

// Channel represents the channel information and items
type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Items       []Item   `xml:"item"`
}

// RSS is the root structure of the feed
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	NS      string   `xml:"xmlns:g,attr"`
	Channel Channel  `xml:"channel"`
}

// AdItem is the internal representation of an ad item that gets converted to XML
type AdItem struct {
	ID           string      `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	Link         string      `json:"link"`
	ImageLink    string      `json:"image_link"`
	Brand        string      `json:"brand"`
	Price        string      `json:"price"`
	Availability string      `json:"availability"`
	GTIN         json.Number `json:"gtin"`
}

// WriteRSSFeedToFile generates and writes the RSS feed to an XML file
func WriteRSSFeedToFile(adItems []AdItem) error {
	// Convert []AdItem to []Item (XML format)
	var items []Item
	for _, ad := range adItems {
		items = append(items, Item{
			ID:           ad.ID,
			Title:        ad.Title,
			Description:  ad.Description,
			Link:         ad.Link,
			ImageLink:    ad.ImageLink,
			Brand:        ad.Brand,
			Price:        ad.Price,
			Availability: ad.Availability,
			GTIN:         ad.GTIN.String(),
		})
	}

	rssFeed := RSS{
		Version: "2.0",
		NS:      "http://base.google.com/ns/1.0", // Google Merchant namespace
		Channel: Channel{
			Title:       "Ayshei",              // Your shop name
			Link:        "https://ayshei.com/", // Your shop link
			Description: "Your one-stop shop for the latest fashion items",
			Items:       items, // List of products
		},
	}

	// Marshal the RSS feed into XML with indentation for readability
	xmlOutput, err := xml.MarshalIndent(rssFeed, "", "  ")
	if err != nil {
		return err
	}

	// Create or overwrite the XML file
	file, err := os.Create("productsfashionaccessories.xml")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the XML header and the generated XML content
	_, err = file.WriteString(xml.Header + string(xmlOutput))
	if err != nil {
		return err
	}

	return nil
}
