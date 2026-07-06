package bdd_test

import (
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/larsartmann/go-output"
)

// BDD specs describing the end-user experience of a CLI developer using
// go-output to render tabular data across multiple formats.
var _ = ginkgo.Describe("a CLI developer formatting tabular data", func() {
	var data *output.Table

	ginkgo.BeforeEach(func() {
		data = output.NewTable([]string{"Name", "Role"})
		data.AddRow([]string{"Ada", "Engineer"})
		data.AddRow([]string{"Grace", "Architect"})
	})

	ginkgo.Describe("rendering via the unified dispatcher", func() {
		ginkgo.Context("when rendering to CSV", func() {
			ginkgo.It("writes headers followed by rows", func() {
				out := mustRender(data, output.FormatCSV)

				gomega.Expect(out).To(gomega.HavePrefix("Name,Role\n"))
				gomega.Expect(out).To(gomega.ContainSubstring("Ada,Engineer"))
				gomega.Expect(out).To(gomega.ContainSubstring("Grace,Architect"))
			})
		})

		ginkgo.Context("when rendering to TSV", func() {
			ginkgo.It("separates columns with tabs", func() {
				out := mustRender(data, output.FormatTSV)

				gomega.Expect(out).To(gomega.ContainSubstring("Name\tRole"))
				gomega.Expect(out).To(gomega.ContainSubstring("Ada\tEngineer"))
			})
		})

		ginkgo.Context("when rendering to Markdown", func() {
			ginkgo.It("produces a pipe-delimited table with a separator row", func() {
				out := mustRender(data, output.FormatMarkdown)

				gomega.Expect(out).To(gomega.ContainSubstring("Name"))
				gomega.Expect(out).To(gomega.ContainSubstring("Ada"))
				gomega.Expect(out).To(gomega.ContainSubstring("Engineer"))
				gomega.Expect(out).To(gomega.ContainSubstring("------"))
			})
		})
	})

	ginkgo.Describe("rendering a footer row", func() {
		ginkgo.BeforeEach(func() {
			data.SetFooter([]string{"Total", "2"})
		})

		ginkgo.Context("in CSV", func() {
			ginkgo.It("appends the footer as the final data row", func() {
				out := mustRender(data, output.FormatCSV)

				lines := strings.Split(strings.TrimSpace(out), "\n")
				gomega.Expect(lines).To(gomega.HaveLen(4))
				gomega.Expect(lines[len(lines)-1]).To(gomega.Equal("Total,2"))
			})
		})

		ginkgo.Context("in Markdown", func() {
			ginkgo.It("renders a bold summary row after the data", func() {
				out := mustRender(data, output.FormatMarkdown)

				gomega.Expect(out).To(gomega.ContainSubstring("Total"))
				gomega.Expect(strings.Count(out, "------")).To(gomega.BeNumerically(">=", 1))
			})
		})
	})

	ginkgo.Describe("validating data consistency", func() {
		ginkgo.Context("when the footer column count mismatches the headers", func() {
			ginkgo.It("returns a descriptive validation error", func() {
				data.SetFooter([]string{"only-one-column"})

				err := data.Validate()

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("column count"))
			})
		})
	})
})

// mustRender renders data in the given format to a string, failing the spec on error.
func mustRender(data *output.Table, format output.Format) string {
	var buf strings.Builder

	err := output.RenderTable(data, format, output.RenderOptions{Writer: &buf})
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	return buf.String()
}
