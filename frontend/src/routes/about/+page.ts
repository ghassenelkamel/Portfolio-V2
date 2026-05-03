import type { PageLoad } from './$types';

export const ssr = false;
export const prerender = true;

export const load: PageLoad = async ({ fetch, url }) => {
	const lang = (url.searchParams.get('lang') || '').toLowerCase().startsWith('fr') ? 'fr' : 'en';
	const qs = lang === 'fr' ? '?lang=fr' : '';
	const fallback = {
		schemaVersion: 1,
		updatedAt: new Date().toISOString(),
		headline: 'Ghassen El Kamel',
		subtitle: 'IoT & Systems Engineer · Backend · Infrastructure',
		summary: [
			'I’m an IoT and Software Systems Engineer passionate about building secure, scalable, and efficient infrastructures for connected devices and platforms.',
			'My work sits at the intersection of embedded systems, backend development, and network security, transforming real-world hardware into reliable digital ecosystems.'
		],
		keywords: [
			'IoT architecture',
			'Go backend',
			'Embedded Linux',
			'MQTT / TLS',
			'VPN networking',
			'Automation APIs'
		],
		specialties: [
			'Designing IoT architectures with secure data flow from device to cloud',
			'Developing back-end services in Go for automation, monitoring, and data processing',
			'Integrating secure communication layers — VPN, MQTTS, and TLS'
		],
		currently: ['Building Go-based IoT backbones and secure APIs for data exchange'],
		goals: ['Contribute to open-source infrastructure and IoT projects'],
		email: 'Ghassenelkamel@live.fr',
		site: 'https://ghassenelkamel.fr',
		github: 'https://github.com/ghassenelkamel'
	};

	const res = await fetch(`/api/content/about${qs}`);
	if (!res.ok) {
		return { ok: false, content: fallback };
	}

	const data = (await res.json()) as { ok?: boolean; data?: unknown };
	return { ok: !!data?.ok, content: data?.data ?? fallback };
};
