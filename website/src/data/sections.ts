import type { PhaseCard, ComparisonItem, UseCase } from "./types";

export const phases: PhaseCard[] = [
  {
    step: "1",
    stepColor: "accent",
    title: "Build — Mutable",
    code: "b := output.NewGraphBuilder()\nb.AddNode(node)\nb.AddEdge(edge)",
    hint: "Fluent API, chainable calls",
    hintColor: "accent",
    hintIcon: "bolt",
  },
  {
    step: "2",
    stepColor: "amber",
    title: "Freeze — Immutable",
    code: "g := b.Build()\n// g.Nodes() / g.Edges() only",
    hint: "Snapshot taken, no further mutation",
    hintColor: "amber",
    hintIcon: "check",
  },
  {
    step: "3",
    stepColor: "violet",
    title: "Render — Pure Functions",
    code: "dot, _     := graph.RenderDOT(g)\nmermaid, _ := graph.RenderMermaid(g)",
    hint: "Same Graph, multiple formats, no side effects",
    hintColor: "violet",
    hintIcon: "chart",
  },
];

export const comparisons: ComparisonItem[] = [
  {
    variant: "DIY",
    price: "DIY",
    accent: false,
    pros: [],
    cons: [
      "Hand-write each format",
      "No type safety on enums",
      "Duplicate logic per renderer",
      "No cross-shape conversion",
    ],
  },
  {
    variant: "go-output",
    price: "Free",
    accent: true,
    pros: [
      "16 formats from one data model",
      "Type-safe enums + branded IDs",
      "CQRS pure-function renderers",
      "Cross-shape projections built-in",
      "NOM real-time progress included",
    ],
    cons: [],
  },
  {
    variant: "Heavy",
    price: "Heavy",
    accent: false,
    pros: [],
    cons: [
      "Pulls lipgloss/bubbletea everywhere",
      "Tightly coupled renderers",
      "No streaming-first design",
      "Large dependency footprint",
    ],
  },
];

export const useCases: UseCase[] = [
  {
    title: "CLI Tools",
    desc: "One data model, runtime format flags. CSV for scripts, tables for terminals, JSON for pipelines.",
    icon: "terminal",
  },
  {
    title: "CI/CD Pipelines",
    desc: "NOM progress tracking with CI-safe degradation. Auto-detects non-TTY and appends frames line-by-line.",
    icon: "refresh",
  },
  {
    title: "Documentation",
    desc: "Render architecture diagrams as D2, Mermaid, DOT, or PlantUML from a single graph builder.",
    icon: "document",
  },
  {
    title: "Dashboards",
    desc: "Interactive SVG DAGs for HTML. Pan, zoom, click-highlight. Zero external dependencies.",
    icon: "chart",
  },
];
