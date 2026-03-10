<script lang="ts">
	import {
		createGetPosts,
		createPostPosts,
		getGetPostsQueryKey,
		type ApiErrorResponse,
		type StorePost
	} from '$lib/api/generated/posts';
	import { queryClient } from '$lib/query-client';

	let title = $state('');
	let content = $state('');
	let submitError = $state<string | null>(null);

	const postsQuery = createGetPosts();
	const createPostMutation = createPostPosts();

	function getErrorMessage(error: unknown, fallback: string) {
		if (error && typeof error === 'object' && 'message' in error && typeof error.message === 'string') {
			return error.message;
		}

		return fallback;
	}

	function getApiMessage(response: { data?: ApiErrorResponse } | undefined, fallback: string) {
		return response?.data?.message?.trim() || fallback;
	}

	function getPostTitle(post: StorePost) {
		return post.title?.trim() || `Post ${post.id ?? ''}`.trim();
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		const nextTitle = title.trim();
		const nextContent = content.trim();

		if (!nextTitle || !nextContent) {
			submitError = 'Title and content are required.';
			return;
		}

		submitError = null;

		try {
			const response = await createPostMutation.mutateAsync({
				data: {
					title: nextTitle,
					content: nextContent
				}
			});

			if (response.status !== 201) {
				submitError = getApiMessage(response, 'Failed to create post.');
				return;
			}

			title = '';
			content = '';
			await queryClient.invalidateQueries({ queryKey: getGetPostsQueryKey() });
		} catch (error) {
			submitError = getErrorMessage(error, 'Failed to create post.');
		}
	}
</script>

<svelte:head>
	<title>Posts</title>
</svelte:head>

<h1>Posts</h1>

<p><a href="/swagger">Swagger UI</a></p>

<h2>Create Post</h2>

<form onsubmit={handleSubmit}>
	<p>
		<label for="title">Title</label><br />
		<input id="title" name="title" bind:value={title} />
	</p>

	<p>
		<label for="content">Content</label><br />
		<textarea id="content" name="content" rows="6" bind:value={content}></textarea>
	</p>

	<p>
		<button type="submit" disabled={createPostMutation.isPending}>Create</button>
	</p>
</form>

{#if createPostMutation.isPending}
	<p>Creating post...</p>
{/if}

{#if submitError}
	<p>{submitError}</p>
{/if}

<h2>All Posts</h2>

{#if postsQuery.isPending}
	<p>Loading posts...</p>
{:else if postsQuery.isError}
	<p>{getErrorMessage(postsQuery.error, 'Failed to load posts.')}</p>
{:else if postsQuery.data?.status !== 200}
	<p>{getApiMessage(postsQuery.data, 'Failed to load posts.')}</p>
{:else}
	<ul>
		{#each postsQuery.data.data as post (post.id)}
			<li>
				<p><a href="/{post.id}">{getPostTitle(post)}</a></p>
				<p>{post.content}</p>
				<p>{post.created_at ?? 'Unknown date'}</p>
			</li>
		{:else}
			<li>No posts yet.</li>
		{/each}
	</ul>
{/if}
