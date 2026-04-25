import type { PageLoad } from './$types';

export const ssr = false;
export const prerender = true;

export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch('/api/profile-summary');
	if (!res.ok) {
		return {
			ok: false,
			summary: {
				headline: 'Ghassen El Kamel',
				subtitle: 'IoT & Systems Engineer',
				about: 'Profile service unavailable.',
				email: 'Ghassenelkamel@live.fr',
				portfolio: 'https://ghassenelkamel.fr'
			}
		};
	}
	return (await res.json()) as any;
};
