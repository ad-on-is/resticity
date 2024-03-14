// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	ssr: false,
	modules: ['@pinia/nuxt', '@nuxt/ui', '@nuxtjs/color-mode'],
	devtools: { enabled: true },
	nitro: { static: true },

	css: ['@fortawesome/fontawesome-svg-core/styles.css'],
	imports: {
		dirs: ['wailsjs/**/*'],
	},

	app: {
		pageTransition: { name: 'page', mode: 'out-in' },
		head: {
			script: [{ src: 'https://unpkg.com/cronstrue@latest/dist/cronstrue.min.js' }],
		},
	},
})
