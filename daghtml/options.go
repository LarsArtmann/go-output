package daghtml

// Option configures the DAG HTML renderer via functional options.
type Option func(*config)

type config struct {
	title       string
	subtitle    string
	containerID string
	dataID      string
	height      int
	footer      string
}

func defaultConfig() config {
	return config{
		title:       "DAG Visualization",
		subtitle:    "",
		containerID: "dag-container",
		dataID:      "dag-data",
		height:      500,
		footer:      "",
	}
}

// WithTitle sets the page <title> and visible header text.
func WithTitle(title string) Option {
	return func(c *config) { c.title = title }
}

// WithSubtitle sets the subtitle text shown below the title.
func WithSubtitle(subtitle string) Option {
	return func(c *config) { c.subtitle = subtitle }
}

// WithContainerID overrides the HTML element ID for the graph container div.
// Use this when embedding multiple graphs on a single page.
func WithContainerID(id string) Option {
	return func(c *config) { c.containerID = id }
}

// WithDataID overrides the HTML element ID for the <script> tag holding the
// DAG JSON data. Must be unique per page when embedding multiple graphs.
func WithDataID(id string) Option {
	return func(c *config) { c.dataID = id }
}

// WithHeight sets the graph container height in pixels (default 500).
func WithHeight(height int) Option {
	return func(c *config) { c.height = height }
}

// WithFooter sets the footer text shown at the bottom of a full page.
func WithFooter(footer string) Option {
	return func(c *config) { c.footer = footer }
}
