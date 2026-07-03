package model

// FieldType é o tipo de campo suportado pelo SchemaForm.vue no MVP (seção 8.2).
type FieldType string

const (
	FieldTypeSwitch    FieldType = "switch"
	FieldTypeSelect    FieldType = "select"
	FieldTypeDimension FieldType = "dimension" // número + unidade in|cm|mm|pt
	FieldTypeText      FieldType = "text"
	FieldTypeNumber    FieldType = "number"
)

// FieldOption é uma opção de um campo "select".
type FieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// FormField é um campo de um FormSchema (seção 8.2).
type FormField struct {
	Key     string        `json:"key"`
	Label   string        `json:"label"`
	Type    FieldType     `json:"type"`
	Default any           `json:"default,omitempty"`
	Options []FieldOption `json:"options,omitempty"`
}

// FormSchema descreve o formulário que SchemaForm.vue renderiza genericamente
// para configurar um plugin (seção 8.2).
type FormSchema struct {
	Fields []FormField `json:"fields"`
}
