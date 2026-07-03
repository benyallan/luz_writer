package latex

import "testing"

func TestEscapeText(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`\`, `\textbackslash{}`},
		{`{`, `\{`},
		{`}`, `\}`},
		{`$`, `\$`},
		{`&`, `\&`},
		{`#`, `\#`},
		{`^`, `\^{}`},
		{`_`, `\_`},
		{`%`, `\%`},
		{`~`, `\textasciitilde{}`},
		{`texto normal`, `texto normal`},
		{`50% de R$100 & outros_valores {x}`, `50\% de R\$100 \& outros\_valores \{x\}`},
		{`a\b{c}`, `a\textbackslash{}b\{c\}`},
	}
	for _, c := range cases {
		if got := EscapeText(c.in); got != c.want {
			t.Errorf("EscapeText(%q) = %q, esperava %q", c.in, got, c.want)
		}
	}
}
