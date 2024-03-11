export const useLogs = defineStore('useLogs', () => {
	const out = ref('')
	const err = ref('')

	function setOut(data: string) {
		out.value += data
	}
	function setErr(data: string) {
		err.value += data
	}

	return {
		out,
		err,
		setOut,
		setErr,
	}
})
