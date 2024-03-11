type Log = {
	[id: string]: string[]
}

export const useLogs = defineStore('useLogs', () => {
	const out = ref<Log>({})
	const err = ref<Log>({})

	function setOut(id: string, data: string) {
		if (out.value[id] === undefined) {
			out.value[id] = []
		}
		out.value[id].push(data)
	}
	function setErr(id: string, data: string) {
		if (err.value[id] === undefined) {
			err.value[id] = []
		}
		err.value[id].push(data)
	}

	return {
		out,
		err,
		setOut,
		setErr,
	}
})
