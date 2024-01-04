export const useSettings = defineStore('useSettings', () => {
	const settings = ref<Awaited<ReturnType<typeof Settings>>>()
	async function init() {
		settings.value = await Settings()
	}
	async function save() {
		console.log('SHOULD SAVE')
		await SaveSettings(settings.value!)
	}
	return {
		settings,
		init,
		save,
	}
})
