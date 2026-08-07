package serialization

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

// stringFromBytes wraps the standard "marshal bytes to string with an error
// context" idiom that every tree/graph renderer uses. It removes the
// four-line copy/paste that toml/yaml/json renderers all carried.
//
// Callers pass the raw multi-value result of a Marshal function directly:
//
//	return stringFromBytes("toml", "tree", toml.Marshal(node))
//
// The helper takes both []byte and error positions of the marshal call's
// return signature, matching the function-style marshal APIs in this module.
func stringFromBytes(format, subject string, data []byte, err error) (string, error) {
	if err != nil {
		return "", fmt.Errorf("marshal %s %s: %w", format, subject, err)
	}

	return string(data), nil
}

// renderTreeNode is the shared Render body for {JSON,TOML,YAML}TreeRenderer.
// It collapses the "if r.root == nil { return <format-empty> }" guard together
// with the toTreeNode projection and the stringFromBytes wrapper, leaving each
// renderer's Render() body as a single line.
//
// The empty-payload value differs per format (JSON: "null", TOML: "", YAML:
// "null\n") so the caller passes the right one. The marshal function is the
// per-format encoder (json.Marshal/yaml.Marshal/toml.Marshal or a closure that
// captures json.Deterministic + indent options).
func renderTreeNode(
	root *output.TreeNode,
	emptyPayload, format string,
	marshal func(any) ([]byte, error),
) (string, error) {
	if root == nil {
		return emptyPayload, nil
	}

	data, err := marshal(toTreeNode(root))
	if err != nil {
		return "", fmt.Errorf("marshal %s %s: %w", format, "tree", err)
	}

	return string(data), nil
}

// writeEmptyArrayPayload writes a minimal empty-collection payload ("[]\n")
// to w, shared between the JSON and YAML streaming encoders. TOML omits this
// because its encoder produces only document metadata when given an empty
// wrapper. Returns nil on a successful write.
func writeEmptyArrayPayload(w io.Writer, format string) error {
	if _, err := io.WriteString(w, "[]\n"); err != nil {
		return fmt.Errorf("write empty %s: %w", format, err)
	}

	return nil
}
