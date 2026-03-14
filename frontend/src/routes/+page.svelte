<script lang="ts">
	import {
		createGetPosts,
		createPostPosts,
		getGetPostsQueryKey,
		type ApiCreatePostRequest,
		type ApiErrorResponse
	} from '$lib/api/generated/posts';
	import { queryClient } from '$lib/query-client';
	import { validateCreatePost, type CreatePostFormErrors } from '$lib/validation/posts';

	let title = $state('');
	let content = $state('');
	let coauthor = $state('');
	let fieldErrors = $state<CreatePostFormErrors>({});
	let submitError = $state<string | null>(null);

	const postsQuery = createGetPosts();
	const createPostMutation = createPostPosts();

	function getErrorMessage(error: unknown, fallback: string) {
		if (error && typeof error === 'object' && 'message' in error && typeof error.message === 'string') {
			return error.message;
		}

		return fallback;
	}

	function getApiMessage(response: ApiErrorResponse, fallback: string) {
		return response.message.trim() || fallback;
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		fieldErrors = {};
		submitError = null;
		const payload: ApiCreatePostRequest = {
			title,
			content,
			...(coauthor.trim() === '' ? {} : { coauthor })
		};
		const validation = validateCreatePost(payload);

		if (!validation.success) {
			fieldErrors = validation.errors;
			submitError = validation.errors.form ?? null;
			return;
		}

		try {
			const response = await createPostMutation.mutateAsync({
				data: validation.data
			});

			if (response.status !== 201) {
				submitError = getApiMessage(response.data, 'Failed to create post.');
				return;
			}

			title = '';
			content = '';
			coauthor = '';
			fieldErrors = {};
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

<h2>Create Post</h2>

<form onsubmit={handleSubmit}>
	<p>
		<label for="title">Title</label><br />
		<input id="title" name="title" bind:value={title} />
	</p>
	{#if fieldErrors.title}
		<p>{fieldErrors.title}</p>
	{/if}

	<p>
		<label for="content">Content</label><br />
		<textarea id="content" name="content" rows="6" bind:value={content}></textarea>
	</p>
	{#if fieldErrors.content}
		<p>{fieldErrors.content}</p>
	{/if}

	<p>
		<label for="coauthor">Coauthor</label><br />
		<input id="coauthor" name="coauthor" bind:value={coauthor} />
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
{:else if !postsQuery.data}
	<p>Failed to load posts.</p>
{:else if postsQuery.data.status !== 200}
	<p>{getApiMessage(postsQuery.data.data, 'Failed to load posts.')}</p>
{:else}
	<ul>
		{#each postsQuery.data.data as post (post.id)}
			<li>
				<p><a href="/{post.id}">{post.title}</a></p>
				<p>{post.content}</p>
				{#if post.coauthor}
					<p>Coauthor: {post.coauthor}</p>
				{/if}
				<p>{post.created_at}</p>
			</li>
		{:else}
			<li>No posts yet.</li>
		{/each}
	</ul>
{/if}
