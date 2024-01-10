<template>
	<div>
		<h1 class="text-purple-500 font-bold"><UIcon name="i-heroicons-server" class="mr-3" />{{ repo.name }}</h1>
		<h2 class="mt-2 text-purple-500">{{ repo.path }}</h2>
		<UDivider class="my-5" />
		<div></div>
		<div>
			<RepositorySnapshots />
		</div>
		<UDivider class="my-10" />
		<div><RepositoryPruneOptions /></div>
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'

	const repo = ref<Repository>({})

	onMounted(async () => {
		repo.value = useSettings().settings?.repositories.find((r: Repository) => r.id === useRoute().params.id)
	})
</script>
