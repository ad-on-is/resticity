export const useSocket = defineStore('useSocket', () => {
	function init() {
		const socket = new WebSocket('ws://127.0.0.1:11278/api/ws')
		socket.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data)
				console.log(data)
				useJobs().running =
					data.map((j: any) => {
						try {
							j.out = JSON.parse(j.out)
						} catch {
							j.out = {}
						}
						return j
					}) || []
				// useLogs().out = data.out
				// useLogs().err = data.err
			} catch (e) {
				useJobs().running = []
				console.error(e)
			}
		}
	}

	return {
		init,
	}
})
