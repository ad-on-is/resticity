<template>
	<div>
		<h1 class="text-info m-0"><FaIcon icon="folder-open" class="mr-3" />{{ backup?.name }}</h1>
		<h2 class="mt-2 opacity-40 text-info">{{ backup?.path }}</h2>
		<BackupTargets @update="(val) => (targets = val)" />
		<BackupExcludeOptions @update="(val) => (excludes = val)" />
		<BackupScheduleOptions />
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import _ from 'lodash'

	const backup = ref<Backup>()
	const targets = ref<string[]>([])
	const excludes = ref<[]>([])

	const update = _.debounce(() => {
		backup.value.targets = targets.value
		backup.value.backup_params = excludes.value
		useSettings().settings!.backups[useSettings().settings!.backups.findIndex((b: Backup) => b.id === backup.value.id)] = backup.value
		useSettings().save()
		console.log(excludes.value)
	}, 300)

	watch([targets, excludes], () => {
		// console.log(targets.value)
		update()
	})

	onMounted(async () => {
		backup.value = useSettings().settings?.backups.find((b: Backup) => b.id === useRoute().params.id)
	})

	function save() {}
</script>
