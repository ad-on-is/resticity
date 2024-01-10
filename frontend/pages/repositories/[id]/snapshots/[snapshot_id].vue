<template>
	<div>
		<h3 class="text-purple-500 mb-3"><FaIcon icon="table-list" class="mr-2" />Browsing snapshot: {{ useRoute().params.snapshot_id }}</h3>
		<UButtonGroup class="mb-5" size="xs">
			<UButton color="indigo" :disabled="history.length === 0" icon="i-heroicons-chevron-left" @click="back"></UButton>
			<UButton color="gray" disabled icon="i-heroicons-folder">{{ path }}</UButton>
			<UInput v-model="subpath" placeholder="/" color="indigo" />
			<UButton color="indigo" @click="setPath(`${path}${subpath}`)">Go to</UButton>
			<UCheckbox v-model="showHidden" color="indigo" label="Show hidden" class="ml-15" />
		</UButtonGroup>
		<UTable :ui="{ td: { padding: 'py-1' } }" :rows="rows" :columns="columns" @select="" :loading="loading" class="bg-gray-950 rounded-xl bg-opacity-50 shadow-lg">
			<template #type-data="{ row }">
				<UIcon
					:class="row.name.startsWith('.') ? 'opacity-30' : ''"
					:name="row.type === 'dir' ? 'i-heroicons-folder' : 'i-heroicons-document'"
					:color="row.type === 'dir' ? 'yellow' : 'white'"
			/></template>
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
	const subpath = ref('')
	const path = ref('')
	const filesdirs = ref([])
	const loading = ref(false)
	const showHidden = ref(false)
	const cache = ref({})
	const setPath = (newPath: string) => {
		console.log('SET PATH')
		history.value.push(path.value)
		path.value = newPath
	}

	const back = () => {
		console.log('BACK')
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
				label: 'Browse',
				icon: 'i-heroicons-document-magnifying-glass',
				click: () => {
					console.log('browse')
				},
			},
		],
	]

	watch(path, async () => {
		loading.value = true
		// if (!cache.value[path.value]) {
		// 	cache.value[path.value] = []
		// 	console.log(cache.value)
		// }
		console.log('HERE')
		// if (typeof cache.value[path.value] !== 'undefined' && cache.value[path.value] !== null && cache.value[path.value].length > 0) {
		// 	console.log('HERE')
		// 	filesdirs.value = cache.value[path.value]
		// } else {

		const res = (await useApi().browseSnapshot(useRoute().params.id as string, useRoute().params.snapshot_id as string, path.value)) as []

		// cache.value[path.value] = res.data.value
		filesdirs.value = res
		// }

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
		console.log('MOUNTED')
		path.value = useRoute().query.path as string
	})
</script>
