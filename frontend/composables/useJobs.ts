export const useJobs = defineStore('useJobs', () => {
	const running = ref([])
	const progress = ref([])

	function scheduleIsRunning(id: string) {
		return running.value?.find((job: any) => job.id === id) ? true : false
	}

	function scheduleProgress(id: string): any | null {
		const job: any = running.value?.find((job: any) => job.id === id)
		if (job) return job.out
		return null
	}

	function repoIsRunning(id: string) {
		return running.value?.find((job: Schedule) => job.schedule.to_repository_id === id) ? true : false
	}
	function repoIsSynching(id: string) {
		return running.value?.find((job: Schedule) => job.schedule.from_repository_id === id) ? true : false
	}

	function backupIsRunning(id: string) {
		return running.value?.find((job: Schedule) => job.schedule.backup_id === id) ? true : false
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
