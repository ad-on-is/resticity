// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	ssr: false,
	modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
	devtools: { enabled: true },
	nitro: { static: true },

	css: ['@fortawesome/fontawesome-svg-core/styles.css'],
	imports: {
		dirs: ['wailsjs/go/**/*'],
	},

	app: {
		pageTransition: { name: 'page', mode: 'out-in' },
	},
})
