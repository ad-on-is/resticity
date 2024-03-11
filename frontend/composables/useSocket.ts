export const useSocket = defineStore('useSocket', () => {
	function init() {
		const getUrl = (): string => {
			const url = useRequestURL()
			return url.protocol === 'wails:' ? 'ws://localhost:11278' : `${url.protocol === 'http:' ? 'ws:' : 'wss:'}//${url.host}`
		}
		const socket = new WebSocket(`${getUrl()}/api/ws`)
		socket.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data)
				data.forEach((j: any) => {
					if (j.out !== undefined) {
						useLogs().setOut(j.id, j.out)
					}
					if (j.err !== undefined) {
						useLogs().setErr(j.id, j.err)
					}
				})
				useJobs().running =
					data.map((j: any) => {
						try {
							j.out = JSON.parse(j.out)
						} catch {
							j.out = {}
						}
						return j
					}) || []
			} catch (e) {
				useJobs().running = []
				console.error(e)
			}
		}

		const interval = setInterval(() => {
			if (socket.readyState === 1) {
				socket.send('ping')
			}
		}, 1000)
	}

	return {
		init,
	}
})
