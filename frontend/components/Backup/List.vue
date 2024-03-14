<template>
	<h1 class="text-sky-500 font-bold mb-3"><FaIcon icon="upload" class="mr-3" />Backups</h1>
	<div class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
		<div
			v-if="props.showNew"
			class="opacity-40 border cursor-pointer border-dashed border-sky-500 border-opacity-40 hover:opacity-100 shadow-lg rounded-lg no-underline hover:bg-sky-500 transition-all hover:bg-opacity-10"
			:class="colorClass"
			@click="isOpen = true"
		>
			<div class="p-5">
				<h3 class="m-0 text-sky-500 font-medium"><FaIcon icon="fa-plus-circle" class="mr-2" />Add Backup</h3>
				<p class="text-sm" :class="textColorClass">Add a new Backup</p>
			</div>
		</div>
		<NuxtLink
			:to="`/backups/${backup.id}`"
			v-for="backup in useSettings().settings?.backups"
			class="shadow-lg rounded-lg no-underline hover:bg-sky-500 transition-all hover:bg-opacity-10"
			:class="colorClass"
		>
			<div class="p-5 pb-0">
				<h3 class="m-0 text-sky-500 font-medium"><FaIcon icon="folder-open" class="mr-2" />{{ backup.name }}</h3>
				<p class="text-xs break-words" :class="textColorClass">{{ backup.path }}</p>
				<div :class="useJobs().backupIsRunning(backup.id) ? 'opacity-100' : 'opacity-0'">
					<span class="loading loading-infinity loading-sm text-orange-500"></span>
				</div>
			</div>
		</NuxtLink>
	</div>

	<UModal v-model="isOpen">
		<UCard>
			<template #header>
				<h2 class="text-sky-500 font-bold">New Backup</h2>
				<p v-if="error === ''" class="text-sm opacity-40">Add a new backup location</p>
				<p v-else class="text-sm text-error"><FaIcon icon="warning" class="mr-2" />{{ error }}</p>
			</template>
			<template #footer>
				<UButton @click="save" icon="i-heroicons-plus-circle" :disabled="newBackup.path === '' || newBackup.name === ''">Add Backup</UButton>
			</template>
			<UInput v-model="newBackup.name" placeholder="Name" class="mb-5" />
			<PathAutocomplete @selected="(p) => (newBackup.path = p)" />
		</UCard>
	</UModal>
</template>

<script setup lang="ts">
	import { generateUUID } from '~/utils'
	const isOpen = ref(false)
	const newBackup = ref({
		id: generateUUID(),
		name: '',
		path: '',
		cron: '',
		backup_params: [],
		targets: [],
	})
	const error = ref('')

	const props = defineProps({
		showNew: {
			type: Boolean,
			default: true,
		},
	})

	const save = async () => {
		useSettings().settings!.backups.push(newBackup.value)
		await useSettings().save()
		isOpen.value = false
	}

	const colorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'bg-base-300' : 'bg-base-300 bg-opacity-10'
	})

	const textColorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'opacity-50' : 'text-black'
	})
</script>
