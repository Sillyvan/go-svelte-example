<script lang="ts">
	import { page } from '$app/state';
	import { createGetPostsId, type ApiErrorResponse } from '$lib/api/generated/posts';

	let postId = $derived((page.params.postId ?? '').trim());

	const postQuery = createGetPostsId(() => postId);

	function getErrorMessage(error: unknown, fallback: string) {
		if (error && typeof error === 'object' && 'message' in error && typeof error.message === 'string') {
			return error.message;
		}

		return fallback;
	}

	function getApiMessage(response: ApiErrorResponse, fallback: string) {
		return response.message.trim() || fallback;
	}
</script>

<svelte:head>
	<title>Post</title>
</svelte:head>

<p><a href="/">Back to posts</a></p>

<h1>Post Detail</h1>

{#if postId === ''}
	<p>Invalid post id.</p>
{:else if postQuery.isPending}
	<p>Loading post...</p>
{:else if postQuery.isError}
	<p>{getErrorMessage(postQuery.error, 'Failed to load post.')}</p>
{:else if !postQuery.data}
	<p>Failed to load post.</p>
{:else if postQuery.data.status !== 200}
	<p>{getApiMessage(postQuery.data.data, 'Failed to load post.')}</p>
{:else}
	<article>
		<h2>{postQuery.data.data.title}</h2>
		<p>{postQuery.data.data.content}</p>
		{#if postQuery.data.data.coauthor}
			<p>Coauthor: {postQuery.data.data.coauthor}</p>
		{/if}
		<p>{postQuery.data.data.created_at}</p>
		<p>ID: {postQuery.data.data.id}</p>
	</article>
{/if}
