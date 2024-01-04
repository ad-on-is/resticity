<template>
	<div v-if="backup">
		<h1 class="text-info m-0"><FaIcon icon="folder-open" class="mr-3" />{{ backup?.name }}</h1>
		<h2 class="mt-2 opacity-40 text-info">{{ backup?.path }} {{ useRoute().params.id }}</h2>
		<BackupTargets @update="(val) => (targets = val)" :targets="targets" />
		<BackupExcludeOptions @update="(val) => (excludes = val)" :excludes="excludes" />
		<BackupScheduleOptions @update="(val) => (cron = val)" :cron="cron" />
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import _ from 'lodash'

	const backup = ref<Backup>()
	const targets = ref<string[]>([])
	const excludes = ref<[]>([])
	const cron = ref<string>('')
	const init = ref(true)
	const idx = ref(-1)

	const update = _.debounce(() => {
		console.log('SHOULD UPDATE')
		if (init.value) {
			init.value = false
			return
		}
		backup.value.targets = targets.value
		backup.value.backup_params = excludes.value
		backup.value.cron = cron.value
		useSettings().settings!.backups[idx.value] = backup.value
		useSettings().save()
	}, 300)

	watch(
		() => [JSON.stringify(targets.value), JSON.stringify(excludes.value), JSON.stringify(cron.value)],
		() => {
			update()
		}
	)

	onMounted(async () => {
		backup.value = useSettings().settings!.backups.find((b: Backup) => b.id === useRoute().params.id)
		idx.value = useSettings().settings!.backups.findIndex((b: Backup) => b.id === backup.value.id)
		targets.value = backup.value.targets
		excludes.value = backup.value.backup_params
		cron.value = backup.value.cron

		console.log(useSettings().settings!.backups)
	})
</script>
