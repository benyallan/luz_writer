package model

// Problem é um item emitido pelo Rule Engine (seção 9 da spec).
type Problem struct {
	Severity string `json:"severity"` // "error" | "warning" | "info"
	Code     string `json:"code"`     // ex.: "R001"
	Message  string `json:"message"`  // em pt-BR
	Source   string `json:"source"`   // "project" | "target:<id>" | "chapter:<id>" | "plugin:<name>" | "styles" | "override:<docId>"
}

// BuildResult é o retorno de App.Compile() e o payload do evento luz:build:done.
type BuildResult struct {
	Success    bool      `json:"success"`
	OutputPath string    `json:"outputPath"` // caminho em dist/
	Problems   []Problem `json:"problems"`
	LogTail    string    `json:"logTail"` // últimas ~40 linhas do log do Tectonic em caso de falha
}

// BuildProgress é o payload do evento luz:build:progress.
type BuildProgress struct {
	Stage   string  `json:"stage"` // "validating" | "generating" | "compiling" | "done"
	Percent float64 `json:"percent"`
}
