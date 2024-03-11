<template>
	<div>
		<UTabs :items="items">
			<template #logs>
				<UAccordion :items="logItems"> </UAccordion>
			</template>
			<template #errors>
				<UAccordion :items="errorItems"> </UAccordion>
			</template>
		</UTabs>
	</div>
</template>

<script lang="ts" setup>
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
		{ slot: 'logs', label: 'Logs' },
		{ slot: 'errors', label: 'Errors' },
	]

	const logItems = Object.keys(useLogs().out).map((l) => {
		return { label: getLabelByScheduleId(l), content: useLogs().out[l], color: 'green' }
	})
	const errorItems = Object.keys(useLogs().out).map((l) => {
		return { label: getLabelByScheduleId(l), content: useLogs().err[l], color: 'red' }
	})
</script>
