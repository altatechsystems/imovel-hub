package main

import (
	"context"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"google.golang.org/api/option"
)

// XMLImovel represents property from XML file
type XMLImovel struct {
	CodigoImobiliaria    string  `xml:"CodigoImobiliaria"`
	Referencia           string  `xml:"Referencia"`
	Tipo                 string  `xml:"Tipo"`
	Finalidade           string  `xml:"Finalidade"`
	Categoria            string  `xml:"Categoria"`
	Titulo               string  `xml:"Titulo"`
	Anuncioparainternet  string  `xml:"Anuncioparainternet"`
	Valorvenda           float64 `xml:"Valorvenda"`
	Valorlocacao         float64 `xml:"Valorlocacao"`
	Valorcondominio      float64 `xml:"Valorcondominio"`
	Dormitorios          int     `xml:"Dormitorios"`
	Suite                int     `xml:"Suite"`
	Banheiro2            int     `xml:"Banheiro2"`
	Garagem              int     `xml:"Garagem"`
	Areacosntruida       float64 `xml:"Areacosntruida"`
	Areatotal            float64 `xml:"Areatotal"`
	Endereco             string  `xml:"Endereco"`
	Numero               string  `xml:"Numero"`
	Complemento          string  `xml:"Complemento"`
	Bairro               string  `xml:"Bairro"`
	Cidade               string  `xml:"Cidade"`
	UnidadeFederativa    string  `xml:"UnidadeFederativa"`
	CEP                  string  `xml:"CEP"`
	Venda                int     `xml:"Venda"`
	Locacao              int     `xml:"Locacao"`
	Aceitafinanciamento  int     `xml:"Aceitafinanciamento"`
	Permuta              int     `xml:"Permuta"`
	Codigoimovel         string  `xml:"Codigoimovel"`
	LinkImovelSite       string  `xml:"LinkImovelSite"`
	Cadastradoem         string  `xml:"Cadastradoem"`
	Atualizadoem         string  `xml:"Atualizadoem"`
	Captador             string  `xml:"Captador"`
	Latitude             string  `xml:"Latitude"`
	Longitude            string  `xml:"Longitude"`
	// Features
	Churrasqueira        int     `xml:"Churrasqueira"`
	Playground           int     `xml:"Playground"`
	Piscina              int     `xml:"Piscina"`
	Arcondicionado       int     `xml:"Arcondicionado"`
	Armariocozinha       int     `xml:"Armariocozinha"`
	Lavanderia           int     `xml:"Lavanderia"`
	Sacada               int     `xml:"Sacada"`
	Varanda              int     `xml:"Varanda"`
}

type XMLImoveis struct {
	XMLName  xml.Name     `xml:"Union"`
	Imoveis  []XMLImovel  `xml:"Imoveis>Imovel"`
}

// Property represents the Firestore property document
type Property struct {
	ID              string                 `firestore:"id" json:"id"`
	TenantID        string                 `firestore:"tenant_id" json:"tenant_id"`
	Reference       string                 `firestore:"reference" json:"reference"`
	Title           string                 `firestore:"title" json:"title"`
	Description     string                 `firestore:"description" json:"description"`
	Type            string                 `firestore:"type" json:"type"`
	Status          string                 `firestore:"status" json:"status"`
	Purpose         string                 `firestore:"purpose" json:"purpose"`
	Category        string                 `firestore:"category" json:"category"`
	SalePrice       float64                `firestore:"sale_price" json:"sale_price"`
	RentalPrice     float64                `firestore:"rental_price" json:"rental_price"`
	CondoFee        float64                `firestore:"condo_fee" json:"condo_fee"`
	Address         map[string]interface{} `firestore:"address" json:"address"`
	Features        map[string]interface{} `firestore:"features" json:"features"`
	OwnerID         string                 `firestore:"owner_id" json:"owner_id"`
	BrokerID        string                 `firestore:"broker_id,omitempty" json:"broker_id,omitempty"`
	Bedrooms        int                    `firestore:"bedrooms" json:"bedrooms"`
	Bathrooms       int                    `firestore:"bathrooms" json:"bathrooms"`
	Suites          int                    `firestore:"suites" json:"suites"`
	ParkingSpaces   int                    `firestore:"parking_spaces" json:"parking_spaces"`
	Area            float64                `firestore:"area" json:"area"`
	TotalArea       float64                `firestore:"total_area" json:"total_area"`
	YearBuilt       int                    `firestore:"year_built,omitempty" json:"year_built,omitempty"`
	AcceptsFinancing bool                  `firestore:"accepts_financing" json:"accepts_financing"`
	AcceptsTrade    bool                   `firestore:"accepts_trade" json:"accepts_trade"`
	ExternalID      string                 `firestore:"external_id,omitempty" json:"external_id,omitempty"`
	ExternalURL     string                 `firestore:"external_url,omitempty" json:"external_url,omitempty"`
	Slug            string                 `firestore:"slug" json:"slug"`
	CreatedAt       time.Time              `firestore:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `firestore:"updated_at" json:"updated_at"`
	CreatedBy       string                 `firestore:"created_by" json:"created_by"`
	UpdatedBy       string                 `firestore:"updated_by" json:"updated_by"`
}

func main() {
	xmlFile := flag.String("xml", "", "Path to XML file")
	xlsFile := flag.String("xls", "", "Path to XLS file")
	tenantID := flag.String("tenant", "", "Tenant ID")
	ownerID := flag.String("owner", "", "Owner ID")
	credentials := flag.String("creds", "backend/config/firebase-adminsdk.json", "Firebase credentials file")
	projectID := flag.String("project", "ecosistema-imob-dev", "Firebase project ID")
	database := flag.String("database", "imob-dev", "Firestore database name")
	dryRun := flag.Bool("dry-run", false, "Dry run (don't actually import)")
	limit := flag.Int("limit", 0, "Limit number of properties to import (0 = no limit)")

	flag.Parse()

	if *tenantID == "" || *ownerID == "" {
		log.Fatal("Tenant ID and Owner ID are required")
	}

	if *xmlFile == "" && *xlsFile == "" {
		log.Fatal("At least one file (XML or XLS) is required")
	}

	ctx := context.Background()

	// Initialize Firebase credentials
	opt := option.WithCredentialsFile(*credentials)

	// Initialize Firestore with named database
	client, err := firestore.NewClientWithDatabase(ctx, *projectID, *database, opt)
	if err != nil {
		log.Fatalf("Error creating Firestore client: %v", err)
	}
	defer client.Close()

	imported := 0
	errors := 0

	// Import from XML
	if *xmlFile != "" {
		log.Printf("Importing from XML file: %s", *xmlFile)
		count, errs := importFromXML(ctx, client, *xmlFile, *tenantID, *ownerID, *dryRun, *limit)
		imported += count
		errors += errs
	}

	// Import from XLS
	if *xlsFile != "" {
		log.Printf("Importing from XLS file: %s", *xlsFile)
		count, errs := importFromXLS(ctx, client, *xlsFile, *tenantID, *ownerID, *dryRun, *limit)
		imported += count
		errors += errs
	}

	log.Printf("\n=== Import Summary ===")
	log.Printf("Successfully imported: %d properties", imported)
	log.Printf("Errors: %d", errors)
	if *dryRun {
		log.Printf("DRY RUN: No data was actually imported")
	}
}

func importFromXML(ctx context.Context, client *firestore.Client, filePath, tenantID, ownerID string, dryRun bool, limit int) (int, int) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening XML file: %v", err)
		return 0, 1
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading XML file: %v", err)
		return 0, 1
	}

	var imoveis XMLImoveis
	if err := xml.Unmarshal(data, &imoveis); err != nil {
		log.Printf("Error parsing XML: %v", err)
		return 0, 1
	}

	total := len(imoveis.Imoveis)
	if limit > 0 && limit < total {
		total = limit
		imoveis.Imoveis = imoveis.Imoveis[:limit]
	}

	log.Printf("Found %d properties in XML, will import %d", len(imoveis.Imoveis), total)

	imported := 0
	errors := 0

	for i, imovel := range imoveis.Imoveis {
		property := convertXMLToProperty(imovel, tenantID, ownerID)

		if dryRun {
			log.Printf("[DRY RUN] Would import property %d/%d: %s - %s", i+1, len(imoveis.Imoveis), property.Reference, property.Title)
		} else {
			if err := saveProperty(ctx, client, property); err != nil {
				log.Printf("Error importing property %s: %v", property.Reference, err)
				errors++
			} else {
				log.Printf("Imported property %d/%d: %s - %s", i+1, len(imoveis.Imoveis), property.Reference, property.Title)
				imported++
			}
		}

		// Add a small delay to avoid rate limiting
		if !dryRun && i > 0 && i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return imported, errors
}

func importFromXLS(ctx context.Context, client *firestore.Client, filePath, tenantID, ownerID string, dryRun bool, limit int) (int, int) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("Error opening XLS file: %v", err)
		return 0, 1
	}
	defer f.Close()

	// Get first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		log.Printf("No sheets found in XLS file")
		return 0, 1
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		log.Printf("Error reading XLS rows: %v", err)
		return 0, 1
	}

	if len(rows) < 2 {
		log.Printf("XLS file has no data rows")
		return 0, 1
	}

	dataRows := rows[1:]
	total := len(dataRows)
	if limit > 0 && limit < total {
		total = limit
		dataRows = dataRows[:limit]
	}

	log.Printf("Found %d properties in XLS, will import %d", len(rows)-1, total)

	// Get header row
	headers := rows[0]

	imported := 0
	errors := 0

	for i, row := range dataRows {
		property := convertXLSToProperty(headers, row, tenantID, ownerID)

		if dryRun {
			log.Printf("[DRY RUN] Would import property %d/%d: %s - %s", i+1, len(rows)-1, property.Reference, property.Title)
		} else {
			if err := saveProperty(ctx, client, property); err != nil {
				log.Printf("Error importing property %s: %v", property.Reference, err)
				errors++
			} else {
				log.Printf("Imported property %d/%d: %s - %s", i+1, len(rows)-1, property.Reference, property.Title)
				imported++
			}
		}

		// Add a small delay to avoid rate limiting
		if !dryRun && i > 0 && i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return imported, errors
}

func convertXMLToProperty(imovel XMLImovel, tenantID, ownerID string) Property {
	now := time.Now()
	id := uuid.New().String()

	// Determine status
	status := "available"
	if imovel.Venda == 0 && imovel.Locacao == 0 {
		status = "inactive"
	}

	// Determine purpose
	purpose := "sale"
	if imovel.Locacao == 1 && imovel.Venda == 0 {
		purpose = "rent"
	} else if imovel.Locacao == 1 && imovel.Venda == 1 {
		purpose = "both"
	}

	// Normalize type
	propertyType := normalizeType(imovel.Tipo)

	// Build features
	features := make(map[string]interface{})
	if imovel.Churrasqueira == 1 {
		features["churrasqueira"] = true
	}
	if imovel.Playground == 1 {
		features["playground"] = true
	}
	if imovel.Piscina == 1 {
		features["piscina"] = true
	}
	if imovel.Arcondicionado == 1 {
		features["ar_condicionado"] = true
	}
	if imovel.Armariocozinha == 1 {
		features["armario_cozinha"] = true
	}
	if imovel.Lavanderia == 1 {
		features["lavanderia"] = true
	}
	if imovel.Sacada == 1 {
		features["sacada"] = true
	}
	if imovel.Varanda == 1 {
		features["varanda"] = true
	}

	// Build address
	address := map[string]interface{}{
		"street":       imovel.Endereco,
		"number":       imovel.Numero,
		"complement":   imovel.Complemento,
		"neighborhood": imovel.Bairro,
		"city":         imovel.Cidade,
		"state":        imovel.UnidadeFederativa,
		"postal_code":  imovel.CEP,
		"country":      "Brasil",
	}

	if imovel.Latitude != "" {
		address["latitude"] = imovel.Latitude
	}
	if imovel.Longitude != "" {
		address["longitude"] = imovel.Longitude
	}

	// Create slug from title or reference
	slug := createSlug(imovel.Titulo, imovel.Referencia)

	return Property{
		ID:               id,
		TenantID:         tenantID,
		Reference:        imovel.Referencia,
		Title:            cleanString(imovel.Titulo),
		Description:      cleanString(imovel.Anuncioparainternet),
		Type:             propertyType,
		Status:           status,
		Purpose:          purpose,
		Category:         strings.ToLower(imovel.Categoria),
		SalePrice:        imovel.Valorvenda,
		RentalPrice:      imovel.Valorlocacao,
		CondoFee:         imovel.Valorcondominio,
		Address:          address,
		Features:         features,
		OwnerID:          ownerID,
		Bedrooms:         imovel.Dormitorios,
		Bathrooms:        imovel.Banheiro2,
		Suites:           imovel.Suite,
		ParkingSpaces:    imovel.Garagem,
		Area:             imovel.Areacosntruida,
		TotalArea:        imovel.Areatotal,
		AcceptsFinancing: imovel.Aceitafinanciamento == 1,
		AcceptsTrade:     imovel.Permuta == 1,
		ExternalID:       imovel.Codigoimovel,
		ExternalURL:      imovel.LinkImovelSite,
		Slug:             slug,
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedBy:        "import",
		UpdatedBy:        "import",
	}
}

func convertXLSToProperty(headers, row []string, tenantID, ownerID string) Property {
	now := time.Now()
	id := uuid.New().String()

	// Helper function to get cell value
	getCell := func(colName string) string {
		for i, h := range headers {
			if strings.Contains(h, colName) {
				if i < len(row) {
					return row[i]
				}
			}
		}
		return ""
	}

	// Helper function to parse float
	parseFloat := func(s string) float64 {
		s = strings.TrimSpace(s)
		if s == "" {
			return 0
		}
		f, _ := strconv.ParseFloat(s, 64)
		return f
	}

	// Helper function to parse int
	parseInt := func(s string) int {
		s = strings.TrimSpace(s)
		if s == "" {
			return 0
		}
		i, _ := strconv.Atoi(s)
		return i
	}

	reference := getCell("Refer")
	propertyType := normalizeType(getCell("Tipo"))
	title := getCell("Link no Site")
	description := getCell("Descri")
	salePrice := parseFloat(getCell("Valor Venda"))
	rentalPrice := parseFloat(getCell("Valor Loca"))
	condoFee := parseFloat(getCell("Valor Condom"))

	// Determine purpose and status
	purpose := "sale"
	status := "available"
	finalidade := getCell("Finalidade")
	if strings.Contains(strings.ToUpper(finalidade), "LOCA") {
		purpose = "rent"
	}
	situacao := getCell("Situa")
	if !strings.Contains(strings.ToUpper(situacao), "ATIVO") {
		status = "inactive"
	}

	// Build address
	address := map[string]interface{}{
		"street":       getCell("Endere"),
		"number":       getCell("mero"),
		"complement":   getCell("Complemento"),
		"neighborhood": getCell("Bairro"),
		"city":         getCell("Cidade"),
		"state":        getCell("UF"),
		"postal_code":  getCell("CEP"),
		"country":      "Brasil",
	}

	// Build features from "Outras características"
	features := make(map[string]interface{})
	otherFeatures := getCell("Outras caracter")
	if otherFeatures != "" {
		featuresList := strings.Split(otherFeatures, ",")
		for _, f := range featuresList {
			f = strings.TrimSpace(f)
			if f != "" {
				features[strings.ToLower(f)] = true
			}
		}
	}

	slug := createSlug(title, reference)

	return Property{
		ID:               id,
		TenantID:         tenantID,
		Reference:        reference,
		Title:            cleanString(title),
		Description:      cleanString(description),
		Type:             propertyType,
		Status:           status,
		Purpose:          purpose,
		Category:         "normal",
		SalePrice:        salePrice * 1000, // Values in XLS are in thousands
		RentalPrice:      rentalPrice,
		CondoFee:         condoFee,
		Address:          address,
		Features:         features,
		OwnerID:          ownerID,
		Bedrooms:         parseInt(getCell("Dorms")),
		Bathrooms:        parseInt(getCell("Banh")),
		Suites:           parseInt(getCell("Su")),
		ParkingSpaces:    parseInt(getCell("Gar.")),
		Area:             parseFloat(getCell("rea Constru")),
		TotalArea:        parseFloat(getCell("rea Total")),
		AcceptsFinancing: false,
		AcceptsTrade:     false,
		ExternalID:       getCell("ID Internet"),
		ExternalURL:      getCell("Link no Site"),
		Slug:             slug,
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedBy:        "import",
		UpdatedBy:        "import",
	}
}

func saveProperty(ctx context.Context, client *firestore.Client, property Property) error {
	_, err := client.Collection("properties").Doc(property.ID).Set(ctx, property)
	return err
}

func normalizeType(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))

	// Map common types
	typeMap := map[string]string{
		"apartamento": "apartment",
		"casa":        "house",
		"terreno":     "land",
		"chácara":     "farm",
		"sítio":       "farm",
		"fazenda":     "farm",
		"comercial":   "commercial",
		"sala":        "commercial",
		"loja":        "commercial",
		"galpão":      "warehouse",
	}

	for key, value := range typeMap {
		if strings.Contains(t, key) {
			return value
		}
	}

	return "other"
}

func cleanString(s string) string {
	// Remove extra whitespace
	s = strings.TrimSpace(s)
	// Replace multiple spaces with single space
	s = strings.Join(strings.Fields(s), " ")
	return s
}

func createSlug(title, reference string) string {
	if title == "" {
		title = reference
	}

	// Simple slug creation
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slug = result.String()

	// Remove multiple hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 100 {
		slug = slug[:100]
	}

	return slug
}
