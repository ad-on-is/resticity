export const useSocket = defineStore('useSocket', () => {
	function init() {
		const jobsocket = new WebSocket('ws://127.0.0.1:11278/api/ws')
		jobsocket.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data)
				useJobs().running = data.jobs
				useLogs().out = data.out
				useLogs().err = data.err
				console.log(data)
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
