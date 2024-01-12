<template>
	<div>
		<h1 class="text-purple-500 font-bold"><UIcon name="i-heroicons-server" class="mr-3" />{{ repo.name }}</h1>
		<h2 class="my-2 opacity-50">{{ repo.path }}</h2>
		<div v-if="stats" class="flex text-xs gap-2">
			<UBadge color="gray"><UIcon name="i-heroicons-server" class="mr-1" />{{ humanFileSize(stats.total_size) }} used</UBadge>
			<UBadge color="gray"><UIcon name="i-heroicons-document-duplicate" class="mr-1" />{{ stats.total_file_count }} files</UBadge>
			<UBadge color="gray"><UIcon name="i-heroicons-queue-list" class="mr-1" />{{ stats.snapshots_count }} snapshots</UBadge>
		</div>

		<UButton v-if="mountPath === ''" color="indigo" variant="outline" icon="i-heroicons-folder" @click="mount">Mount</UButton>
		<UButtonGroup v-else>
			<UButton color="gray" disabled icon="i-heroicons-folder">{{ mountPath }}</UButton>
			<UButton @click="unmount" color="indigo">Unmount</UButton>
		</UButtonGroup>
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
	const stats = ref()
	const mountPath = ref('')

	const mount = async () => {
		const dir = await SelectDirectory('Select a mount point')
		if (!dir) return
		mountPath.value = dir
		useSettings().settings.mounts.push({
			id: useRoute().params.id as string,
			path: mountPath.value,
		})
		useSettings().save()
		useApi().mount(useRoute().params.id as string, mountPath.value)
	}
	const unmount = () => {
		useSettings().settings.mounts = useSettings().settings.mounts.filter((m: Mount) => m.id !== useRoute().params.id)
		useSettings().save()
		useApi().unmount(useRoute().params.id as string, mountPath.value)
		mountPath.value = ''
	}

	const repo = ref<Repository>({})

	onMounted(async () => {
		repo.value = useSettings().settings?.repositories.find((r: Repository) => r.id === useRoute().params.id)
		// stats.value = await useApi().statRepository(useRoute().params.id as string)
		console.log(stats.value)
		mountPath.value = useSettings().settings.mounts.find((m: Mount) => m.id === useRoute().params.id)?.path ?? ''
	})
</script>
