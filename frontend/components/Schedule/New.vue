<template>
	<div class="bg-gray-950 p-5 rounded-md mt-10 bg-opacity-50">
		<h1 class="text-primary">Add a new task to schedule</h1>
		<p class="opacity-40 mb-3">Create a new schedule to either run manually or in the background at a specific time.</p>
		<UButtonGroup class="">
			<USelectMenu v-model="selectedType" :options="typeOptions" class="w-44">
				<template #label>
					<UIcon :name="selectedType.icon || 'i-heroicons-ellipsis-horizontal-circle'" />
					<span>{{ selectedType.label }}</span>
				</template>
			</USelectMenu>
			<USelectMenu v-model="selectedFrom" :options="fromOptions" option-attribute="name" class="w-48">
				<template #label>
					<UIcon :name="selectedFrom.icon || 'i-heroicons-ellipsis-horizontal-circle'" />
					<span>{{ selectedFrom.name }}</span>
				</template>
			</USelectMenu>
			<USelectMenu v-model="selectedTo" :options="toOptions" option-attribute="name" class="w-48">
				<template #label>
					<UIcon :name="'i-heroicons-server'" />
					<span>{{ selectedTo.name }}</span>
				</template>
			</USelectMenu>
			<USelectMenu v-model="selectedCron" :options="cronOptions" class="w-48"></USelectMenu>
			<UInput class="w-48" v-model="cron" placeholder="" />
			<UButton @click="addSchedule" icon="i-heroicons-plus-circle">Add Schedule</UButton>
		</UButtonGroup>
	</div>
</template>

<script setup lang="ts">
	type Option = { name: string; id: string; icon?: string; disabled?: boolean }
	const defaultFrom = { id: '', name: 'Select Location', disabled: true }
	const fromOptions = ref<Array<Option>>([defaultFrom])
	const toOptions: Array<Option> = [
		...[{ id: '', name: 'Select Destination', disabled: true }],
		...useSettings().settings!.repositories.map((o: Repository) => ({ name: 'to ' + o.name, id: o.id, icon: 'i-heroicons-server' })),
	]

	const typeOptions = [
		{ id: '', label: 'Select Action', disabled: true },
		{ id: 'backup', label: 'Run Backup', icon: 'i-heroicons-arrow-up-tray' },
		{ id: 'repository', label: 'Copy Snapshots', icon: 'i-heroicons-server' },
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
	const selectedType = ref(typeOptions[0])
	const selectedFrom = ref(fromOptions.value[0])
	const selectedTo = ref(toOptions[0])
	const selectedCron = ref(cronOptions[0])
	const cron = ref('')
	watch(selectedType, () => {
		if (selectedType.value.id === 'backup') {
			fromOptions.value = [defaultFrom, ...useSettings().settings!.backups.map((o: Backup) => ({ name: 'from ' + o.name, id: o.id, icon: 'i-heroicons-folder' }))]
		} else if (selectedType.value.id === 'repository') {
			fromOptions.value = [defaultFrom, ...useSettings().settings!.repositories.map((o: Repository) => ({ name: 'from ' + o.name, id: o.id, icon: 'i-heroicons-server' }))]
		} else {
			fromOptions.value = [defaultFrom]
		}
		selectedFrom.value = fromOptions.value[0]
	})

	watch(selectedFrom, () => {
		selectedTo.value = toOptions[0]
	})

	watch(selectedCron, () => {
		cron.value = selectedCron.value.value || ''
	})
	const addSchedule = () => {
		useSettings().settings!.schedules.push({
			id: generateUUID(),
			backup_id: selectedType.value.id === 'backup' ? selectedFrom.value.id : '',
			to_repository_id: selectedTo.value.id,
			from_repository_id: selectedType.value.id === 'repository' ? selectedFrom.value.id : '',
			cron: cron.value,
		})

		selectedType.value = typeOptions[0]
		selectedCron.value = cronOptions[0]

		cron.value = ''
		useSettings().save()
	}
</script>
