<template>
	<div class="bg-gray-950 p-5 rounded-lg mt-10 bg-opacity-50 shadow-lg">
		<h1 class="text-yellow-700 font-bold">Add a new task to schedule</h1>
		<p class="opacity-40 mb-3">Create a new schedule to either run manually or in the background at a specific time.</p>
		<UButtonGroup class="">
			<USelectMenu v-model="selectedAction" :options="actionOptions" class="w-44">
				<template #label>
					<UIcon :name="selectedAction.icon" />
					<span>{{ selectedAction.label }}</span>
				</template>
			</USelectMenu>
			<USelectMenu v-if="selectedAction.id === 'backup'" v-model="selectedBackup" :options="backups()" option-attribute="name" class="w-48">
				<template #label>
					<UIcon :name="selectedBackup ? 'i-heroicons-folder' : 'i-heroicons-ellipsis-horizontal-circle'" />
					<span>{{ selectedBackup.name }}</span>
				</template>
			</USelectMenu>
			<USelectMenu v-if="selectedAction.id === 'copy-snapshots'" v-model="selectedFromRepository" :options="repositories('From Repository')" option-attribute="name" class="w-48">
				<template #label>
					<UIcon name="i-heroicons-server" />
					<span>{{ selectedFromRepository.name }}</span>
				</template>
			</USelectMenu>
			<USelectMenu
				v-if="selectedBackup.id !== '' || selectedFromRepository.id !== '' || selectedAction.id === 'prune-repository'"
				v-model="selectedToRepository"
				:options="repositories('Repository', '')"
				option-attribute="name"
				class="w-48"
			>
				<template #label>
					<UIcon name="i-heroicons-server" />
					<span>{{ selectedToRepository.name }}</span>
				</template>
			</USelectMenu>
			<USelectMenu v-model="selectedCron" :options="cronOptions" class="w-48"></USelectMenu>
			<UInput class="w-32" v-model="cron" placeholder="" />
			<UButton @click="addSchedule" color="yellow" icon="i-heroicons-plus-circle">Add Schedule</UButton>
		</UButtonGroup>
	</div>
</template>

<script setup lang="ts">
	const backups = (title: string = 'From Backup') => [
		{ id: '', name: title, disabled: true },
		...useSettings().settings!.backups.map((o: any) => ({ name: 'from ' + o.name, id: o.id, icon: 'i-heroicons-folder' })),
	]
	const repositories = (title: string = 'Select Destination', prefix: string = 'from ') => [
		{ id: '', name: title, disabled: true },
		...useSettings().settings!.repositories.map((o: any) => ({ name: prefix + o.name, id: o.id, icon: 'i-heroicons-server' })),
	]

	const actionOptions = [
		{ id: '', label: 'Select Action', disabled: true, icon: 'i-heroicons-ellipsis-horizontal-circle' },
		{ id: 'backup', label: 'Run Backup', icon: 'i-heroicons-arrow-up-tray' },
		{ id: 'copy-snapshots', label: 'Copy Snapshots', icon: 'i-heroicons-server' },
		{ id: 'prune-repository', label: 'Prune repository', icon: 'i-heroicons-server' },
	]
	const cronOptions = [
		{ label: 'Run manually', disabled: true },
		{ label: 'Every 5 Minutes', value: '*/5 * * * *' },
		{ label: 'Every 10 Minutes', value: '*/10 * * * *' },
		{ label: 'Every 15 Minutes', value: '*/15 * * * *' },
		{ label: 'Every 30 Minutes', value: '*/30 * * * *' },
		{ label: 'Every Hour', value: '0 * * * *' },
		{ label: 'Every Day', value: '0 0 * * *' },
		{ label: 'Every Week', value: '0 0 * * 0' },
		{ label: 'Every Month', value: '0 0 1 * *' },
		{ label: 'Every Year', value: '0 0 1 1 *' },
		{ label: 'Custom schedule', value: '* * * * *' },
	]

	const selectedAction = ref(actionOptions[0])
	const selectedCron = ref(cronOptions[0])
	const selectedBackup = ref(backups()[0])
	const selectedFromRepository = ref(repositories('From Repository')[0])
	const selectedToRepository = ref(repositories('To Repository')[0])

	const cron = ref('')

	watch(selectedAction, () => {
		selectedBackup.value = backups()[0]
		selectedFromRepository.value = repositories('From Repository')[0]
		selectedToRepository.value = repositories('To Repository')[0]
		selectedCron.value = cronOptions[0]
	})

	watch(selectedCron, () => {
		cron.value = selectedCron.value.value || ''
	})
	const addSchedule = () => {
		useSettings().settings!.schedules.push({
			id: generateUUID(),
			action: selectedAction.value.id,
			backup_id: selectedBackup.value.id,
			to_repository_id: selectedToRepository.value.id,
			from_repository_id: selectedFromRepository.value.id,
			cron: cron.value,
			active: false,
		})
		selectedAction.value = actionOptions[0]
		useSettings().save()
	}
</script>
