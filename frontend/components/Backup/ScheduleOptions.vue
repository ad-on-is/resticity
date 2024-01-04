<template>
	<div class="collapse bg-base-200 mb-5">
		<input type="radio" name="backup-accordion" />
		<h3 class="collapse-title m-0 text-warning"><FaIcon icon="clock" class="mr-2" />Schedule</h3>

		<div class="collapse-content">
			<p>Set a cronjob for this backup.</p>
			<div class="join">
				<select class="select select-bordered select-sm join-item w-48" v-model="predefined">
					<option value="">Never</option>
					<option value="* * * * *">Every minute</option>
					<option value="0 * * * *">Every hour</option>
					<option value="0 */2 * * *">Every 2 hours</option>
					<option value="0 0 * * *">Every day</option>
					<option value="0 8 * * *">Every day at 8 am</option>
					<option value="custom">Custom</option>
				</select>
				<input class="input input-bordered join-item input-sm w-48 disabled:input-bordered" placeholder="" :disabled="predefined !== 'custom'" v-model="cron" />
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	const props = defineProps({
		cron: {
			type: String,
			default: '',
		},
	})
	const predefined = ref<string>(props.cron === '' ? '' : 'custom')
	const cron = ref<string>(props.cron)
	const emit = defineEmits(['update'])
	watch([predefined, cron], () => {
		if (predefined.value !== 'custom') {
			cron.value = predefined.value
		}
		emit('update', cron.value)
	})
</script>
