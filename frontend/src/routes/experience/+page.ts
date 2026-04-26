import type { PageLoad } from './$types';

export const ssr = false;
export const prerender = true;

export const load: PageLoad = async ({ fetch }) => {
	const fallback = {
		schemaVersion: 1,
		updatedAt: new Date().toISOString(),
		eyebrow: 'Experience',
		title: 'Professional Journey',
		description: 'Select an experience to explore scope, impact, and tools used.',
		items: []
	};

	const res = await fetch('/api/content/experience');
	if (!res.ok) {
		return { ok: false, content: fallback };
	}

	const data = (await res.json()) as { ok?: boolean; data?: unknown };
	return { ok: !!data?.ok, content: data?.data ?? fallback };
};
