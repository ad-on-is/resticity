<template>
	<div>
		<h1 class="text-yellow-500 font-bold mb-3"><UIcon name="i-heroicons-clock" class="mr-2" dynamic />Schedules</h1>

		<UTable :rows="useSettings().settings?.schedules" :columns="columns" class="bg-gray-950 rounded-xl bg-opacity-50 shadow-lg" @select="">
			<template #task-data="{ row }">
				<div class="inline-flex items-center gap-1">
					<span v-if="row.backup_id !== ''"
						>Backup folders from

						<span class="text-primary">{{ useSettings().settings?.backups.find((b: Backup) => b.id === row.backup_id)?.name || '' }}</span>
					</span>
					<span v-if="row.from_repository_id !== ''"
						>Copy snapshots from

						<span class="text-purple-500"> {{ useSettings().settings?.repositories.find((r: Repository) => r?.id === row.from_repository_id)?.name || '' }}</span></span
					>

					<UIcon name="i-heroicons-chevron-double-right" />
					<span class="text-purple-500">
						<span>{{ useSettings().settings?.repositories.find((r: Repository) => r?.id === row.to_repository_id)?.name || '' }}</span></span
					>
				</div>
			</template>
			<template #status-data="{ row }">
				<div v-if="row.active">
					<UBadge v-if="useJobs().scheduleIsRunning(row.id)" color="green">Running</UBadge>
					<UBadge v-else color="primary" variant="outline">Scheduled</UBadge>
				</div>
				<div v-else>
					<UBadge v-if="useJobs().scheduleIsRunning(row.id)" color="green">Running</UBadge>
					<UBadge v-else color="gray" class="opacity-40">Inactive</UBadge>
				</div>
			</template>
			<template #cron-data="{ row }">
				<UBadge color="gray" v-if="row.cron !== ''">{{ cronToHuman(row.cron) }}</UBadge>
				<UBadge color="indigo" variant="outline" v-else>Manually</UBadge>
			</template>
			<template #actions-data="{ row }">
				<UToggle v-model="row.active" color="green" />
				<UDropdown :items="items(row)">
					<UButton color="gray" variant="ghost" class="ml-3" icon="i-heroicons-ellipsis-horizontal-20-solid" />
				</UDropdown>
			</template>
		</UTable>
	</div>
</template>

<script setup lang="ts">
	import cronstrue from 'cronstrue'
	const columns = [
		{ key: 'status', label: 'Status', class: 'w-32' },
		{ key: 'task', label: 'Task' },
		{ key: 'cron', label: 'Scheduled' },
		{ key: 'actions', class: 'w-10' },
	]

	const cronToHuman = (cron: string) => {
		try {
			return cronstrue.toString(cron)
		} catch (e) {
			return cron
		}
	}

	const items = (row: any) => [
		[
			{
				label: 'Run now',
				icon: 'i-heroicons-arrow-uturn-right',
				click: async () => {
					const t = await useApi().runSchedule(row.id)
					console.log(t)
				},
			},
			{
				label: 'Delete',
				icon: 'i-heroicons-trash',
				click: () => {
					useSettings().settings!.schedules = useSettings().settings!.schedules.filter((item: any) => item.id !== row.id)
					useSettings().save()
				},
			},
		],
	]
</script>
