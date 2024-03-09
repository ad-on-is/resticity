export const useSocket = defineStore('useSocket', () => {
	function init() {
		const getUrl = (): string => {
			const url = useRequestURL()
			return url.host.includes('wails.localhost') ? 'ws://localhost:11278' : `${url.protocol === 'http:' || url.protocol === 'wails:' ? 'ws:' : 'wss:'}//${url.host}`
		}
		const socket = new WebSocket(`${getUrl()}/api/ws`)
		socket.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data)
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
