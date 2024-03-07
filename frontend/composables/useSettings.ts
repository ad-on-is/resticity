export const useSettings = defineStore('useSettings', () => {
	const settings = ref()
	async function init() {
		refresh()
	}
	async function refresh() {
		settings.value = await useApi().getConfig()
	}
	async function save() {
		await useApi().saveConfig(settings.value)
	}
	return {
		settings,
		init,
		save,
		refresh,
	}
})
