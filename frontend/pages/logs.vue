<template>
	<div>
		<UTabs :items="items">
			<template #schedulelogs>
				<UAccordion :items="scheduleLogs"> </UAccordion>
			</template>
			<template #scheduleerrors>
				<UAccordion :items="scheduleErrors"> </UAccordion>
			</template>
			<template #archivelogs>
				<UAccordion :items="fileLogs"> </UAccordion>
			</template>
			<template #archiveerrors>
				<UAccordion :items="fileErrors"> </UAccordion>
			</template>
		</UTabs>
	</div>
</template>

<script lang="ts" setup>
	const fileLogs = ref([])
	const fileErrors = ref([])
	const getLabelByScheduleId = (id: string) => {
		let label = ''
		const schedule = useSettings().settings.schedules.find((s: Schedule) => s.id === id)
		if (!schedule) return id
		label = schedule.action
		if (schedule.from_repository_id !== '') {
			label += ` from ${useSettings().settings.repositories.find((r: Repository) => r.id === schedule.from_repository_id).name} to`
		}
		if (schedule.backup_id !== '') {
			label += ` from ${useSettings().settings.backups.find((r: Backup) => r.id === schedule.backup_id).name} to`
		}
		if (schedule.to_repository_id !== '') {
			label += ` ${useSettings().settings.repositories.find((r: Repository) => r.id === schedule.to_repository_id).name}`
		}

		// return id
		return label
	}

	const items = [
		{ slot: 'schedulelogs', label: 'Schedule: Logs' },
		{ slot: 'scheduleerrors', label: 'Schedule: Errors' },
		{ slot: 'archivelogs', label: 'Archive: Logs' },
		{ slot: 'archiveerrors', label: 'Archive: Errors' },
	]

	const scheduleLogs = Object.keys(useLogs().out).map((l) => {
		return { label: `Schedule: ${getLabelByScheduleId(l)}`, content: useLogs().out[l], color: 'blue' }
	})
	const scheduleErrors = Object.keys(useLogs().out).map((l) => {
		return { label: `Schedule: ${getLabelByScheduleId(l)}`, content: useLogs().err[l], color: 'red' }
	})

	onMounted(async () => {
		const { logs, errors } = await useApi().getLogs()
		fileLogs.value = logs.map((l: string) => {
			return { slot: 'content', file: l, label: `File: ${l}`, color: 'blue' }
		})
		fileErrors.value = errors.map((l: string) => {
			return { slot: 'content', file: l, label: `File: ${l}`, color: 'red' }
		})
	})
</script>
