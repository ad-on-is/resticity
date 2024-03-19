export const useMounts = defineStore('useMounts', () => {
	const mounts = ref([])
	const progress = ref([])

	function repoIsMounted(id: string) {
		return mounts.value?.find((repo: any) => repo.id === id)
	}

	return {
		mounts,
		repoIsMounted,
	}
})
