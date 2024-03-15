<template>
	<div>
		<div class="flex justify-between">
			<div>
				<h1 class="text-green-500 font-bold"><UIcon name="i-heroicons-cog-6-tooth" class="mr-2" />Settings</h1>
			</div>
		</div>

		<div class="grid grid-cols-2 gap-10 p-10 bg-opacity-70 rounded-lg shadow-lg mt-5" :class="colorClass">
			<div>
				<h4 class="text-green-500 mb-2">Theme</h4>
				<USelect v-model="theme" :options="['system', 'light', 'dark']" />
				<h4 class="text-green-500 mb-2 mt-5">Hooks</h4>
				<p class="mb-3">Hooks can be used to run custom scripts on specific events.</p>
				<div class="text-sm" :class="textColorClass">Execute on schedule start</div>
				<UInput v-model="hookOnScheduleStart" placeholder="/absolute/path/to/start.sh" />
				<div class="text-sm mt-3" :class="textColorClass">Execute when schedule finishes successfully</div>
				<UInput v-model="hookOnScheduleSuccess" placeholder="/absolute/path/to/success.sh" />
				<div class="text-sm mt-3" :class="textColorClass">Execute when schedule finishes with errors</div>
				<UInput v-model="hookOnScheduleError" placeholder="/absolute/path/to/error.sh" />
			</div>
			<div>
				<h4 class="text-green-500 mb-2">Notifications</h4>
				<p class="mb-3" :class="textColorClass">Send Desktop notifications on specific events.</p>
				<UCheckbox v-model="notifiyOnScheduleStart" name="notifiyOnScheduleStart" color="green" label="Notify on schedule start" />
				<UCheckbox v-model="notifiyOnScheduleSuccess" name="notifiyOnScheduleSuccess" color="green" label="Notify when schedule finishes successfully" />
				<UCheckbox v-model="notifiyOnScheduleError" name="notifiyOnScheduleError" color="green" label="Notify when schedule finishes with errors" />
				<h4 class="text-green-500 mb-2 mt-5">Preserve error log files for X days.</h4>
				<UInput placeholder="7" v-model="preserveErrorLogsDays" />
				<UAlert title="Notes" class="mt-5" icon="i-heroicons-information-circle">
					<template #description>
						<ul>
							<li>- Notifications only work on Desktop.</li>
							<li>- Every hook gets a JSON object passed to it.</li>
						</ul>
						<pre class="mt-2">{schedule: {}, "backup": {}, "to_repository": {}, "from_repository": {} }</pre>
					</template>
				</UAlert>
			</div>
		</div>
		<div class="text-xs text-center mt-10">
			Resticity<br />Version: {{ version }}<br />Build: {{ build }} <br />Server: {{ `${useRequestURL().protocol}//${useRequestURL().host}` }}
		</div>
	</div>
</template>

<script lang="ts" setup>
	import _ from 'lodash'
	const theme = ref('auto')
	const notifiyOnScheduleError = ref(false)
	const notifiyOnScheduleSuccess = ref(false)
	const notifiyOnScheduleStart = ref(false)

	const hookOnScheduleError = ref('')
	const hookOnScheduleSuccess = ref('')
	const hookOnScheduleStart = ref('')

	const preserveErrorLogsDays = ref(7)

	const version = ref('')
	const build = ref('')

	const update = _.debounce(() => {
		saveSettings()
	}, 300)

	onMounted(async () => {
		theme.value = useSettings().settings.app_settings.theme
		notifiyOnScheduleError.value = useSettings().settings.app_settings.notifications.on_schedule_error
		notifiyOnScheduleStart.value = useSettings().settings.app_settings.notifications.on_schedule_start
		notifiyOnScheduleSuccess.value = useSettings().settings.app_settings.notifications.on_schedule_success
		hookOnScheduleError.value = useSettings().settings.app_settings.hooks.on_schedule_error
		hookOnScheduleStart.value = useSettings().settings.app_settings.hooks.on_schedule_start
		hookOnScheduleSuccess.value = useSettings().settings.app_settings.hooks.on_schedule_success
		preserveErrorLogsDays.value = useSettings().settings.app_settings.preserve_error_logs_days
		watch(
			[theme, notifiyOnScheduleError, notifiyOnScheduleStart, notifiyOnScheduleSuccess, hookOnScheduleError, hookOnScheduleStart, hookOnScheduleSuccess, preserveErrorLogsDays],
			async () => {
				update()
				useColorMode().preference = theme.value
			}
		)
		const vb = await useApi().getVersion()
		version.value = vb.version
		build.value = vb.build
	})

	function saveSettings() {
		useSettings().settings.app_settings = {
			theme: theme.value,
			notifications: {
				on_schedule_error: notifiyOnScheduleError.value,
				on_schedule_start: notifiyOnScheduleStart.value,
				on_schedule_success: notifiyOnScheduleSuccess.value,
			},
			hooks: {
				on_schedule_error: hookOnScheduleError.value,
				on_schedule_start: hookOnScheduleStart.value,
				on_schedule_success: hookOnScheduleSuccess.value,
			},
		}
		useSettings().save()
	}

	const colorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'bg-gray-950' : 'bg-white'
	})

	const textColorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'opacity-50' : 'text-black'
	})
</script>
