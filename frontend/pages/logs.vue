<template>
	<div>
		<div class="flex justify-between">
			<div>
				<h1 class="text-teal-500 font-bold"><UIcon name="i-heroicons-cog-6-tooth" class="mr-2" />Logs</h1>
			</div>
		</div>

		<div class="grid grid-cols-2 gap-10 p-10 bg-opacity-70 rounded-lg shadow-lg mt-5" :class="colorClass">
			<div>
				<h4 class="text-teal-500">Schedule logs</h4>
				<UAccordion v-if="scheduleLogs.length > 0" :items="scheduleLogs">
					<template #content="{ item }">
						<pre class="text-xs overflow-scroll h-72 bordered">{{ item.content }}</pre>
					</template>
				</UAccordion>
				<p v-else>No logs</p>
				<h4 class="text-teal-500 mt-5">Schedule errors</h4>
				<UAccordion v-if="scheduleErrors.length > 0" :items="scheduleErrors">
					<template #content="{ item }">
						<pre class="text-xs overflow-scroll h-72">{{ item.content }}</pre>
					</template>
				</UAccordion>
				<p v-else>No errors</p>
				<h4 class="text-teal-500 mt-5">Server errors</h4>
				<pre v-if="useLogs().serverErr.length > 0" class="text-xs overflow-scroll h-72">{{ useLogs().serverErr.join('\n') }}</pre>
				<p v-else>No errors</p>
			</div>
			<div>
				<h4 class="text-teal-500">Archive logs</h4>
				<div v-for="file in fileLogs" :key="file" class="cursor-pointer" @click="loadLogFile(file)">
					<UIcon name="i-heroicons-document" class="mr-2" />
					{{ file }}
				</div>
				<h4 class="text-teal-500 mt-5">Archive errors</h4>
				<div v-for="file in fileErrors" :key="file" class="cursor-pointer" @click="loadLogFile(file)">
					<UIcon name="i-heroicons-document" class="mr-2" />
					{{ file }}
				</div>
			</div>
		</div>
		<UModal v-model="isOpen" fullscreen>
			<UCard>
				<template #header>
					<div class="flex items-center justify-between">
						<h1 class="text-teal-500 font-bold mb-3"><UIcon name="i-heroicons-document" class="mr-2" />{{ logFile }}</h1>
						<UButton color="gray" variant="ghost" icon="i-heroicons-x-mark-20-solid" class="-my-1" @click="isOpen = false" />
					</div>
				</template>
				<pre class="text-xs">
					{{ logFileContent }}
				</pre
				>
			</UCard>
		</UModal>
	</div>
</template>

<script lang="ts" setup>
	const isOpen = ref(false)
	const logFile = ref('')
	const logFileContent = ref('')

	const loadLogFile = async (file: string) => {
		logFile.value = file
		logFileContent.value = await useApi().getLogFile(file)
		isOpen.value = true
	}
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

	const scheduleLogs = Object.keys(useLogs().out).map((l) => {
		return { slot: 'content', label: `Schedule: ${getLabelByScheduleId(l)}`, content: useLogs().out[l] }
	})
	const scheduleErrors = Object.keys(useLogs().out).map((l) => {
		return { slot: 'content', label: `Schedule: ${getLabelByScheduleId(l)}`, content: useLogs().err[l] }
	})

	onMounted(async () => {
		const { logs, errors } = await useApi().getLogs()
		fileLogs.value = logs
		fileErrors.value = errors
	})

	const colorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'bg-gray-950' : 'bg-white'
	})
</script>
