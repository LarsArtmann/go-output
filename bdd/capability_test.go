package bdd_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/larsartmann/go-output"
)

// BDD specs describing the format-parsing and shape-capability experience.
var _ = ginkgo.Describe("format discovery and capabilities", func() {
	ginkgo.Describe("parsing format strings from CLI flags", func() {
		ginkgo.DescribeTable(
			"accepts known formats",
			func(input string, expected output.Format) {
				got, err := output.ParseFormat(input)

				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(got).To(gomega.Equal(expected))
			},
			ginkgo.Entry("json", "json", output.FormatJSON),
			ginkgo.Entry("csv", "csv", output.FormatCSV),
			ginkgo.Entry("markdown", "markdown", output.FormatMarkdown),
			ginkgo.Entry("yaml", "yaml", output.FormatYAML),
			ginkgo.Entry("tree", "tree", output.FormatTree),
			ginkgo.Entry("d2", "d2", output.FormatD2),
		)

		ginkgo.Context("when given an unknown format", func() {
			ginkgo.It("returns an error listing allowed values", func() {
				_, err := output.ParseFormat("xml-but-typoed")

				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("allowed:"))
			})
		})
	})

	ginkgo.Describe("querying the shape capability matrix", func() {
		ginkgo.Context("for table-capable formats", func() {
			ginkgo.It("reports JSON and CSV support table shape", func() {
				gomega.Expect(output.FormatJSON.Supports(output.ShapeTable)).To(gomega.BeTrue())
				gomega.Expect(output.FormatCSV.Supports(output.ShapeTable)).To(gomega.BeTrue())
			})

			ginkgo.It("reports CSV does not support graph shape", func() {
				gomega.Expect(output.FormatCSV.Supports(output.ShapeGraph)).To(gomega.BeFalse())
			})
		})

		ginkgo.Context("reverse lookup by shape", func() {
			ginkgo.It("includes table formats when querying for table shape", func() {
				tableFormats := output.FormatsForShape(output.ShapeTable)

				gomega.Expect(tableFormats).To(gomega.ContainElement(output.FormatJSON))
				gomega.Expect(tableFormats).To(gomega.ContainElement(output.FormatCSV))
			})
		})
	})

	ginkgo.Describe("format validation", func() {
		ginkgo.It("confirms a valid format", func() {
			gomega.Expect(output.FormatJSON.IsValid()).To(gomega.BeTrue())
		})

		ginkgo.It("rejects a bogus format", func() {
			gomega.Expect(output.Format("nonsense").IsValid()).To(gomega.BeFalse())
		})

		ginkgo.It("exposes all 16 formats via AllowedValues", func() {
			gomega.Expect(output.AllFormats).To(gomega.HaveLen(16))
		})
	})
})
