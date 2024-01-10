export const useApi = defineStore('useApi', () => {
	const browseSnapshot = async (repoId: string, snapshotId: string, path: string) =>
		(await useHttp.post(`/repositories/${repoId}/snapshots/${snapshotId}/browse`, { path: path })) ?? []
	const getSnapshots = async (repoId: string) => (await useHttp.get(`/repositories/${repoId}/snapshots`)) ?? []
	const runSchedule = async (scheduleId: string) => (await useHttp.get(`/schedules/${scheduleId}/run`)) ?? {}
	const getConfig = async () => (await useHttp.get(`/config`)) ?? {}
	const saveConfig = async (config: any) => (await useHttp.post(`/config`, config, { title: 'Settings', text: 'Settings saved successfully' })) ?? {}
	const checkRepository = async (repo: any) => (await useHttp.post(`/check`, repo, { title: 'Check Repository', text: 'Repository can be used' })) ?? {}
	const initRepository = async (repo: any) => (await useHttp.post(`/init`, repo, { title: 'Init Repository', text: 'Repository initialized' })) ?? {}
	return {
		browseSnapshot,
		getSnapshots,
		runSchedule,
		getConfig,
		saveConfig,
		checkRepository,
		initRepository,
	}
})
