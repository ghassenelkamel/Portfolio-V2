export type SkillItem = {
  name: string;
  pct: number;
};

export const skillsData: SkillItem[] = [
  { name: "Go (Backend)", pct: 0.88 },
  { name: "MQTT / VPN / TLS", pct: 0.86 },
  { name: "Linux / Unix", pct: 0.86 },
  { name: "PostgreSQL", pct: 0.82 },
  { name: "Docker", pct: 0.78 },
  { name: "JavaScript / Svelte", pct: 0.78 },
  { name: "C / C++", pct: 0.72 },
  { name: "Python", pct: 0.70 }
];

export function avgSkill(skills: SkillItem[]): number {
  if (!skills.length) return 0;
  const total = skills.reduce((sum, x) => sum + x.pct, 0);
  return total / skills.length;
}

export function topSkill(skills: SkillItem[]): SkillItem | null {
  if (!skills.length) return null;
  return [...skills].sort((a, b) => b.pct - a.pct)[0] ?? null;
}
