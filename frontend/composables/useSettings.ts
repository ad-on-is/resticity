export const useSettings = defineStore('useSettings', () => {
	const settings = ref()
	async function init() {
		const res = await useFetch('http://127.0.0.1:11278/api/config', { method: 'GET' })
		settings.value = res.data.value
		console.log(settings.value)
	}
	async function save() {
		console.log('SHOULD SAVE')
		// await SaveSettings(settings.value!)
		await useFetch('http://localhost:11278/api/config', { method: 'POST', body: settings.value })
	}
	return {
		settings,
		init,
		save,
	}
})
