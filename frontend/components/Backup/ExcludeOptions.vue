<template>
	<h3 class="text-sky-500 mb-3"><UIcon name="i-heroicons-folder-minus" class="mr-2" />Exclude options</h3>
	<div class="grid grid-cols-2 gap-10 bg-gray-950 p-10 bg-opacity-70 rounded-lg shadow-lg">
		<div class="">
			<h4 class="text-indigo-500">Exclude files and folders</h4>
			<p class="opacity-50 text-sm mb-3">Pattern for files and folders to exclude. One per line.</p>
			<UTextarea :rows="7" placeholder="foo/**/bar" v-model="filesAndFolders"></UTextarea>
			<h4 class="text-indigo-500 mt-10">File</h4>
			<p class="opacity-50 text-sm mb-3">Exclude items listed in specific files.</p>
			<UTextarea :padded="true" :rows="7" placeholder="exclude.txt" v-model="listedInFiles"></UTextarea>
		</div>
		<div>
			<h4 class="text-indigo-500">Exclude if present</h4>
			<p class="opacity-50 text-sm mb-3">Excludes a folder if it contains any of these files.</p>
			<UTextarea :rows="7" placeholder=".nobackup" v-model="ifPresent"></UTextarea>
			<UCheckbox v-model="cacheDir" class="mt-10 mb-5">
				<template #label>
					<span
						>Exclude if
						<code class="text-lime-500">CACHEDIR.TAG</code>
						file is present</span
					>
				</template>
			</UCheckbox>

			<h4 class="text-indigo-500">Exclude files larger than</h4>
			<p class="opacity-50 text-sm mb-3">Exclude files if they exceed a specific file size</p>
			<UButtonGroup>
				<UInput v-model="largerThan" type="number" placeholder="0" />
				<USelect v-model="largerThanUnit" :options="units" option-attribute="name" class="w-20"></USelect>
			</UButtonGroup>
		</div>
	</div>
</template>

<script lang="ts" setup>
	const props = defineProps({
		excludes: {
			default: [[]],
		},
	})

	const units = [
		{ name: 'KiB', value: 'K' },
		{ name: 'MiB', value: 'M' },
		{ name: 'GiB', value: 'G' },
		{ name: 'TiB', value: 'T' },
	]

	const filesAndFolders = ref(fromPropsArray('--exclude'))
	const ifPresent = ref(fromPropsArray('--exclude-if-present'))
	const listedInFiles = ref(fromPropsArray('--exclude-file'))
	const cacheDir = ref(props.excludes.some((e) => e[0] === '--exclude-caches'))
	const largerThan = ref(fromPropsArray('--exclude-if-larger-than', '').replace(/[^0-9]/g, '') || 0)
	const largerThanUnit = ref(fromPropsArray('--exclude-if-larger-than', '').replace(/[0-9]/g, '') || 'K')
	const emit = defineEmits(['update'])

	// onMounted(() => {
	//
	// 	listedInFiles.value = props.excludes
	// 		.filter((e) => e[0] === '--exclude-file')
	// 		.map((e) => e[1])
	// 		.join('\n')
	// 	cacheDir.value = props.excludes.some((e) => e[0] === '--exclude-caches')
	// 	largerThan.value =
	// 		parseInt(
	// 			props.excludes
	// 				.filter((e) => e[0] === '--exclude-if-larger-than')
	// 				.map((e) => e[1])
	// 				.join('')
	// 				.replace(/[^0-9]/g, '')
	// 		) || 0
	// 	largerThanUnit.value = props.excludes
	// 		.filter((e) => e[0] === '--exclude-if-larger-than')
	// 		.map((e) => e[1])
	// 		.join('')
	// 		.replace(/[0-9]/g, '')
	// })

	function toParamArray(str: string, param: string): any {
		return str
			.split('\n')
			.map((f) => f.trim())
			.filter((f) => f !== '')
			.map((f) => [param, f])
	}

	function fromPropsArray(param: string, j: string = '\n'): string {
		return props.excludes
			.filter((e) => e[0] === param)
			.map((e) => e[1])
			.join(j)
	}

	watch([filesAndFolders, ifPresent, listedInFiles, cacheDir, largerThan, largerThanUnit], () => {
		emit('update', [
			...toParamArray(filesAndFolders.value, '--exclude'),
			...toParamArray(ifPresent.value, '--exclude-if-present'),
			...toParamArray(listedInFiles.value, '--exclude-file'),
			...(cacheDir.value ? [['--exclude-caches', '']] : []),
			...(parseInt(largerThan.value as string) > 0 ? [['--exclude-if-larger-than', `${largerThan.value}${largerThanUnit.value}`]] : []),
		])
	})
</script>
