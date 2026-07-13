import type { Feature } from "./types";

export const features: Feature[] = [
  {
    icon: "layers",
    title: "16 Formats, One API",
    desc: "Same Table, TreeNode, or GraphNode. JSON, CSV, Markdown, D2, Mermaid — no format-specific code paths.",
  },
  {
    icon: "shield",
    title: "Type-Safe Everything",
    desc: "Format, ColorMode, ActivityStatus — validated at parse time. Branded IDs prevent mixing types at compile time.",
  },
  {
    icon: "feather",
    title: "Zero Heavy Deps in Root",
    desc: "go get pulls only x/term. YAML, lipgloss, bubbletea, and diagram renderers are opt-in sub-modules.",
  },
  {
    icon: "activity",
    title: "NOM Real-Time Progress",
    desc: "Dependency trees, activity counts, timing estimates, and inline terminal rendering. O(1) summary bars.",
  },
  {
    icon: "stream",
    title: "Streaming for Large Data",
    desc: "CQRS pure-function renderers stream via standard encoders. No intermediate string allocations.",
  },
  {
    icon: "palette",
    title: "Zero-Config Color",
    desc: "ColorModeAuto detects TTY, respects NO_COLOR, CI, FORCE_COLOR. Three wiring mechanisms, one result.",
  },
];
