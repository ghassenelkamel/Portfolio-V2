import type { PageLoad } from './$types';

export const ssr = false;
export const prerender = true;

export const load: PageLoad = async ({ fetch }) => {
	const fallback = {
		schemaVersion: 1,
		updatedAt: new Date().toISOString(),
		eyebrow: 'Contact',
		title: 'Professional Reach',
		description: 'Select a contact channel to discuss scope, impact, and collaboration.',
		email: 'Ghassenelkamel@live.fr',
		linkedin: 'https://linkedin.com/in/ghassenelkamel',
		github: 'https://github.com/ghassenelkamel',
		portfolio: 'https://ghassenelkamel.fr'
	};

	const res = await fetch('/api/content/contact');
	if (!res.ok) {
		return { ok: false, content: fallback };
	}

	const data = (await res.json()) as { ok?: boolean; data?: unknown };
	return { ok: !!data?.ok, content: data?.data ?? fallback };
};
