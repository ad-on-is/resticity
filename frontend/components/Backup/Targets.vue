<template>
	<div class="collapse bg-base-200 mb-5">
		<input type="radio" name="backup-accordion" />
		<h3 class="collapse-title m-0 text-success"><FaIcon icon="server" class="mr-2" />Targets</h3>

		<div class="collapse-content">
			<div class="grid grid-cols-5 gap-5">
				<div
					v-for="repo in useSettings().settings?.repositories"
					class="shadow-lg bg-base-300 rounded-lg no-underline hover:bg-primary transition-all hover:bg-opacity-10 cursor-pointer"
					@click="toggleTarget(repo.id)"
					:key="repo.id"
				>
					<div class="p-3">
						<div class="form-control">
							<span class="label cursor-pointer justify-normal">
								<input type="checkbox" :checked="isSelected(repo.id)" class="checkbox checkbox-xs mr-3" :class="isSelected(repo.id) ? 'checkbox-info' : ''" />
								<span class="label-text" :class="isSelected(repo.id) ? 'text-info' : ''"><FaIcon icon="fa-hard-drive" class="mr-2" size="sm" />{{ repo.name }}</span>
							</span>
						</div>
						<p class="text-xs break-words m-0 opacity-40">{{ repo.path }}</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
	const props = defineProps({
		targets: {
			type: Array as PropType<string[]>,
			default: [],
		},
	})
	const targets = ref<string[]>(props.targets)
	const emit = defineEmits(['update'])
	function isSelected(id: string) {
		return targets.value.includes(id)
	}
	function toggleTarget(id: string) {
		if (targets.value.includes(id)) {
			targets.value = targets.value.filter((t) => t !== id)
		} else {
			targets.value.push(id)
		}

		emit('update', targets.value)
	}
</script>
