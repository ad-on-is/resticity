<template>
	<div>
		<h1 class="text-yellow-500 font-bold mb-3"><UIcon name="i-heroicons-clock" class="mr-2" />Schedules</h1>

		<UTable :rows="useSettings().settings?.schedules" :columns="columns" class="bg-gray-950 rounded-xl bg-opacity-50 shadow-lg" @select="">
			<template #id-data="{ row }">
				<div class="inline-flex items-center gap-1">
					<UTooltip :text="row.id"
						><span class="text-gray-400">{{ row.id.split('-')[0] }}...</span></UTooltip
					>
				</div>
			</template>
			<template #task-data="{ row }">
				<div class="inline-flex items-center gap-1">
					<span v-if="row.action === 'backup'"
						>Backup folders from

						<span class="text-primary">{{ useSettings().settings?.backups.find((b: Backup) => b.id === row.backup_id)?.name || '' }}</span>
					</span>
					<span v-if="row.action === 'copy-snapshots'"
						>Copy snapshots from

						<span class="text-purple-500"> {{ useSettings().settings?.repositories.find((r: Repository) => r?.id === row.from_repository_id)?.name || '' }}</span>
					</span>
					<span v-if="row.action === 'prune-repository'"
						>Prune

						<span class="text-purple-500"> {{ useSettings().settings?.repositories.find((r: Repository) => r?.id === row.from_repository_id)?.name || '' }}</span></span
					>
					<UIcon name="i-heroicons-chevron-double-right" />
					<span class="text-purple-500">
						<span>{{ useSettings().settings?.repositories.find((r: Repository) => r?.id === row.to_repository_id)?.name || '' }}</span></span
					>
				</div>
				<div v-if="useJobs().scheduleProgress(row.id) !== null">
					<div v-if="useJobs().scheduleProgress(row.id).percent_done">
						<UProgress :value="useJobs().scheduleProgress(row.id).percent_done * 100" class="mt-2" color="sky" />
						<div class="text-xs opacity-50 flex justify-between mt-2">
							<span>{{ useJobs().scheduleProgress(row.id).files_done }}/{{ useJobs().scheduleProgress(row.id).total_files }} files</span>
							<span>{{ humanFileSize(useJobs().scheduleProgress(row.id).bytes_done) }}/{{ humanFileSize(useJobs().scheduleProgress(row.id).total_bytes) }}</span>
							<span>{{ useJobs().scheduleProgress(row.id).seconds_remaining || 'unknown' }} seconds remaining</span>
						</div>
					</div>
					<div v-else>
						<UProgress animation="carousel" />
						<div class="text-xs opacity-50 flex justify-between mt-2">In progress</div>
					</div>
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
			<template #lastrun-data="{ row }">
				<UTooltip :text="row.last_error">
					<span v-if="row.last_run !== null && row.last_run !== ''">
						<UIcon v-if="row.last_error === ''" name="i-heroicons-check-circle" class="text-green-500 mr-1" />
						<UIcon v-else name="i-heroicons-x-circle" class="text-red-500 mr-1" />

						<span>{{ new Date(row.last_run).toUTCString() }}</span></span
					>
					<span v-else>Never</span>
				</UTooltip>
			</template>
			<template #cron-data="{ row }">
				<UBadge color="gray" v-if="row.cron !== ''">{{ cronToHuman(row.cron) }}</UBadge>
				<UBadge color="indigo" variant="outline" v-else>Manually</UBadge>
			</template>
			<template #actions-data="{ row }">
				<UToggle v-model="row.active" color="green" @update:model-value="useSettings().save()" :disabled="row.cron === ''" />
				<UDropdown :items="items(row)">
					<UButton color="gray" variant="ghost" class="ml-3" icon="i-heroicons-ellipsis-horizontal-20-solid" />
				</UDropdown>
			</template>
		</UTable>
		<UModal v-model="openDelete">
			<UCard>
				<template #header><span class="text-red-500">Delete schedule</span> </template>
				<p>Do you really want to delete this schedule?</p>
				<template #footer><UButton color="red" icon="i-heroicons-trash" @click="deleteSchedule">Yes, delete</UButton></template>
			</UCard>
		</UModal>
	</div>
</template>

<script setup lang="ts">
	import cronstrue from 'cronstrue'

	const openDelete = ref(false)
	let toDelete: any = null
	const columns = [
		{ key: 'id', class: 'w-32', label: 'ID' },
		{ key: 'status', label: 'Status', class: 'w-32' },
		{ key: 'task', label: 'Task', class: 'w-128' },
		{ key: 'lastrun', label: 'Last run', class: 'w-48' },
		{ key: 'cron', label: 'Scheduled', class: 'w-32' },
		{ key: 'actions', class: 'w-10' },
	]

	const cronToHuman = (cron: string) => {
		try {
			return cronstrue.toString(cron)
		} catch (e) {
			return cron
		}
	}

	const deleteSchedule = async () => {
		if (toDelete === null) {
			openDelete.value = false
			return
		}
		useSettings().settings!.schedules = useSettings().settings!.schedules.filter((item: any) => item.id !== toDelete.id)
		await useSettings().save()
		openDelete.value = false
		toDelete = null
	}

	const items = (row: any) => [
		[
			!useJobs().scheduleIsRunning(row.id)
				? {
						label: 'Run now',
						icon: 'i-heroicons-arrow-uturn-right',
						click: async () => {
							const t = await useApi().runSchedule(row.id)
						},
				  }
				: {
						label: 'Stop',
						icon: 'i-heroicons-arrow-uturn-right',
						click: async () => {
							const t = await useApi().stopSchedule(row.id)
						},
				  },

			{
				label: 'Delete',
				icon: 'i-heroicons-trash',
				click: () => {
					toDelete = row
					openDelete.value = true
				},
			},
		],
	]

	onUnmounted(() => {
		useSettings().refresh()
	})
</script>
