<template>
	<h3 class="text-purple-500 mb-3"><FaIcon icon="table-list" class="mr-2" />Snapshots</h3>

	<UTable :rows="snapshots" v-model="selected" :columns="columns" @select="" :loading="loading" class="bg-gray-950 rounded-xl bg-opacity-50 shadow-lg">
		<template #tags-data="{ row }">
			<UBadge v-for="tag in row.tags" variant="outline" color="indigo">{{ tag }}</UBadge>
		</template>
		<template #time-data="{ row }">
			{{ format(new Date(row.time), 'dd.MM.yyyy H:I:s') }}
		</template>
		<template #paths-data="{ row }">
			{{ row.paths.join(',') }}
		</template>
		<template #actions-data="{ row }">
			<UDropdown :items="items(row)">
				<UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" :disabled="selected.length > 0" />
			</UDropdown>
		</template>
	</UTable>
	<UButton
		icon="i-heroicons-trash"
		class="mt-3"
		:disabled="selected.length === 0"
		:color="selected.length === 0 ? 'gray' : 'indigo'"
		:variant="selected.length === 0 ? 'solid' : 'outline'"
		>Prune snapshots</UButton
	>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import { format } from 'date-fns'
	import _ from 'lodash'
	const snapshots = ref<Array<Snapshot>>([])
	const loading = ref(true)
	const selected = ref<Array<Snapshot>>([])
	const mounted = ref<Array<string>>([])

	const paths = ref<Array<string>>([])

	const items = (row: any) => [
		[
			...paths.value.map((p) => ({
				label: 'Browse ' + p,
				icon: 'i-heroicons-document-magnifying-glass',
				click: () => {
					navigateTo({ path: `/repositories/${useRoute().params.id}/snapshots/${row.id}`, query: { path: p } })
				},
			})),
			{
				label: mounted.value.includes(row.id) ? 'Unmount' : 'Mount',
				icon: mounted.value.includes(row.id) ? 'i-heroicons-server' : 'i-heroicons-server',
				click: () => {
					if (mounted.value.includes(row.id)) {
						mounted.value = mounted.value.filter((item) => item !== row.id)
						// Mount(row.id, false)
					} else {
						mounted.value.push(row.id)
						// Mount(row.id, true)
					}
				},
			},
		],
	]

	function select(row: Snapshot) {
		const index = selected.value.findIndex((item: Snapshot) => item.id === row.id)
		if (index === -1) {
			selected.value.push(row)
		} else {
			selected.value.splice(index, 1)
		}
	}
	const columns = [
		{
			key: 'short_id',
			label: 'ID',
		},
		{
			key: 'hostname',
			label: 'Host',
		},

		{
			key: 'time',
			label: 'Time',
			format: (value: string) => format(new Date(value), 'dd.MM.yyyy H:I:s'),
		},
		{
			key: 'paths',
			label: 'Paths',
			format: (value: string[]) => value.join(','),
		},
		{
			key: 'tags',
			label: 'Tags',
			format: (value: string[]) => value.join(','),
		},
		{
			key: 'actions',
		},
	]

	onMounted(async () => {
		const res = await useApi().getSnapshots(useRoute().params.id as string)
		snapshots.value = res || []
		paths.value = _.uniq(snapshots.value.map((snapshot: Snapshot) => snapshot.paths).flat())
		loading.value = false
	})
</script>
