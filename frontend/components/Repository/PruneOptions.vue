<template>
	<h3 class="text-purple-500 mb-3"><UIcon name="i-heroicons-trash" class="mr-2" />Prune options</h3>
	<div class="grid grid-cols-2 gap-5 p-10 bg-opacity-70 rounded-lg shadow-lg" :class="colorClass">
		<div>
			<h4 class="text-indigo-500 font-medium">Keep tags</h4>
			<p class="text-xs mb-3">Specify the tags to keep. One per line</p>
			<UTextarea placeholder="tags" v-model="keep.tags"></UTextarea>
		</div>
		<div>
			<h4 class="text-indigo-500 font-medium">Keep</h4>
			<p class="text-xs mb-3">Keep range</p>
			<div class="grid grid-cols-3 gap-3">
				<UButtonGroup>
					<UInput v-model="keep.last" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Last</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keep.hourly" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Hourly</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keep.daily" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Daily</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keep.weekly" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Weekly</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keep.monthly" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Monthly</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keep.yearly" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Yearly</UButton>
				</UButtonGroup>
			</div>
			<h4 class="text-indigo-500 mt-5">Keep within</h4>
			<p class="text-xs mb-3">Keep within a range</p>
			<div class="grid grid-cols-3 gap-5">
				<UButtonGroup>
					<UInput v-model="keepWithin.last" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Last</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keepWithin.hourly" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Hourly</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keepWithin.daily" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Daily</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keepWithin.weekly" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Weekly</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keepWithin.monthly" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Monthly</UButton>
				</UButtonGroup>
				<UButtonGroup>
					<UInput v-model="keepWithin.yearly" placeholder="0" class="w-32" />
					<UButton disabled color="gray" class="!cursor-default w-20">Yearly</UButton>
				</UButtonGroup>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
	const props = defineProps({
		prunes: {
			default: [],
		},
	})

	const emit = defineEmits(['update'])

	const keep = ref({
		tags: fromPropsArray('--keep-tag'),
		last: fromPropsArrayToNumber('--keep-last'),
		hourly: fromPropsArrayToNumber('--keep-hourly'),
		daily: fromPropsArrayToNumber('--keep-daily'),
		weekly: fromPropsArrayToNumber('--keep-weekly'),
		monthly: fromPropsArrayToNumber('--keep-monthly'),
		yearly: fromPropsArrayToNumber('--keep-yearly'),
	})

	const keepWithin = ref({
		last: fromPropsArrayToNumber('--keep-within-last'),
		hourly: fromPropsArrayToNumber('--keep-within-hourly'),
		daily: fromPropsArrayToNumber('--keep-within-daily'),
		weekly: fromPropsArrayToNumber('--keep-within-weekly'),
		monthly: fromPropsArrayToNumber('--keep-within-monthly'),
		yearly: fromPropsArrayToNumber('--keep-within-yearly'),
	})

	function toParamArray(str: string, param: string): any {
		return str
			.split('\n')
			.map((f) => f.trim())
			.filter((f) => f !== '')
			.map((f) => [param, f])
	}

	function fromPropsArrayToNumber(param: string, j: string = '\n'): number {
		return props.prunes.filter((e) => e[0] === param).flat()[1]
	}

	function fromPropsArray(param: string, j: string = '\n'): string {
		return props.prunes
			.filter((e) => e[0] === param)
			.map((e) => e[1])
			.join(j)
	}

	watch(
		() => [JSON.stringify(keep), JSON.stringify(keepWithin)],
		() => {
			emit(
				'update',
				[
					...toParamArray(keep.value.tags, '--keep-tag'),
					keep.value.last > 0 ? ['--keep-last', keep.value.last] : [],
					keep.value.hourly > 0 ? ['--keep-hourly', keep.value.hourly] : [],
					keep.value.daily > 0 ? ['--keep-daily', keep.value.daily] : [],
					keep.value.weekly > 0 ? ['--keep-weekly', keep.value.weekly] : [],
					keep.value.monthly > 0 ? ['--keep-monthly', keep.value.monthly] : [],
					keep.value.yearly > 0 ? ['--keep-yearly', keep.value.yearly] : [],
					keepWithin.value.last > 0 ? ['--keep-within-last', keepWithin.value.last] : [],
					keepWithin.value.hourly > 0 ? ['--keep-within-hourly', keepWithin.value.hourly] : [],
					keepWithin.value.daily > 0 ? ['--keep-within-daily', keepWithin.value.daily] : [],
					keepWithin.value.weekly > 0 ? ['--keep-within-weekly', keepWithin.value.weekly] : [],
					keepWithin.value.monthly > 0 ? ['--keep-within-monthly', keepWithin.value.monthly] : [],
					keepWithin.value.yearly > 0 ? ['--keep-within-yearly', keepWithin.value.yearly] : [],
				].filter((e) => e.length > 0)
			)
		}
	)
	const colorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'bg-gray-950' : 'bg-white'
	})
</script>
