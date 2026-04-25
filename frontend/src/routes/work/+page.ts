import type { PageLoad } from './$types';

export const ssr = false;
export const prerender = true;

export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch('/api/github-repos?user=ghassenelkamel&limit=100');
	if (!res.ok) {
		const txt = await res.text().catch(() => '');
		return { repos: [], error: `API error: ${res.status} ${res.statusText}${txt ? ` — ${txt.slice(0, 200)}` : ''}` };
	}

	const data = (await res.json()) as { repos: any[]; error?: string | null };
	return { repos: data.repos ?? [], error: data.error ?? '' };
};
