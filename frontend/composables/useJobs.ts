export const useJobs = defineStore('useJobs', () => {
	const running = ref<Awaited<ReturnType<typeof GetBackupJobs>>>([])
	function init() {
		// const interval = setInterval(async () => {
		// 	running.value = await GetBackupJobs()
		// }, 300)
	}

	function repoIsRunning(id: string) {
		return running.value?.find((job: BackupJob) => job.repository_id === id) ? true : false
	}

	function backupIsRunning(id: string) {
		return running.value?.find((job: BackupJob) => job.backup_id === id) ? true : false
	}

	return {
		running,
		init,
		repoIsRunning,
		backupIsRunning,
	}
})
