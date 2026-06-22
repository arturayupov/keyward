package envfile

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Parse reads KEY=VALUE lines, honoring `export ` prefixes, # comments,
// and single/double quoted values.
func Parse(b []byte) (map[string]string, error) {
	out := map[string]string{}
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		v = unquote(v)
		out[k] = v
	}
	return out, sc.Err()
}

// ParseFile parses the env file at path.
func ParseFile(path string) (map[string]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b)
}

// Set writes key=value into the env file at path, replacing an existing line
// for key or appending. Creates the file 0600. Never returns or logs value.
func Set(path, key, value string) error {
	var lines []string
	if b, err := os.ReadFile(path); err == nil {
		lines = strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	} else if !os.IsNotExist(err) {
		return err
	}
	replaced := false
	for i, ln := range lines {
		trimmed := strings.TrimPrefix(strings.TrimSpace(ln), "export ")
		if k, _, ok := strings.Cut(trimmed, "="); ok && strings.TrimSpace(k) == key {
			lines[i] = fmt.Sprintf("%s=%s", key, quote(value))
			replaced = true
			break
		}
	}
	if !replaced {
		if len(lines) == 1 && lines[0] == "" {
			lines = nil
		}
		lines = append(lines, fmt.Sprintf("%s=%s", key, quote(value)))
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o600)
}

// quote renders a value for a single .env line. Values with whitespace, quotes,
// '#', or newlines are double-quoted with \", \\, and \n escaped, so a value
// (including a multi-line PEM key) round-trips through Parse unchanged.
func quote(v string) string {
	if !strings.ContainsAny(v, " \t\"'#\n\r\\") {
		return v
	}
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`)
	return `"` + r.Replace(v) + `"`
}

// unquote reverses quote: it strips matching surrounding quotes and, for
// double-quoted values, unescapes \\, \", \n, and \r.
func unquote(v string) string {
	if len(v) < 2 || v[0] != v[len(v)-1] || (v[0] != '"' && v[0] != '\'') {
		return v
	}
	inner := v[1 : len(v)-1]
	if v[0] == '\'' {
		return inner // single quotes are literal
	}
	// Single pass, no rescanning: \\ is tried before the others, so an escaped
	// backslash is consumed whole and never re-interpreted.
	return strings.NewReplacer(`\\`, `\`, `\"`, `"`, `\n`, "\n", `\r`, "\r").Replace(inner)
}
