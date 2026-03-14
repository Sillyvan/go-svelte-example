# sv

Everything you need to build a Svelte project, powered by [`sv`](https://github.com/sveltejs/cli).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```sh
# create a new project
npx sv create my-app
```

To recreate this project with the same configuration:

```sh
# recreate this project
npx sv@0.12.5 create --template minimal --types ts --install npm frontend
```

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```sh
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```sh
npm run build
```

For a deployed frontend that talks to the deployed backend, set `PUBLIC_API_BASE_URL` at build time:

```sh
PUBLIC_API_BASE_URL=https://go-svelte-example-backend.jubati911.workers.dev npm run build
```

For Cloudflare Pages / Workers static hosting:

- Root directory: `frontend`
- Build command: `npm run build`
- Output directory: `build`
- Node version: `22.18.0` or newer
- Environment variable: `PUBLIC_API_BASE_URL=https://go-svelte-example-backend.jubati911.workers.dev`

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.
