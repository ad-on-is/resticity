<template>
	<h3 class="text-purple-500 mb-3"><UIcon name="i-heroicons-queue-list" />Snapshots</h3>
	<UButtonGroup>
		<UButton color="gray" disabled>Group by</UButton>
		<USelectMenu v-model="selectedGroupBy" class="w-48" :options="groupByOptions" color="purple">
			<template #label> <UIcon :name="selectedGroupBy.icon" /> {{ selectedGroupBy.label }} </template>
		</USelectMenu>
		<UButton v-if="selected.length > 0" icon="i-heroicons-trash" color="orange" variant="outline">Prune snapshots</UButton>
	</UButtonGroup>
	<div v-if="loading" class="mt-10">LOADING Snapshots</div>
	<div v-else class="mt-10">
		<UAccordion :items="snapshotGroups" color="gray" variant="outline">
			<template #default="{ item, index, open }">
				<UButton color="gray" variant="ghost" :class="open ? 'bg-gray-950/50' : ''" :ui="{ padding: { sm: 'p-3' } }">
					<template #leading>
						<UIcon :name="item.icon" class="w-4 h-4 text-teal-500" />
					</template>

					<span class="truncate"
						>{{ item.label }} <span class="opacity-40">({{ item.snapshots.length }} snapshots)</span></span
					>

					<template #trailing>
						<UIcon name="i-heroicons-chevron-right-20-solid" class="w-5 h-5 ms-auto transform transition-transform duration-200" :class="[open && 'rotate-90']" />
					</template>
				</UButton>
			</template>
			<template #item="{ item }">
				<UTable :rows="_.orderBy(item.snapshots, 'time', 'desc')" v-model="selected" :columns="columns" @select="" class="bg-gray-950 rounded-xl bg-opacity-50 shadow-lg">
					<template #time-data="{ row }">
						<span class="text-teal-600">{{ formatISO9075(new Date(row.time)) }}</span>
					</template>
					<template #hostname-data="{ row }"
						><span class="text-pink-500"><UIcon name="i-heroicons-tv" /> {{ row.hostname }}</span></template
					>
					<template #info-data="{ row }">
						<div class="mb-2">
							<UButton
								@click="navigateTo({ path: `/repositories/${useRoute().params.id}/snapshots/${row.id}`, query: { path: path } })"
								v-for="path in row.paths"
								size="xs"
								color="yellow"
								variant="link"
								icon="i-heroicons-folder"
								>{{ path }}</UButton
							>
						</div>
						<div class="gap-2 flex">
							<UBadge v-for="tag in row.tags" :variant="tag === 'resticity' ? 'outline' : 'solid'" :color="tag === 'resticity' ? 'sky' : 'gray'" size="xs"
								><UIcon name="i-heroicons-tag-solid" class="mr-1" />{{ tag }}</UBadge
							>
						</div>
					</template>
				</UTable>
			</template>
		</UAccordion>
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import { formatISO9075 } from 'date-fns'
	import _ from 'lodash'
	const groupByOptions = [
		{ id: 'host', label: 'Host', icon: 'i-heroicons-tv' },
		{ id: 'path', label: 'Path', icon: 'i-heroicons-folder' },
		{ id: 'tag', label: 'Tag', icon: 'i-heroicons-tag' },
	]
	const snapshotGroups = ref<Array<SnapshotGroup>>([])
	const loading = ref(true)
	const selected = ref<Array<Snapshot>>([])
	const selectedGroupBy = ref(groupByOptions[0])
	const mounted = ref<Array<string>>([])
	const groupBy = ref('')

	const paths = ref<Array<string>>([])

	const items = (row: any) => [
		[
			...row.paths.map((p) => ({
				label: 'Browse ' + p,
				icon: 'i-heroicons-document-magnifying-glass',
				click: () => {
					navigateTo({ path: `/repositories/${useRoute().params.id}/snapshots/${row.id}`, query: { path: p } })
				},
			})),
			{
				label: 'Quick Restore',
				icon: 'i-heroicons-arrow-top-right-on-square',
				click: () => {},
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
			class: 'w-24',
		},
		{
			key: 'time',
			label: 'Time',
			class: 'w-32',
			format: (value: string) => format(new Date(value), 'dd.MM.yyyy H:I:s'),
		},
		{
			key: 'hostname',
			label: 'Host',
		},

		{
			key: 'info',
			label: 'Info',
			format: (value: string[]) => value.join(','),
		},
	]

	function getLabel(gk: GroupKey) {
		let label = ''
		let icon = ''
		if (gk.hostname && gk.hostname !== '') {
			label = gk.hostname
			icon = 'i-heroicons-tv'
		}

		if (gk.paths) {
			label = gk.paths.join(',')
			icon = 'i-heroicons-folder'
		}

		if (selectedGroupBy.value.id === 'tag') {
			label = gk.tags && gk.tags.length > 0 ? gk.tags.join(',') : '(Untagged)'
			icon = 'i-heroicons-tag'
		}

		return { label, icon }
	}

	async function load() {
		loading.value = true
		const res = await useApi().getSnapshots(useRoute().params.id as string, selectedGroupBy.value.id)
		console.log(res)
		snapshotGroups.value = res.map((r: SnapshotGroup) => ({ ...r, ...getLabel(r.group_key) })) || []
		paths.value = _.uniq(snapshotGroups.value.map((snapshot: Snapshot) => snapshot.paths).flat())
		loading.value = false
	}

	watch(selectedGroupBy, async () => {
		load()
	})

	onMounted(async () => {
		load()
	})
</script>
