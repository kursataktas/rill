import chroma, { type Color } from "chroma-js";

export const TailwindColorSpacing = [
  50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950,
];

// Pulled from https://tailwindcss.com/docs/customizing-colors
// Autogenerated using this script in console,
/*
  colors = {}
  document.querySelector(".grid").childNodes.forEach((c) => {
    const palette = [];
    for (const color of c.childNodes[1].childNodes) {
      palette.push(color.querySelector(".text-slate-500").textContent)
    }
    colors[c.childNodes[0].textContent] = palette
  })
  JSON.stringify(colors)
*/
const TailwindColorPresets: Record<string, string[]> = {
  Slate: [
    "#f8fafc",
    "#f1f5f9",
    "#e2e8f0",
    "#cbd5e1",
    "#94a3b8",
    "#64748b",
    "#475569",
    "#334155",
    "#1e293b",
    "#0f172a",
    "#020617",
  ],
  Gray: [
    "#f9fafb",
    "#f3f4f6",
    "#e5e7eb",
    "#d1d5db",
    "#9ca3af",
    "#6b7280",
    "#4b5563",
    "#374151",
    "#1f2937",
    "#111827",
    "#030712",
  ],
  Zinc: [
    "#fafafa",
    "#f4f4f5",
    "#e4e4e7",
    "#d4d4d8",
    "#a1a1aa",
    "#71717a",
    "#52525b",
    "#3f3f46",
    "#27272a",
    "#18181b",
    "#09090b",
  ],
  Neutral: [
    "#fafafa",
    "#f5f5f5",
    "#e5e5e5",
    "#d4d4d4",
    "#a3a3a3",
    "#737373",
    "#525252",
    "#404040",
    "#262626",
    "#171717",
    "#0a0a0a",
  ],
  Stone: [
    "#fafaf9",
    "#f5f5f4",
    "#e7e5e4",
    "#d6d3d1",
    "#a8a29e",
    "#78716c",
    "#57534e",
    "#44403c",
    "#292524",
    "#1c1917",
    "#0c0a09",
  ],
  Red: [
    "#fef2f2",
    "#fee2e2",
    "#fecaca",
    "#fca5a5",
    "#f87171",
    "#ef4444",
    "#dc2626",
    "#b91c1c",
    "#991b1b",
    "#7f1d1d",
    "#450a0a",
  ],
  Orange: [
    "#fff7ed",
    "#ffedd5",
    "#fed7aa",
    "#fdba74",
    "#fb923c",
    "#f97316",
    "#ea580c",
    "#c2410c",
    "#9a3412",
    "#7c2d12",
    "#431407",
  ],
  Amber: [
    "#fffbeb",
    "#fef3c7",
    "#fde68a",
    "#fcd34d",
    "#fbbf24",
    "#f59e0b",
    "#d97706",
    "#b45309",
    "#92400e",
    "#78350f",
    "#451a03",
  ],
  Yellow: [
    "#fefce8",
    "#fef9c3",
    "#fef08a",
    "#fde047",
    "#facc15",
    "#eab308",
    "#ca8a04",
    "#a16207",
    "#854d0e",
    "#713f12",
    "#422006",
  ],
  Lime: [
    "#f7fee7",
    "#ecfccb",
    "#d9f99d",
    "#bef264",
    "#a3e635",
    "#84cc16",
    "#65a30d",
    "#4d7c0f",
    "#3f6212",
    "#365314",
    "#1a2e05",
  ],
  Green: [
    "#f0fdf4",
    "#dcfce7",
    "#bbf7d0",
    "#86efac",
    "#4ade80",
    "#22c55e",
    "#16a34a",
    "#15803d",
    "#166534",
    "#14532d",
    "#052e16",
  ],
  Emerald: [
    "#ecfdf5",
    "#d1fae5",
    "#a7f3d0",
    "#6ee7b7",
    "#34d399",
    "#10b981",
    "#059669",
    "#047857",
    "#065f46",
    "#064e3b",
    "#022c22",
  ],
  Teal: [
    "#f0fdfa",
    "#ccfbf1",
    "#99f6e4",
    "#5eead4",
    "#2dd4bf",
    "#14b8a6",
    "#0d9488",
    "#0f766e",
    "#115e59",
    "#134e4a",
    "#042f2e",
  ],
  Cyan: [
    "#ecfeff",
    "#cffafe",
    "#a5f3fc",
    "#67e8f9",
    "#22d3ee",
    "#06b6d4",
    "#0891b2",
    "#0e7490",
    "#155e75",
    "#164e63",
    "#083344",
  ],
  Sky: [
    "#f0f9ff",
    "#e0f2fe",
    "#bae6fd",
    "#7dd3fc",
    "#38bdf8",
    "#0ea5e9",
    "#0284c7",
    "#0369a1",
    "#075985",
    "#0c4a6e",
    "#082f49",
  ],
  Blue: [
    "#eff6ff",
    "#dbeafe",
    "#bfdbfe",
    "#93c5fd",
    "#60a5fa",
    "#3b82f6",
    "#2563eb",
    "#1d4ed8",
    "#1e40af",
    "#1e3a8a",
    "#172554",
  ],
  Indigo: [
    "#eef2ff",
    "#e0e7ff",
    "#c7d2fe",
    "#a5b4fc",
    "#818cf8",
    "#6366f1",
    "#4f46e5",
    "#4338ca",
    "#3730a3",
    "#312e81",
    "#1e1b4b",
  ],
  Violet: [
    "#f5f3ff",
    "#ede9fe",
    "#ddd6fe",
    "#c4b5fd",
    "#a78bfa",
    "#8b5cf6",
    "#7c3aed",
    "#6d28d9",
    "#5b21b6",
    "#4c1d95",
    "#2e1065",
  ],
  Purple: [
    "#faf5ff",
    "#f3e8ff",
    "#e9d5ff",
    "#d8b4fe",
    "#c084fc",
    "#a855f7",
    "#9333ea",
    "#7e22ce",
    "#6b21a8",
    "#581c87",
    "#3b0764",
  ],
  Fuchsia: [
    "#fdf4ff",
    "#fae8ff",
    "#f5d0fe",
    "#f0abfc",
    "#e879f9",
    "#d946ef",
    "#c026d3",
    "#a21caf",
    "#86198f",
    "#701a75",
    "#4a044e",
  ],
  Pink: [
    "#fdf2f8",
    "#fce7f3",
    "#fbcfe8",
    "#f9a8d4",
    "#f472b6",
    "#ec4899",
    "#db2777",
    "#be185d",
    "#9d174d",
    "#831843",
    "#500724",
  ],
  Rose: [
    "#fff1f2",
    "#ffe4e6",
    "#fecdd3",
    "#fda4af",
    "#fb7185",
    "#f43f5e",
    "#e11d48",
    "#be123c",
    "#9f1239",
    "#881337",
    "#4c0519",
  ],
};

// Usable preset converted from the above generated one
export const TailwindColorPresetsConverted: Color[][] = Object.keys(
  TailwindColorPresets,
).map((colorName) =>
  TailwindColorPresets[colorName].map((shade) => chroma(shade)),
);

/**
 * This type represents a color palette, where the keys are
 * the color lightness numbers (50, 100, 200, etc) and the
 * values are css color strings.
 */
export type LightnessMap = { [key: number]: string };

/**
 * The three categories of colors that could be rethemed.
 * (though we only support primary and secondary for now)
 */
export type ThemeColorKind = "primary" | "secondary" | "muted";

/**
 * Rill primary brand colors.
 */
export const defaultPrimaryColors: LightnessMap = {
  50: "227 100% 96%",
  100: "228 100% 93%",
  200: "229 100% 88%",
  300: "231 100% 81%",
  400: "236 100% 73%",
  500: "240 100% 67%",
  600: "246 91% 58%",
  700: "235 58% 49%",
  800: "245 65% 41%",
  900: "244 57% 34%",
  950: "246 57% 20%",
};

// backup pallette of primary colors
// useful for testing application of colors
// export const defaultPrimaryColors = Object.fromEntries(
//   [50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950].map((n) => [
//     n,
//     `lch(${(100 * (1000 - n)) / 1000}% 64 139)`,
//   ]),
// );

/**
 * Rill secondary brand colors.
 */
export const defaultSecondaryColors = {
  50: "199 100% 97%",
  100: "202 100% 94%",
  200: "198 100% 86%",
  300: "196 100% 73%",
  400: "195 100% 59%",
  500: "197 100% 50%",
  600: "200 100% 42%",
  700: "200 100% 34%",
  800: "199 100% 28%",
  900: "200 90% 24%",
  950: "202 90% 16%",
};

// backup pallette of secondary colors (red spectrum)",
// useful for testing application of colors
// export const defaultSecondaryColors = Object.fromEntries(
//   [50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950].map((n) => [
//     n,
//     `lch(${(100 * (1000 - n)) / 1000}% 78 13)`,
//   ]),
// );

/*
 * colors for greyed-out elements. For now, using tailwind's
 * standard "Gray", but using semantic color vars will
 * allow us to change this to a custom palette
 * or use e.g. tailwind's "slate", "zinc", etc if we want.
 *
 * Copied from https://github.com/shadcn-ui/ui/issues/669#issue-1771280130
 *
 * Visit that link if we want to copy/paste if
 * we switch to "slate", "zinc", etc
 *
 */
export const mutedColors = {
  50: "210 20% 98%",
  100: "220 14.3% 95.9%",
  200: "220 13% 91%",
  300: "216 12.2% 83.9%",
  400: "217.9 10.6% 64.9%",
  500: "220 8.9% 46.1%",
  600: "215 13.8% 34.1%",
  700: "216.9 19.1% 26.7%",
  800: "215 27.9% 16.9%",
  900: "220.9 39.3% 11%",
  950: "224 71.4% 4.1%;",
};

export function getRandomBgColor(): string {
  const colorList = [
    "bg-blue-500",
    "bg-green-500",
    "bg-red-500",
    "bg-orange-500",
    "bg-yellow-500",
    "bg-amber-500",
    "bg-pink-500",
    "bg-lime-500",
    "bg-emerald-500",
    "bg-teal-500",
    "bg-cyan-500",
    "bg-sky-500",
    "bg-indigo-500",
    "bg-violet-500",
    "bg-purple-500",
    "bg-fuchsia-500",
    "bg-rose-500",
  ];
  return colorList[Math.floor(Math.random() * colorList.length)];
}
