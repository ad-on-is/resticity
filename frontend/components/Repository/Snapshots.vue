<template>
	<div class="collapse bg-base-200 mb-5">
		<input type="radio" name="repository-accordion" checked />
		<h3 class="collapse-title m-0 text-primary"><FaIcon icon="table-list" class="mr-2" />Snapshots</h3>
		<div class="collapse-content">
			<UButton :color="selected.length > 0 ? 'primary' : 'primary'" :disabled="selected.length === 0" :variant="selected.length === 0 ? 'outline' : 'solid'"
				>Prune selected Snapshots</UButton
			>
			<UTable :rows="snapshots" v-model="selected" :columns="columns" @select="" :loading="loading">
				<template #tags-data="{ row }">
					<UBadge v-for="tag in row.tags" variant="outline" color="blue">{{ tag }}</UBadge>
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
		</div>
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import { format } from 'date-fns'
	const snapshots = ref<Array<Snapshot>>([])
	const loading = ref(true)
	const selected = ref<Array<Snapshot>>([])
	const mounted = ref<Array<string>>([])

	const items = (row: any) => [
		[
			{
				label: mounted.value.includes(row.id) ? 'Unmount' : 'Mount',
				icon: mounted.value.includes(row.id) ? 'i-iconify-folder' : 'i-heroicons-folder-open',
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
		[
			{
				label: 'Prune',
				icon: 'i-heroicons-trash',
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
		const res = await useFetch(`http://localhost:11278/api/snapshots/${useRoute().params.id}`, { cache: 'no-cache' })
		console.log(res.data.value)
		snapshots.value = res.data.value || []
		loading.value = false
	})
</script>
