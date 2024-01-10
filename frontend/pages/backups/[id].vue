<template>
	<div v-if="backup">
		<h1 class="text-sky-500 font-bold m-0"><UIcon name="i-heroicons-folder-open" class="mr-3" />{{ backup?.name }}</h1>
		<h2 class="text-sky-500">{{ backup?.path }}</h2>
		<UDivider class="my-5" />
		<BackupExcludeOptions @update="(val) => (excludes = val)" :excludes="excludes" />
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import _ from 'lodash'

	const backup = ref<Backup>()
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
		backup.value.backup_params = excludes.value
		backup.value.cron = cron.value
		useSettings().settings!.backups[idx.value] = backup.value
		useSettings().save()
	}, 300)

	watch(
		() => [JSON.stringify(excludes.value), JSON.stringify(cron.value)],
		() => {
			update()
		}
	)

	onMounted(async () => {
		backup.value = useSettings().settings!.backups.find((b: Backup) => b.id === useRoute().params.id)
		idx.value = useSettings().settings!.backups.findIndex((b: Backup) => b.id === backup.value.id)
		excludes.value = backup.value.backup_params
		cron.value = backup.value.cron

		console.log(useSettings().settings!.backups)
	})
</script>
