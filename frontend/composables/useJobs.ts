export const useJobs = defineStore('useJobs', () => {
	const running = ref<Awaited<ReturnType<typeof GetBackupJobs>>>([])

	function scheduleIsRunning(id: string) {
		return running.value?.find((job: BackupJob) => job.schedule.id === id) ? true : false
	}

	function repoIsRunning(id: string) {
		return running.value?.find((job: BackupJob) => job.schedule.to_repository_id === id) ? true : false
	}
	function repoIsSynching(id: string) {
		return running.value?.find((job: BackupJob) => job.schedule.from_repository_id === id) ? true : false
	}

	function backupIsRunning(id: string) {
		return running.value?.find((job: BackupJob) => job.schedule.backup_id === id) ? true : false
	}

	return {
		running,
		scheduleIsRunning,
		repoIsRunning,
		repoIsSynching,
		backupIsRunning,
	}
})
