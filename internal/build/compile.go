package build

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"luz-writer/internal/model"
	"luz-writer/internal/rules"
	"luz-writer/internal/workspace"
)

const compileTimeout = 120 * time.Second

// Compile executa o pipeline da seção 10: (1) Rule Engine — qualquer error
// aborta antes de gerar qualquer LaTeX; (2)-(3) preâmbulo e conteúdo, usando
// o target ativo real (seção 8.4) e os plugins habilitados em plugins.json;
// (4)-(6) escrita e compilação via Tectonic. onProgress pode ser nil.
func Compile(ws *workspace.Workspace, onProgress func(model.BuildProgress)) (model.BuildResult, error) {
	report := func(stage string, percent float64) {
		if onProgress != nil {
			onProgress(model.BuildProgress{Stage: stage, Percent: percent})
		}
	}

	report("validating", 5)

	ruleCtx, err := rules.NewContext(ws)
	if err != nil {
		return model.BuildResult{}, err
	}
	problems := rules.Validate(ruleCtx)
	if hasError(problems) {
		return model.BuildResult{Success: false, Problems: problems}, nil
	}

	report("generating", 10)

	in := BuildInputs{
		Workspace:       ws,
		Project:         ruleCtx.Project,
		Target:          ruleCtx.Target,
		EnabledOptional: ruleCtx.EnabledOptional,
		Styles:          ruleCtx.Styles,
		Chapters:        ruleCtx.Chapters,
	}

	// Passo 3 da seção 10: converte e agrupa todos os documentos primeiro,
	// para então (passo 2) montar o preâmbulo já sabendo quais idiomas
	// adicionais o babel precisa carregar.
	body, conversionProblems, langsUsed, err := AssembleBody(in)
	if err != nil {
		return model.BuildResult{}, err
	}
	problems = append(problems, conversionProblems...)
	if hasError(problems) {
		return model.BuildResult{Success: false, Problems: problems}, nil
	}

	preamble, err := BuildPreamble(in, langsUsed)
	if err != nil {
		return model.BuildResult{}, err
	}
	tex := preamble + "\n\\begin{document}\n\n" + titlePageBody(ruleCtx.Target) + body + "\n\\end{document}\n"

	report("generating", 30)

	tmpDir := filepath.Join(ws.Root, ".tmp")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return model.BuildResult{}, err
	}
	texPath := filepath.Join(tmpDir, "main.tex")
	if err := os.WriteFile(texPath, []byte(tex), 0o644); err != nil {
		return model.BuildResult{}, err
	}

	// O Tectonic resolve caminhos relativos (\includegraphics, \includepdf)
	// em relação à pasta do .tex de entrada (não ao cwd do processo) — por
	// isso espelhamos imagens/ e anexos/ dentro de .tmp/ a cada build,
	// mantendo "imagens/<arquivo>"/"anexos/<arquivo>" como os caminhos usados
	// tanto nos JSONs do workspace quanto no LaTeX gerado.
	if err := mirrorDir(filepath.Join(ws.Root, "imagens"), filepath.Join(tmpDir, "imagens")); err != nil {
		return model.BuildResult{}, err
	}
	if err := mirrorDir(filepath.Join(ws.Root, "anexos"), filepath.Join(tmpDir, "anexos")); err != nil {
		return model.BuildResult{}, err
	}

	distDir := filepath.Join(ws.Root, "dist")
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return model.BuildResult{}, err
	}

	report("compiling", 50)

	// Checado aqui (não só na inicialização do app) porque CheckTectonic()
	// extrai a cópia embutida sob demanda (seção 10) — resolve o caminho de
	// verdade a usar em vez de depender do "tectonic" bare estar no PATH, o
	// que falharia sem nenhuma pista em processos com PATH restrito.
	found, tectonicPath := CheckTectonic()
	if !found {
		return model.BuildResult{
			Success: false,
			Problems: append(problems, model.Problem{
				Severity: "error",
				Code:     "TECTONIC",
				Message:  "O binário 'tectonic' não foi encontrado — a exportação está desabilitada até ele ficar acessível para o Luz Writer.",
				Source:   "project",
			}),
			LogTail: "tectonic não encontrado (nem embutido, nem no PATH).",
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), compileTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, tectonicPath, "-o", "dist", filepath.Join(".tmp", "main.tex"))
	cmd.Dir = ws.Root
	output, runErr := cmd.CombinedOutput()

	if runErr != nil {
		logTail := tailLines(string(output), 40)
		if logTail == "" {
			logTail = runErr.Error()
		}
		return model.BuildResult{
			Success:  false,
			Problems: problems,
			LogTail:  logTail,
		}, nil
	}

	report("compiling", 90)

	slug := workspace.Slugify(ruleCtx.Project.Title)
	finalName := slug + ".pdf"
	if ruleCtx.Target.ID != "" {
		finalName = slug + "-" + ruleCtx.Target.ID + ".pdf"
	}
	finalPath := filepath.Join(distDir, finalName)
	if err := os.Rename(filepath.Join(distDir, "main.pdf"), finalPath); err != nil {
		return model.BuildResult{}, err
	}

	report("done", 100)

	return model.BuildResult{
		Success:    true,
		OutputPath: finalPath,
		Problems:   problems,
	}, nil
}

// mirrorDir substitui dst por uma cópia completa de src. Não é erro src não
// existir (projeto sem nenhuma imagem/anexo importado ainda).
func mirrorDir(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	return os.CopyFS(dst, os.DirFS(src))
}

func hasError(problems []model.Problem) bool {
	for _, p := range problems {
		if p.Severity == "error" {
			return true
		}
	}
	return false
}

func tailLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) <= n {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}
