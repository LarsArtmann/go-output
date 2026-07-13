export const featureIconKeys = [
  "layers",
  "shield",
  "feather",
  "activity",
  "stream",
  "palette",
] as const;
export type FeatureIcon = (typeof featureIconKeys)[number];

export interface Feature {
  icon: FeatureIcon;
  title: string;
  desc: string;
}

export interface PhaseCard {
  step: string;
  stepColor: "accent" | "amber" | "violet";
  title: string;
  code: string;
  hint: string;
  hintColor: "accent" | "amber" | "violet";
  hintIcon: UseCaseIcon;
}

export type ComparisonVariant = "DIY" | "go-output" | "Heavy";

export interface ComparisonItem {
  variant: ComparisonVariant;
  price: string;
  pros: string[];
  cons: string[];
  accent: boolean;
}

export const useCaseIconKeys = [
  "terminal",
  "refresh",
  "document",
  "chart",
  "bolt",
  "check",
] as const;
export type UseCaseIcon = (typeof useCaseIconKeys)[number];

export interface UseCase {
  title: string;
  desc: string;
  icon: UseCaseIcon;
}

export const uiIconKeys = [
  "arrow-external",
  "arrow-right",
  "github",
  "menu",
  "close",
  "sun",
  "moon",
  "star",
] as const;
export type UIIcon = (typeof uiIconKeys)[number];

export type IconName = FeatureIcon | UseCaseIcon | UIIcon;
