import * as v from 'valibot';

import type { ApiCreatePostRequest } from '$lib/api/generated/posts';

export type CreatePostFormErrors = {
	title?: string;
	content?: string;
	form?: string;
};

export const createPostSchema = v.object({
	title: v.pipe(v.string(), v.trim(), v.nonEmpty('Title is required.')),
	content: v.pipe(v.string(), v.trim(), v.nonEmpty('Content is required.')),
	coauthor: v.optional(v.pipe(v.string(), v.trim(), v.nonEmpty('Coauthor cannot be empty.')))
}) satisfies v.GenericSchema<ApiCreatePostRequest>;

export function validateCreatePost(input: ApiCreatePostRequest) {
	const result = v.safeParse(createPostSchema, input);

	if (result.success) {
		return {
			success: true as const,
			data: result.output,
			errors: {}
		};
	}

	const flattened = v.flatten<typeof createPostSchema>(result.issues);

	return {
		success: false as const,
		data: null,
		errors: {
			title: flattened.nested?.title?.[0],
			content: flattened.nested?.content?.[0],
			form: flattened.root?.[0]
		} satisfies CreatePostFormErrors
	};
}
