package model

// CustomStyle é um estilo de texto nomeado, composto a partir de blocos
// seguros (seção 5.7) — equivalente amigável ao \newcommand de formatação.
type CustomStyle struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Italic    bool    `json:"italic"`
	Bold      bool    `json:"bold"`
	SmallCaps bool    `json:"smallCaps"`
	Color     *string `json:"color"` // hex ou nil
}

// StylesFile é o conteúdo de .luz/styles.json (seção 5.7).
type StylesFile struct {
	LuzVersion int           `json:"luzVersion"`
	Styles     []CustomStyle `json:"styles"`
}
