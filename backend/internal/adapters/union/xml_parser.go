package union

import (
	"encoding/xml"
	"io"
)

// XMLUnion represents the root XML structure from Union CRM
type XMLUnion struct {
	XMLName xml.Name     `xml:"Union"`
	Imoveis []XMLImovel  `xml:"Imoveis>Imovel"`
}

// XMLImovel represents a property from Union XML
type XMLImovel struct {
	// Identificação
	CodigoImobiliaria   string `xml:"CodigoImobiliaria"`
	Codigoimovel        string `xml:"Codigoimovel"` // external_id principal
	Referencia          string `xml:"Referencia"`   // ex: AP00335

	// Tipo e Categoria
	Tipo       string `xml:"Tipo"`       // APARTAMENTO, CASA, TERRENO, etc.
	Finalidade string `xml:"Finalidade"` // RESIDENCIAL, COMERCIAL, RURAL
	Categoria  string `xml:"Categoria"`  // NORMAL, DESTAQUE, etc.

	// Títulos e Descrições
	Titulo              string `xml:"Titulo"`
	Anuncioparainternet string `xml:"Anuncioparainternet"`
	MetadescriptionSEO  string `xml:"MetadescriptionSEO"`
	PalavraschavesSEO   string `xml:"PalavraschavesSEO"`

	// Valores
	Valorvenda      float64 `xml:"Valorvenda"`
	Valorlocacao    float64 `xml:"Valorlocacao"`
	Valortemporada  float64 `xml:"Valortemporada"`
	Valorcondominio float64 `xml:"Valorcondominio"`
	Valoriptu       float64 `xml:"Valoriptu"`

	// Características
	Dormitorios    int     `xml:"Dormitorios"`
	Suite          int     `xml:"Suite"`
	Banheiro       int     `xml:"Banheiro"`
	Banheiro2      int     `xml:"Banheiro2"` // Union usa este campo
	Garagem        int     `xml:"Garagem"`
	Areacosntruida float64 `xml:"Areacosntruida"` // typo no XML original
	Areatotal      float64 `xml:"Areatotal"`
	Areautil       float64 `xml:"Areautil"`

	// Endereço
	Endereco          string `xml:"Endereco"`
	Numero            string `xml:"Numero"`
	Complemento       string `xml:"Complemento"`
	Bairro            string `xml:"Bairro"`
	Cidade            string `xml:"Cidade"`
	UnidadeFederativa string `xml:"UnidadeFederativa"` // UF
	CEP               string `xml:"CEP"`
	Latitude          string `xml:"Latitude"`
	Longitude         string `xml:"Longitude"`

	// Flags de Transação
	Venda                int `xml:"Venda"`
	Locacao              int `xml:"Locacao"`
	Temporada            int `xml:"Temporada"`
	Aceitafinanciamento  int `xml:"Aceitafinanciamento"`
	Permuta              int `xml:"Permuta"`
	Valorsobconsulta     int `xml:"Valorsobconsulta"`

	// Features (checklist extenso)
	Churrasqueira    int `xml:"Churrasqueira"`
	Playground       int `xml:"Playground"`
	Piscina          int `xml:"Piscina"`
	Arcondicionado   int `xml:"Arcondicionado"`
	Armariocozinha   int `xml:"Armariocozinha"`
	Lavanderia       int `xml:"Lavanderia"`
	Sacada           int `xml:"Sacada"`
	Varanda          int `xml:"Varanda"`
	Elevador         int `xml:"Elevador"`
	Portaria24horas  int `xml:"Portaria24horas"`
	Salafesta        int `xml:"Salafesta"`
	Quadrapoliesportiva int `xml:"Quadrapoliesportiva"`
	Salacinema       int `xml:"Salacinema"`
	Salaginastica    int `xml:"Salaginastica"`
	Sauna            int `xml:"Sauna"`
	Piscinavaquecida int `xml:"Piscinaaquecida"`
	Jardim           int `xml:"Jardim"`
	Gourmet          int `xml:"Gourmet"`

	// Metadata
	Cadastradoem  string `xml:"Cadastradoem"`
	Atualizadoem  string `xml:"Atualizadoem"`
	Captador      string `xml:"Captador"`      // Corretor captador
	Cadastradopor string `xml:"Cadastradopor"` // Quem cadastrou
	Atualizadopor string `xml:"Atualizadopor"` // Quem atualizou

	// Links
	LinkImovelSite string `xml:"LinkImovelSite"`

	// Fotos (podem ser múltiplas tags <Foto>)
	Fotos []string `xml:"Fotos>Foto"`

	// Condomínio/Empreendimento
	Condominio     int    `xml:"Condominio"`     // 0 ou 1
	Condominionome string `xml:"Condominionome"` // Nome do condomínio
	Empreendimento string `xml:"Empreendimento"`
	Construtora    string `xml:"Construtora"`
	Edificio       string `xml:"Edificio"`

	// Ano de Construção
	AnoConstrucao string `xml:"AnoConstrucao"`
}

// ParseXML parses Union XML file
func ParseXML(reader io.Reader) (*XMLUnion, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var union XMLUnion
	if err := xml.Unmarshal(data, &union); err != nil {
		return nil, err
	}

	return &union, nil
}
