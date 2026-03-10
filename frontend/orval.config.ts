import { defineConfig } from 'orval';

export default defineConfig({
	postsApi: {
		input: {
			target: '../backend/docs/swagger.json'
		},
		output: {
			target: 'src/lib/api/generated/posts.ts',
			client: 'svelte-query',
			httpClient: 'fetch',
			clean: true,
			override: {
				query: {
					useQuery: true,
					useMutation: true
				}
			}
		}
	}
});
