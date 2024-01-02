<template>
	<div class="collapse bg-base-200 mb-5">
		<input type="radio" name="backup-accordion" />
		<h3 class="collapse-title m-0 text-error"><FaIcon icon="folder-minus" class="mr-2" />Exclude</h3>
		<div class="collapse-content">
			<div class="grid grid-cols-2 gap-10">
				<div>
					<h4 class="text-primary">Files and Folders</h4>
					<p class="opacity-50">Pattern for files and folders to exclude. One per line.</p>
					<textarea class="textarea textarea-bordered w-full h-32" placeholder="foo/**/bar" v-model="filesAndFolders"></textarea>
					<h4 class="text-primary">File</h4>
					<p class="opacity-50">Exclude items listed in specific files.</p>
					<textarea class="textarea textarea-bordered w-full h-32" placeholder="exclude.txt" v-model="listedInFiles"></textarea>
				</div>
				<div>
					<h4 class="text-primary">Exclude if present</h4>
					<p class="opacity-50">Excludes a folder if it contains any of these files.</p>
					<textarea class="textarea textarea-bordered w-full h-32" placeholder=".nobackup" v-model="ifPresent"></textarea>
					<div class="form-control">
						<label class="cursor-pointer label justify-normal">
							<input type="checkbox" class="checkbox checkbox-info mr-3" v-model="cacheDir" />
							<span class="label-text"
								>Exclude if
								<code class="text-warning">CACHEDIR.TAG</code>
								file is present</span
							>
						</label>
					</div>
					<h4 class="text-primary">Exclude files larger than</h4>
					<p class="opacity-50">Exclude files if they exceed a specific file size</p>
					<div class="join">
						<input class="input input-bordered join-item input-sm w-32" placeholder="0" v-model="largerThan" />
						<select class="select select-bordered select-sm join-item w-32" v-model="largerThanUnit">
							<option value="K" selected>KiB</option>
							<option value="M">MiB</option>
							<option value="G">GiB</option>
							<option value="T">TiB</option>
						</select>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	const filesAndFolders = ref('')
	const ifPresent = ref('')
	const listedInFiles = ref('')
	const cacheDir = ref(false)
	const largerThan = ref(0)
	const largerThanUnit = ref('K')
	const emit = defineEmits(['update'])

	function toParamArray(str: string, param: string): any {
		return str
			.split('\n')
			.map((f) => f.trim())
			.filter((f) => f !== '')
			.map((f) => [param, f])
	}

	watch([filesAndFolders, ifPresent, listedInFiles, cacheDir, largerThan, largerThanUnit], () => {
		emit('update', [
			...toParamArray(filesAndFolders.value, '--exclude'),
			...toParamArray(ifPresent.value, '--exclude-if-present'),
			...toParamArray(listedInFiles.value, '--exclude-file'),
			...(cacheDir.value ? [['--exclude-caches', '']] : []),
			...(largerThan.value > 0 ? [['--exclude-if-larger-than', `${largerThan.value}${largerThanUnit.value}`]] : []),
		])
	})
</script>
