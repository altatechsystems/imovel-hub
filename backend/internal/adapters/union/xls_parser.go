package union

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"golang.org/x/net/html"
)

// XLSRecord represents a row from Union XLS export
type XLSRecord struct {
	// Identificação
	Referencia  string
	IDInternet  string
	CodigoImovel string

	// Tipo
	Tipo       string
	Finalidade string
	Categoria  string

	// Valores
	ValorVenda     float64
	ValorLocacao   float64
	ValorTemporada float64
	ValorCondominio float64
	ValorIPTU      float64

	// Características
	Dormitorios    int
	Suites         int
	Garagens       int
	Banheiros      int
	AreaConstruida float64
	AreaUtil       float64
	AreaTotal      float64

	// Endereço
	Endereco   string
	Numero     string
	Complemento string
	Bairro     string
	Cidade     string
	UF         string
	CEP        string

	// Proprietário (OWNER DATA)
	Proprietario     string // Nome do proprietário
	Empresa          string
	CelularTelefone  string
	Email            string

	// Metadata
	DataCadastro     string
	DataAtualizacao  string
	Captador         string
	Equipe           string
	LocalChaves      string

	// Links
	LinkSite       string
	FotoPrincipal  string

	// Descrição
	Descricao string

	// Features (como texto)
	DetalhesBasico       string
	DetalhesServicos     string
	DetalhesLazer        string
	DetalhesSocial       string
	OutrasCaracteristicas string
}

// ParseXLS parses Union XLS file (supports both XLSX and HTML-based XLS)
func ParseXLS(filePath string) ([]XLSRecord, error) {
	// Try to detect if this is an HTML file
	isHTML, err := isHTMLFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to detect file format: %w", err)
	}

	if isHTML {
		fmt.Printf("📄 Detected HTML-based XLS file, using HTML parser\n")
		return parseHTMLXLS(filePath)
	}

	// Try to parse as XLSX
	fmt.Printf("📊 Attempting to parse as XLSX file\n")
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XLS: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in XLS")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("XLS has no data rows")
	}

	// Parse header to create column map
	headers := rows[0]
	colMap := makeColumnMap(headers)

	var records []XLSRecord
	for i, row := range rows[1:] {
		record, err := parseXLSRow(row, colMap)
		if err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to parse row %d: %v\n", i+2, err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

// isHTMLFile checks if the file is HTML by reading the first few bytes
func isHTMLFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read first 512 bytes to detect format
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

// parseHTMLXLS parses HTML-based XLS file (common export format)
func parseHTMLXLS(filePath string) ([]XLSRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open HTML file: %w", err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract table rows
	rows := extractTableRows(doc)
	if len(rows) < 2 {
		return nil, fmt.Errorf("HTML XLS has no data rows (found %d rows)", len(rows))
	}

	fmt.Printf("📋 Extracted %d rows from HTML table\n", len(rows))

	// Parse header to create column map
	headers := rows[0]
	colMap := makeColumnMap(headers)

	fmt.Printf("📌 Found %d columns: %v\n", len(headers), headers)

	var records []XLSRecord
	for i, row := range rows[1:] {
		record, err := parseXLSRow(row, colMap)
		if err != nil {
			// Log error but continue
			fmt.Printf("⚠️  Warning: failed to parse row %d: %v\n", i+2, err)
			continue
		}
		records = append(records, record)
	}

	fmt.Printf("✅ Successfully parsed %d records from HTML XLS\n", len(records))
	return records, nil
}

// extractTableRows extracts all rows from an HTML table
func extractTableRows(n *html.Node) [][]string {
	var rows [][]string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			row := extractTableCells(n)
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return rows
}

// extractTableCells extracts all cells from a table row
func extractTableCells(tr *html.Node) []string {
	var cells []string

	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
			cells = append(cells, getNodeText(c))
		}
	}

	return cells
}

// getNodeText extracts text content from an HTML node
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

// makeColumnMap creates a map from column name to index
// Handles variations in column names (with/without accents, different cases)
func makeColumnMap(headers []string) map[string]int {
	colMap := make(map[string]int)

	for i, header := range headers {
		// Normalize: lowercase, remove accents, trim
		normalized := normalizeColumnName(header)
		colMap[normalized] = i
	}

	return colMap
}

// normalizeColumnName normalizes column name for flexible matching
func normalizeColumnName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	// Remove common accents (basic approach)
	replacements := map[string]string{
		"á": "a", "à": "a", "ã": "a", "â": "a",
		"é": "e", "ê": "e",
		"í": "i",
		"ó": "o", "õ": "o", "ô": "o",
		"ú": "u", "ü": "u",
		"ç": "c",
	}

	for old, new := range replacements {
		name = strings.ReplaceAll(name, old, new)
	}

	return name
}

// getCell gets cell value by column name (flexible matching)
func getCell(row []string, colMap map[string]int, variations ...string) string {
	for _, variation := range variations {
		normalized := normalizeColumnName(variation)
		if idx, ok := colMap[normalized]; ok {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
		}
	}
	return ""
}

// parseFloat parses float with fallback to 0
func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".") // Handle comma decimal separator
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// parseInt parses int with fallback to 0
func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

// parseXLSRow parses a single XLS row into XLSRecord
func parseXLSRow(row []string, colMap map[string]int) (XLSRecord, error) {
	record := XLSRecord{
		// Identificação
		Referencia:   getCell(row, colMap, "Referência", "Referencia", "Ref"),
		IDInternet:   getCell(row, colMap, "ID Internet", "IDInternet", "Codigo"),
		CodigoImovel: getCell(row, colMap, "Código Imóvel", "Codigo Imovel", "CodImovel"),

		// Tipo
		Tipo:       getCell(row, colMap, "Tipo"),
		Finalidade: getCell(row, colMap, "Finalidade"),
		Categoria:  getCell(row, colMap, "Categoria"),

		// Valores
		ValorVenda:      parseFloat(getCell(row, colMap, "Valor Venda", "ValorVenda")),
		ValorLocacao:    parseFloat(getCell(row, colMap, "Valor Locação", "Valor Locacao", "ValorLocacao")),
		ValorTemporada:  parseFloat(getCell(row, colMap, "Valor Temporada", "ValorTemporada")),
		ValorCondominio: parseFloat(getCell(row, colMap, "Valor Condomínio", "Valor Condominio", "ValorCondominio")),
		ValorIPTU:       parseFloat(getCell(row, colMap, "Valor IPTU", "ValorIPTU")),

		// Características
		Dormitorios:    parseInt(getCell(row, colMap, "Dorms.", "Dormitorios", "Quartos")),
		Suites:         parseInt(getCell(row, colMap, "Suíte", "Suite", "Suites")),
		Garagens:       parseInt(getCell(row, colMap, "Gar.", "Garagem", "Garagens", "Vagas")),
		Banheiros:      parseInt(getCell(row, colMap, "Banh.", "Banheiro", "Banheiros")),
		AreaConstruida: parseFloat(getCell(row, colMap, "Área Construída", "Area Construida", "AreaConstruida")),
		AreaUtil:       parseFloat(getCell(row, colMap, "Área Útil", "Area Util", "AreaUtil")),
		AreaTotal:      parseFloat(getCell(row, colMap, "Área Total", "Area Total", "AreaTotal")),

		// Endereço
		Endereco:    getCell(row, colMap, "Endereço", "Endereco"),
		Numero:      getCell(row, colMap, "Número", "Numero", "Nº", "N."),
		Complemento: getCell(row, colMap, "Complemento"),
		Bairro:      getCell(row, colMap, "Bairro"),
		Cidade:      getCell(row, colMap, "Cidade"),
		UF:          getCell(row, colMap, "UF", "Estado"),
		CEP:         getCell(row, colMap, "CEP"),

		// Proprietário (CRITICAL FOR OWNER ENRICHMENT)
		Proprietario:    getCell(row, colMap, "Proprietário", "Proprietario", "Owner"),
		Empresa:         getCell(row, colMap, "Empresa"),
		CelularTelefone: getCell(row, colMap, "Celular/Telefone", "Telefone", "Celular"),
		Email:           getCell(row, colMap, "E-mail", "Email"),

		// Metadata
		DataCadastro:    getCell(row, colMap, "Data Cadastro", "DataCadastro"),
		DataAtualizacao: getCell(row, colMap, "Data Atualização", "Data Atualizacao", "DataAtualizacao"),
		Captador:        getCell(row, colMap, "Captador"),
		Equipe:          getCell(row, colMap, "Equipe"),
		LocalChaves:     getCell(row, colMap, "Local das Chaves", "LocalChaves"),

		// Links
		LinkSite:      getCell(row, colMap, "Link no Site", "LinkSite"),
		FotoPrincipal: getCell(row, colMap, "Foto principal", "FotoPrincipal"),

		// Descrição
		Descricao: getCell(row, colMap, "Descrição", "Descricao"),

		// Features
		DetalhesBasico:        getCell(row, colMap, "Detalhes Básico", "Detalhes Basico"),
		DetalhesServicos:      getCell(row, colMap, "Detalhes Serviços", "Detalhes Servicos"),
		DetalhesLazer:         getCell(row, colMap, "Detalhes Lazer"),
		DetalhesSocial:        getCell(row, colMap, "Detalhes Social"),
		OutrasCaracteristicas: getCell(row, colMap, "Outras características", "Outras caracteristicas"),
	}

	return record, nil
}

// FindXLSRecordByCode finds XLS record by matching external code
// Priority: 1) CodigoImovel, 2) IDInternet, 3) Referencia
func FindXLSRecordByCode(records []XLSRecord, xmlImovel *XMLImovel) *XLSRecord {
	// Try matching by Codigoimovel (strongest match)
	if xmlImovel.Codigoimovel != "" {
		for i := range records {
			if records[i].CodigoImovel == xmlImovel.Codigoimovel ||
				records[i].IDInternet == xmlImovel.Codigoimovel {
				return &records[i]
			}
		}
	}

	// Try matching by Referencia (e.g., AP00335)
	if xmlImovel.Referencia != "" {
		for i := range records {
			if records[i].Referencia == xmlImovel.Referencia {
				return &records[i]
			}
		}
	}

	// No match found
	return nil
}
