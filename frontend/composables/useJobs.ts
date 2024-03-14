export const useJobs = defineStore('useJobs', () => {
	const running = ref([])
	const progress = ref([])

	function scheduleIsRunning(id: string) {
		const j = running.value?.find((job: any) => job.id === id)

		if (j) {
			if (j['out']['running'] !== undefined) {
				return j['out']['running']
			}
			return true
		}
		return false
	}

	function scheduleProgress(id: string): any | null {
		const job: any = running.value?.find((job: any) => job.id === id)
		if (job) return job.out
		return null
	}

	function repoIsRunning(id: string) {
		// const j = running.value?.find((job: any) => job.schedule.to_repository_id === id)
		// if (j) {
		// 	if (j['out']['running'] !== undefined) {
		// 		return j['out']['running']
		// 	}
		// 	return true
		// }
		return false
	}
	function repoIsSynching(id: string) {
		// const j = running.value?.find((job: Schedule) => job.schedule.from_repository_id === id)
		// if (j) {
		// 	if (j['out']['running'] !== undefined) {
		// 		return j['out']['running']
		// 	}
		// 	return true
		// }
		return false
	}

	function backupIsRunning(id: string) {
		// const j = running.value?.find((job: Schedule) => job.schedule.backup_id === id)
		// if (j) {
		// 	if (j['out']['running'] !== undefined) {
		// 		return j['out']['running']
		// 	}
		// 	return true
		// }
		return false
	}

	return {
		running,
		progress,
		scheduleIsRunning,
		scheduleProgress,
		repoIsRunning,
		repoIsSynching,
		backupIsRunning,
	}
})
