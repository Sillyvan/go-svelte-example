<script lang="ts">
	import { page } from '$app/state';
	import { createGetPostsId, type ApiErrorResponse } from '$lib/api/generated/posts';

	let postId = $derived(Number(page.params.postId));

	const postQuery = createGetPostsId(() => postId);

	function getErrorMessage(error: unknown, fallback: string) {
		if (error && typeof error === 'object' && 'message' in error && typeof error.message === 'string') {
			return error.message;
		}

		return fallback;
	}

	function getApiMessage(response: { data?: ApiErrorResponse } | undefined, fallback: string) {
		return response?.data?.message?.trim() || fallback;
	}
</script>

<svelte:head>
	<title>Post</title>
</svelte:head>

<p><a href="/">Back to posts</a></p>

<h1>Post Detail</h1>

{#if Number.isNaN(postId) || postId < 1}
	<p>Invalid post id.</p>
{:else if postQuery.isPending}
	<p>Loading post...</p>
{:else if postQuery.isError}
	<p>{getErrorMessage(postQuery.error, 'Failed to load post.')}</p>
{:else if postQuery.data?.status !== 200}
	<p>{getApiMessage(postQuery.data, 'Failed to load post.')}</p>
{:else}
	<article>
		<h2>{postQuery.data.data.title}</h2>
		<p>{postQuery.data.data.content}</p>
		<p>{postQuery.data.data.created_at ?? 'Unknown date'}</p>
		<p>ID: {postQuery.data.data.id}</p>
	</article>
{/if}
