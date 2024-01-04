<template>
	<h1 class="text-primary"><FaIcon icon="upload" class="mr-3" />Backups</h1>
	<div class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
		<NuxtLink
			to="#"
			class="opacity-40 border border-dashed border-success border-opacity-40 hover:opacity-100 shadow-lg bg-base-300 rounded-lg no-underline hover:bg-success transition-all hover:bg-opacity-10"
			@click="dialog?.showModal()"
		>
			<div class="p-5">
				<h3 class="m-0 text-success"><FaIcon icon="fa-plus-circle" class="mr-2" />Add Backup</h3>
				<p class="text-sm opacity-40">Add a new Backup</p>
			</div>
		</NuxtLink>
		<NuxtLink
			:to="`/backups/${backup.id}`"
			v-for="backup in useSettings().settings?.backups"
			class="shadow-lg bg-base-300 rounded-lg no-underline hover:bg-primary transition-all hover:bg-opacity-10"
		>
			<div class="p-5">
				<h3 class="m-0 text-info"><FaIcon icon="folder-open" class="mr-2" />{{ backup.name }}</h3>
				<p class="text-sm">{{ backup.path }}</p>
				<div class="flex animate-pulse" v-if="useJobs().backupIsRunning(backup.id)">
					<span class="loading loading-infinity loading-sm text-warning"></span><span class="text-sm ml-2 text-warning">Backup running</span>
				</div>
			</div>
		</NuxtLink>
	</div>

	<dialog class="modal" ref="dialog">
		<div class="modal-box">
			<h2 class="m-0 mb-3">New Backup</h2>
			<p v-if="error === ''" class="text-sm opacity-40">Add a new Backup</p>
			<p v-else class="text-sm text-error"><FaIcon icon="warning" class="mr-2" />{{ error }}</p>
			<div class="">
				<input class="input input-bordered input-sm mb-5 min-w-full" placeholder="Name" v-model="newBackup.name" />

				<div class="join w-full">
					<input class="input input-bordered join-item input-sm flex-grow" disabled placeholder="Location" v-model="newBackup.path" />
					<button class="btn join-item btn-sm input-bordered" @click="openDir()"><FaIcon icon="folder" /></button>
				</div>
				<div class="modal-action">
					<button class="btn btn-sm btn-success" @click="save"><FaIcon icon="plus-circle" />Add Backup</button>
				</div>
			</div>
		</div>
	</dialog>
</template>

<script setup lang="ts">
	import { generateUUID } from '~/utils'

	const dialog = ref<HTMLDialogElement>()
	const newBackup = ref({
		id: generateUUID(),
		name: '',
		path: '',
		cron: '',
		backup_params: [],
		targets: [],
	})
	const error = ref('')

	const openDir = async () => {
		const dir = await SelectDirectory('Select a folder to backup')
		if (useSettings().settings?.backups.find((b: any) => b.path === dir)) {
			error.value = `${dir} is already a backup`
			return
		}
		if (dir !== '') {
			error.value = ''
			newBackup.value.path = dir
		}
	}

	const save = async () => {
		useSettings().settings?.backups.push(newBackup.value)
		await useSettings().save()
		dialog.value?.close()
	}
</script>
