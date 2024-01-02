export const useSettings = defineStore('useSettings', () => {
	const settings = ref<Awaited<ReturnType<typeof Settings>>>()
	async function init() {
		settings.value = await Settings()
		console.log('SETTINGS', settings)
	}
	async function save() {
		await SaveSettings(settings.value!)
	}
	return {
		settings,
		init,
		save,
	}
})
