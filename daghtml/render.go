package daghtml

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

// pageData is the data model for the HTML templates.
type pageData struct {
	Title       string
	Subtitle    string
	CSS         template.CSS
	JS          template.JS
	JSON        template.JS
	ContainerID string
	DataID      string
	Height      int
	Footer      string
}

const pageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; base-uri 'none'; frame-ancestors 'none';">
<title>{{.Title}}</title>
<style>
{{.CSS}}
</style>
</head>
<body>
<header>
<h1><span class="logo-dot"></span>{{.Title}}</h1>
{{if .Subtitle}}<p class="subtitle">{{.Subtitle}}</p>{{end}}
</header>
<div id="{{.ContainerID}}" class="dag-container" style="min-height:{{.Height}}px">
  <div class="graph-controls">
    <button class="graph-zoom-in" title="Zoom in" aria-label="Zoom in">+</button>
    <button class="graph-zoom-out" title="Zoom out" aria-label="Zoom out">&#8722;</button>
    <button class="graph-fit" title="Fit to view" aria-label="Fit to view">&#8962;</button>
  </div>
  <div class="graph-info">Scroll/pinch to zoom &middot; Drag to pan &middot; Click node to highlight</div>
</div>
<script type="application/json" id="{{.DataID}}">{{.JSON}}</script>
<script>{{.JS}}</script>
<script>initDAGGraph("{{.ContainerID}}", "{{.DataID}}");</script>
{{if .Footer}}<div class="footer"><span>{{.Footer}}</span></div>{{end}}
</body>
</html>`

const graphSectionTemplate = `<div id="{{.ContainerID}}" class="dag-container" style="min-height:{{.Height}}px">
  <div class="graph-controls">
    <button class="graph-zoom-in" title="Zoom in" aria-label="Zoom in">+</button>
    <button class="graph-zoom-out" title="Zoom out" aria-label="Zoom out">&#8722;</button>
    <button class="graph-fit" title="Fit to view" aria-label="Fit to view">&#8962;</button>
  </div>
  <div class="graph-info">Scroll/pinch to zoom &middot; Drag to pan &middot; Click node to highlight</div>
</div>
<script type="application/json" id="{{.DataID}}">{{.JSON}}</script>
<script>{{.JS}}</script>
<script>initDAGGraph("{{.ContainerID}}", "{{.DataID}}");</script>`

var (
	//nolint:gochecknoglobals // Package-level compiled template (parsed once at init).
	compiledPage = template.Must(template.New("page").Parse(pageTemplate))
	//nolint:gochecknoglobals // Package-level compiled template (parsed once at init).
	compiledGraphSection = template.Must(template.New("graph").Parse(graphSectionTemplate))
)

// prepareRender applies options to the default config and serializes the DAG
// to JSON. Used by [Render] and [GraphHTML] to share their preamble.
func prepareRender(opts []Option, dag DAG) (config, string, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	jsonData, err := dagToJSON(dag)
	if err != nil {
		return config{}, "", fmt.Errorf("serialize DAG: %w", err)
	}

	return cfg, jsonData, nil
}

// Render produces a complete, self-contained HTML page with an interactive
// Sugiyama-layered DAG visualization. The output includes all CSS and
// JavaScript inline — no external dependencies, no network requests.
func Render(dag DAG, opts ...Option) (string, error) {
	cfg, jsonData, err := prepareRender(opts, dag)
	if err != nil {
		return "", err
	}

	data := pageData{
		Title:       cfg.title,
		Subtitle:    cfg.subtitle,
		CSS:         template.CSS(graphCSS), //nolint:gosec // G203: trusted embedded CSS, not user input
		JS:          template.JS(graphJS),   //nolint:gosec // G203: trusted embedded JS, not user input
		JSON:        template.JS(jsonData),  //nolint:gosec // G203: HTML-escaped JSON via SetEscapeHTML(true)
		ContainerID: cfg.containerID,
		DataID:      cfg.dataID,
		Height:      cfg.height,
		Footer:      cfg.footer,
	}

	var buf bytes.Buffer
	if err := compiledPage.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render HTML: %w", err)
	}

	return buf.String(), nil
}

// Write writes a complete, self-contained HTML page to w.
func Write(w io.Writer, dag DAG, opts ...Option) error {
	html, err := Render(dag, opts...)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, html)
	if err != nil {
		return fmt.Errorf("write HTML: %w", err)
	}

	return nil
}

// GraphHTML returns just the graph container, controls, JSON data, and
// JavaScript — without the surrounding <html>, <head>, and <body> tags.
//
// Use this to embed the DAG visualization inside an existing HTML page or
// dashboard tab. The host page must define the CSS custom properties
// (--bg, --surface, --accent, --success, --error, etc.) or include the
// graph CSS via [StyleSheet].
func GraphHTML(dag DAG, opts ...Option) (string, error) {
	cfg, jsonData, err := prepareRender(opts, dag)
	if err != nil {
		return "", err
	}

	data := pageData{
		JS:          template.JS(graphJS),  //nolint:gosec // G203: trusted embedded JS
		JSON:        template.JS(jsonData), //nolint:gosec // G203: HTML-escaped JSON
		ContainerID: cfg.containerID,
		DataID:      cfg.dataID,
		Height:      cfg.height,
	}

	var buf bytes.Buffer
	if err := compiledGraphSection.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render graph section: %w", err)
	}

	return buf.String(), nil
}

// StyleSheet returns the CSS used by the DAG visualization as a string.
// Use this when embedding the graph in a host page via [GraphHTML].
func StyleSheet() string {
	return graphCSS
}

// Script returns the JavaScript used by the DAG visualization as a string.
// Use this when embedding the graph in a host page via [GraphHTML].
func Script() string {
	return graphJS
}
