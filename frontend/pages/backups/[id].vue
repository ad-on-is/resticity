<template>
	<div v-if="backup">
		<div class="flex justify-between">
			<div>
				<h1 class="text-sky-500 font-bold m-0"><UIcon name="i-heroicons-folder-open" class="mr-3" />{{ backup?.name }}</h1>
				<h2 class="">{{ backup?.path }}</h2>
			</div>
			<div class="mt-3">
				<UButton icon="i-heroicons-trash" color="red" class="ml-2" @click="openDelete = true">Delete</UButton>
			</div>
		</div>

		<UDivider class="my-5" />
		<BackupExcludeOptions @update="(val) => (excludes = val)" :excludes="excludes" />
		<UModal v-model="openDelete">
			<UCard>
				<template #header><span class="text-red-500">Delete backup</span> </template>
				<p>Do you really want to delete this backup and its associated schedules? <br /><span class="opacity-50 text-sm">This will not delete actual files.</span></p>
				<template #footer><UButton color="red" icon="i-heroicons-trash" @click="deleteBackup">Yes, delete</UButton></template>
			</UCard>
		</UModal>
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import _ from 'lodash'
	const openDelete = ref(false)
	const backup = ref<Backup>()
	const excludes = ref<[]>([])
	const idx = ref(-1)

	const deleteBackup = async () => {
		useSettings().settings!.backups = useSettings().settings!.backups.filter((item: Backup) => item.id !== backup.value.id)
		useSettings().settings!.schedules = useSettings().settings!.schedules.filter((item: Schedule) => item.backup_id !== backup.value.id)
		await useSettings().save()
		openDelete.value = false
		return navigateTo('/backups')
	}

	const update = _.debounce(() => {
		backup.value.backup_params = excludes.value
		useSettings().settings!.backups[idx.value] = backup.value
		useSettings().save()
	}, 300)

	onMounted(async () => {
		backup.value = useSettings().settings!.backups.find((b: Backup) => b.id === useRoute().params.id)
		idx.value = useSettings().settings!.backups.findIndex((b: Backup) => b.id === backup.value.id)
		excludes.value = backup.value.backup_params
		watch(
			() => [JSON.stringify(excludes.value)],
			() => {
				update()
			}
		)
	})
</script>
