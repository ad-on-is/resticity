import _ from 'lodash'

export const useApi = defineStore('useApi', () => {
	const browseSnapshot = async (repoId: string, snapshotId: string, path: string): Promise<FileDescriptor[]> =>
		(await useHttp.post(`/repositories/${repoId}/snapshots/${snapshotId}/browse`, { path: path })) ?? []
	const restoreFromSnapshot = async (repoId: string, snapshotId: string, rootPath: string, fromPath: string, toPath: string) =>
		(await useHttp.post(
			`/repositories/${repoId}/snapshots/${snapshotId}/restore`,
			{ root_path: rootPath, from_path: fromPath.replaceAll(rootPath, ''), to_path: toPath },
			{ title: 'Restoring', text: 'Successfully restored' }
		)) ?? []
	const getSnapshots = async (repoId: string, groupBy: string = 'host'): Promise<SnapshotGroup[]> => {
		const data = (await useHttp.post(`/repositories/${repoId}/snapshots?group_by=${groupBy}`)) ?? []
		console.log(data)
		return _.orderBy(data, ['time'], ['desc'])
	}
	const mount = async (repoId: string, path: string) => (await useHttp.post(`/repositories/${repoId}/mount`, { path: path }, { title: 'Mount', text: `Mounted to ${path}` })) ?? {}
	const unmount = async (repoId: string, path: string) =>
		(await useHttp.post(`/repositories/${repoId}/unmount`, { path: path }, { title: 'Unmount', text: `Unmounted: ${path}` })) ?? {}

	const statRepository = async (repoId: string) => (await useHttp.get(`/repositories/${repoId}/stats`)) ?? {}
	const runSchedule = async (scheduleId: string) => (await useHttp.get(`/schedules/${scheduleId}/run`)) ?? {}
	const stopSchedule = async (scheduleId: string) => (await useHttp.get(`/schedules/${scheduleId}/stop`)) ?? {}
	const getConfig = async (): Promise<Config> => (await useHttp.get(`/config`)) ?? {}
	const saveConfig = async (config: any) => (await useHttp.post(`/config`, config, { title: 'Settings', text: 'Settings saved successfully' })) ?? {}
	const checkRepository = async (repo: any) => (await useHttp.post(`/check`, repo, { title: 'Check Repository', text: 'Repository can be used' })) ?? {}
	const initRepository = async (repo: any) => (await useHttp.post(`/init`, repo, { title: 'Init Repository', text: 'Repository initialized' })) ?? {}

	return {
		browseSnapshot,
		restoreFromSnapshot,
		getSnapshots,
		runSchedule,
		stopSchedule,
		mount,
		unmount,
		getConfig,
		saveConfig,
		checkRepository,
		initRepository,
		statRepository,
	}
})
