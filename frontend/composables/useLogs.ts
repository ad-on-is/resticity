export const useLogs = defineStore('useLogs', () => {
	const out = ref('')
	const err = ref('')

	return {
		out,
		err,
	}
})
