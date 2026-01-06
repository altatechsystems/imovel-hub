package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: test-xls <file.xls>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Check if it's HTML
	isHTML, err := isHTMLFile(filePath)
	if err != nil {
		fmt.Printf("Error checking file type: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📄 File: %s\n", filePath)
	fmt.Printf("📄 Is HTML: %v\n\n", isHTML)

	if !isHTML {
		fmt.Println("❌ File is not HTML format")
		os.Exit(1)
	}

	// Read and display file content
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📄 File size: %d bytes\n", len(content))
	fmt.Printf("📄 First 500 chars:\n%s\n\n", string(content[:min(500, len(content))]))

	// Parse HTML
	file.Seek(0, 0)
	doc, err := html.Parse(file)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		os.Exit(1)
	}

	// Extract table rows
	rows := extractTableRows(doc)
	fmt.Printf("\n📊 Summary: found %d rows\n\n", len(rows))

	if len(rows) > 0 {
		fmt.Println("First 3 rows:")
		for i, row := range rows {
			if i >= 3 {
				break
			}
			fmt.Printf("  Row %d (%d cells): %v\n", i+1, len(row), row)
		}
	}
}

func isHTMLFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}

	content := strings.ToLower(string(buf[:n]))
	return strings.Contains(content, "<html") ||
		strings.Contains(content, "<!doctype") ||
		strings.Contains(content, "<table"), nil
}

func extractTableRows(n *html.Node) [][]string {
	var rows [][]string
	var foundTables int
	var foundTRs int
	var f func(*html.Node, int)

	f = func(n *html.Node, depth int) {
		if n.Type == html.ElementNode {
			if n.Data == "table" {
				foundTables++
				fmt.Printf("🔍 Found <table> at depth %d\n", depth)
			}
			if n.Data == "tr" {
				foundTRs++
				row := extractTableCells(n)
				fmt.Printf("🔍 Found <tr> at depth %d with %d cells\n", depth, len(row))
				if len(row) > 0 {
					rows = append(rows, row)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, depth+1)
		}
	}

	f(n, 0)
	fmt.Printf("📊 Found %d tables, %d <tr> elements\n", foundTables, foundTRs)
	return rows
}

func extractTableCells(tr *html.Node) []string {
	var cells []string

	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
			cells = append(cells, getNodeText(c))
		}
	}

	return cells
}

func getNodeText(n *html.Node) string {
	var text strings.Builder

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return strings.TrimSpace(text.String())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
