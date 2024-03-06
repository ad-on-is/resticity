<template>
	<div class="p-5">
		<h3 class="text-purple-500 mb-3">{{ props.snapshotId }}</h3>
		<div class="flex justify-between">
			<UButtonGroup class="mb-5" size="xs">
				<UButton color="indigo" :disabled="history.length === 0" icon="i-heroicons-chevron-left" @click="back"></UButton>
				<UButton color="gray" disabled icon="i-heroicons-folder">{{ path }}</UButton>
			</UButtonGroup>
			<div class="ml-5 pt-1"><UCheckbox v-model="showHidden" color="indigo" label="Show hidden" /></div>
		</div>

		<UTable :ui="{ td: { padding: 'py-1' } }" :rows="rows" :columns="columns" @select="" :loading="loading" class="bg-gray-950 rounded-xl bg-opacity-50 shadow-lg">
			<template #type-data="{ row }">
				<span :class="row.type === 'dir' ? 'text-yellow-500' : 'text-white'"
					><UIcon :class="row.name.startsWith('.') ? 'opacity-40' : ''" :name="row.type === 'dir' ? 'i-heroicons-folder' : 'i-heroicons-document'" /></span
			></template>
			<template #name-data="{ row }">
				<div @click="setPath(row.path)">{{ row.name }}</div>
			</template>
			<template #mtime-data="{ row }"
				><div class="text-right text-xs">{{ formatISO9075(new Date(row.mtime)) }}</div></template
			>
			<template #size-data="{ row }"
				><div class="text-right text-xs">{{ humanFileSize(row.size) }}</div></template
			>
			<template #actions-data="{ row }">
				<UDropdown :items="items(row)"> <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" /> </UDropdown
			></template>
		</UTable>
	</div>
</template>

<script setup lang="ts">
	import { formatISO9075 } from 'date-fns'
	import _ from 'lodash'
	const history = ref<Array<string>>([])
	const path = ref('')
	const filesdirs = ref([])
	const loading = ref(false)
	const showHidden = ref(false)
	const setPath = (newPath: string) => {
		history.value.push(path.value)
		path.value = newPath
	}

	const props = defineProps({
		path: {
			type: String,
			required: true,
		},
		repositoryId: {
			type: String,
			required: true,
		},
		snapshotId: {
			type: String,
			required: true,
		},
	})

	const back = () => {
		path.value = history.value.pop() as string
	}
	const columns = [
		{ key: 'type', class: 'w-4' },
		{ key: 'name', label: 'Name', class: '' },
		{ key: 'mtime', label: 'Modified', class: 'w-30' },
		{ key: 'size', label: 'Size', class: 'w-30' },
		{ key: 'actions', class: 'w-10' },
	]
	const items = (row: any) => [
		[
			{
				label: 'Restore',
				disabled: true,
			},
		],
		[
			{
				label: 'Select folder',
				icon: 'i-heroicons-folder',
				click: async () => {
					const dir = await SelectDirectory('Select a repository')
					if (!dir) return
					restore(row.path, dir)
				},
			},
			{
				label: 'Replace original',
				icon: 'i-heroicons-document-duplicate',
				click: async () => {
					// todo: windows paths
					const dir = row.path.split('/').slice(0, -1).join('/')
					restore(row.path, dir)
				},
			},
		],
	]

	function restore(from: string, to: string) {
		useApi().restoreFromSnapshot(props.repositoryId, props.snapshotId, props.path, from, to)
	}

	watch(path, async () => {
		loading.value = true

		const res = (await useApi().browseSnapshot(props.repositoryId, props.snapshotId, path.value)) as []

		filesdirs.value = res

		loading.value = false
	})

	const rows = computed(() => {
		if (filesdirs.value.length === 0) return []
		const dirs = filesdirs.value.filter((file: any) => file.type === 'dir')
		dirs.shift()
		const files = filesdirs.value.filter((file: any) => file.type === 'file')
		const items = [...dirs, ...files]
		return items.filter((item: any) => {
			if (showHidden.value) return true
			return !item.name.startsWith('.')
		})
	})

	onMounted(() => {
		path.value = props.path
		// path.value = useRoute().query.path as string
	})
</script>
