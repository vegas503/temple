package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil" // nolint: staticcheck
	"os"
	"strings"
	"text/template"
)

func exit(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(2)
}

func main() {
	var (
		inputFile  string
		outputFile string
	)

	flag.StringVar(&inputFile, "i", "-", "Input file ('-' for stdin)")
	flag.StringVar(&outputFile, "o", "-", "Output file ('-' for stdout)")
	flag.Parse()

	if inputFile == "" || outputFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var err error

	var input io.ReadCloser
	if inputFile == "-" {
		input = os.Stdin
	} else {
		input, err = os.Open(inputFile)
		if err != nil {
			exit(err)
		}
	}

	content, err := ioutil.ReadAll(input)
	if err != nil {
		exit(err)
	}

	input.Close()

	var output io.WriteCloser
	if outputFile == "-" {
		output = os.Stdout
	} else {
		output, err = os.Create(outputFile)
		if err != nil {
			exit(err)
		}
	}

	defer output.Close()

	mustRender(string(content), output)
}

func mustRender(input string, w io.Writer) {
	tpl, err := template.New("source").Funcs(template.FuncMap{
		"env":        envFunc,
		"envdefault": envdefaultFunc,
		"split":      splitFunc,
		"contains":   strings.Contains,
		"join":       strings.Join,
		"coalesce":   coalesceFunc,
		"append":     appendFunc,
		"uniq":       uniqFunc,
		"replace":    replaceFunc,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"istrue":     istrueFunc,
		"error":      errorFunc,
	}).Parse(input)
	if err != nil {
		exit(fmt.Errorf("parse template: %w", err))
	}

	if err := tpl.Execute(w, nil); err != nil {
		exit(err)
	}
}

func envFunc(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		exit(fmt.Errorf("required env variable %s not provided", key))
	}

	return v
}

func envdefaultFunc(key string, def string) string {
	if s := os.Getenv(key); s != "" {
		return s
	}

	return def
}

func splitFunc(s string, sep string) []string {
	if s == "" {
		return nil
	}

	res := []string{}

	for _, part := range strings.Split(s, sep) {
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			res = append(res, part)
		}
	}

	return res
}

func coalesceFunc(arg ...string) string {
	for _, s := range arg {
		if s != "" {
			return s
		}
	}

	return ""
}

func appendFunc(arg ...[]string) []string {
	res := []string{}
	for _, s := range arg {
		res = append(res, s...)
	}

	return res
}

func replaceFunc(s string, a string, b string) string {
	// Why not strings.Replace? Because we want compatibility with ancient Go versions.
	return strings.Replace(s, a, b, -1) // nolint: gocritic
}

func uniqFunc(a []string) []string {
	res := []string{}
	set := map[string]struct{}{}

	for _, s := range a {
		if _, ok := set[s]; !ok {
			res = append(res, s)
			set[s] = struct{}{}
		}
	}

	return res
}

func istrueFunc(val string) bool {
	s := strings.TrimSpace(strings.ToLower(val))
	if s == "" {
		return false
	}

	return s[0] == '1' || s[0] == 't' || s[0] == 'y' || s == "on"
}

func errorFunc(s string) error {
	exit(fmt.Errorf("%s", s))

	return nil
}
