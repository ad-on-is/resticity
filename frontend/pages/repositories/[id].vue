<template>
	<div>
		<h1 class="text-info m-0"><FaIcon icon="hard-drive" class="mr-3" />{{ repo.name }}</h1>
		<h2 class="mt-2 opacity-40 text-info">{{ repo.path }}</h2>
		<RepositoryPruneOptions />
		<RepositorySnapshots :repo="repo" />
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'

	const repo = ref<Repository>({})

	onMounted(async () => {
		repo.value = useSettings().settings?.repositories.find((r: Repository) => r.id === useRoute().params.id)
	})
</script>
