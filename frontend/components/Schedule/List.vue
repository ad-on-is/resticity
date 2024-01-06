import type settingsVue from '~/pages/settings.vue';
<template>
	<div>
		<h1>Schedules</h1>

		<UTable :rows="useSettings().settings?.schedules" :columns="columns" class="mt-10">
			<template #task-data="{ row }">
				<div class="inline-flex items-center gap-1">
					<UBadge v-if="useJobs().scheduleIsRunning(row.id)" color="orange">Running</UBadge>
					<span v-if="row.backup_id !== ''">Backup from</span>
					<span v-if="row.backup_id !== ''" class="">
						<span class="text-primary"
							><UIcon name="i-heroicons-folder" class="h-2.5" />{{ useSettings().settings?.backups.find((b: Backup) => b.id === row.backup_id)?.name || '' }}</span
						>
					</span>
					<span v-if="row.from_repository_id !== ''">Copy snapshots from</span>
					<span v-if="row.from_repository_id !== ''">
						<span class="text-primary"
							><UIcon name="i-heroicons-server" class="h-2.5 text-primary" />
							{{ useSettings().settings?.repositories.find((r: Repository) => r?.id === row.from_repository_id)?.name || '' }}</span
						>
					</span>
					<span>to</span>
					<span class="text-primary">
						<UIcon name="i-heroicons-server" class="h-2.5 text-primary" /><span>{{
							useSettings().settings?.repositories.find((r: Repository) => r?.id === row.to_repository_id)?.name || ''
						}}</span></span
					>
				</div>
			</template>
			<template #cron-data="{ row }">
				<code class="text-warning bg-gray-950 p-1">{{ row.cron }}</code>
			</template>
			<template #actions-data="{ row }">
				<UDropdown :items="items(row)">
					<UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
				</UDropdown>
			</template>
		</UTable>
		<ScheduleNew />
	</div>
</template>

<script setup lang="ts">
	const columns = [{ key: 'task', label: 'Task' }, { key: 'cron', label: 'Cron' }, { key: 'actions' }]

	const items = (row: any) => [
		[
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
