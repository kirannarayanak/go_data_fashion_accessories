package util

import (
	"go_data_fashion_accessories/model/output"
	"log"
	"os"
)

// GenerateXML manually creates the XML file, bypassing &amp; issues
func GenerateXML(ads []output.Item) error {
	// Create the file
	file, err := os.Create("productsfashionaccessories.xml")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the XML header
	file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	file.WriteString("\n<rss version=\"2.0\" xmlns:g=\"http://base.google.com/ns/1.0\">\n")
	file.WriteString("  <channel>\n")
	file.WriteString("    <title>Ayshei</title>\n")
	file.WriteString("    <link>https://ayshei.com/</link>\n")
	file.WriteString("    <description>Your one-stop shop for the latest fashion items</description>\n")

	// Manually write the ad items to the file
	for _, ad := range ads {
		file.WriteString("    <item>\n")
		file.WriteString("      <g:id>" + ad.ID + "</g:id>\n")
		file.WriteString("      <g:title>" + ad.Title + "</g:title>\n")
		file.WriteString("      <g:description>" + ad.Description + "</g:description>\n")
		file.WriteString("      <g:link>" + ad.Link + "</g:link>\n")
		// Manually write the image link without escaping
		file.WriteString("      <g:image_link>" + ad.ImageLink + "</g:image_link>\n")
		file.WriteString("      <g:brand>" + ad.Brand + "</g:brand>\n")
		file.WriteString("      <g:price>" + ad.Price + "</g:price>\n")
		file.WriteString("      <g:availability>" + ad.Availability + "</g:availability>\n")
		file.WriteString("      <g:gtin>" + ad.GTIN + "</g:gtin>\n")
		file.WriteString("    </item>\n")
	}

	// Close the root elements
	file.WriteString("  </channel>\n")
	file.WriteString("</rss>\n")

	log.Println("Successfully written XML to file.")
	return nil
}
