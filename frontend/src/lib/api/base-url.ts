import { PUBLIC_API_BASE_URL } from '$env/static/public';

const publicApiBaseUrl = PUBLIC_API_BASE_URL.trim();

export function apiUrl(path: string): string {
	if (!publicApiBaseUrl) {
		return path;
	}

	const normalizedBase = publicApiBaseUrl.replace(/\/+$/, '');
	const normalizedPath = path.startsWith('/') ? path : `/${path}`;
	return `${normalizedBase}${normalizedPath}`;
}
