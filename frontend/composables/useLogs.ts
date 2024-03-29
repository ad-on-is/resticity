type Log = {
	[id: string]: string[]
}

export const useLogs = defineStore('useLogs', () => {
	const out = ref<Log>({})
	const err = ref<Log>({})
	const serverErr = ref<string[]>([])

	function setOut(id: string, data: string) {
		if (out.value[id] === undefined) {
			out.value[id] = []
		}
		if (data !== '') {
			out.value[id].push(data)
		}
	}
	function setErr(id: string, data: string) {
		if (err.value[id] === undefined) {
			err.value[id] = []
		}
		if (data !== '') {
			err.value[id].push(data)
		}
	}

	function setServerError(data: string) {
		serverErr.value.push(data)
	}

	return {
		out,
		err,
		setOut,
		setErr,
		serverErr,
		setServerError,
	}
})
