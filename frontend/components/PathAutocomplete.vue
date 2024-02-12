<template>
	<UButtonGroup class="flex">
		<!-- <UInput v-model="path" placeholder="/" class="flex-grow" /> -->
		<USelectMenu v-model="selected" :searchable="search" class="flex-grow" searchable-placeholder="Type to autocomplete">
			<template #label>
				<span v-if="path !== ''">{{ `${path}` }}</span
				><span v-else>Select a directory</span>
			</template>
		</USelectMenu>
		<UButton icon="i-heroicons-folder-open" color="indigo" @click="openDir()" />
	</UButtonGroup>
</template>

<script setup lang="ts">
	const loading = ref(false)
	const emit = defineEmits(['selected'])
	const path = ref('')
	const selected = ref('')
	let previousPaths: string[] = []
	async function search(q: string) {
		let paths = []
		if (q.endsWith('/') || q.endsWith('\\')) {
			path.value = q
			loading.value = true
			const res = await useApi().autoCompletePath(q)
			previousPaths = res || []
			paths = res || []
		} else {
			const last = q.split(/\/|\\/).pop()
			paths = previousPaths.filter((p) => p.startsWith(last as string))
		}

		loading.value = false
		return paths
	}

	watch(selected, (v) => {
		if (v !== '') {
			path.value = path.value + v
		}
	})

	watch(path, (p) => {
		if (p !== '') {
			emit('selected', p)
		}
	})

	const openDir = async () => {
		try {
			const dir = await SelectDirectory('Select a directory')
			if (dir !== '') {
				path.value = dir
			}
		} catch (e) {
			console.log('Not supported in browser')
		}
	}
</script>
