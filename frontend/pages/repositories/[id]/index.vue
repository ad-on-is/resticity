<template>
	<div v-if="repo">
		<div class="flex justify-between">
			<div>
				<h1 class="text-purple-500 font-bold"><UIcon :name="getRepoIcon(repo)" class="mr-3" />{{ repo.name }}</h1>
				<h2 class="mb-5">{{ repo.path }}</h2>
			</div>

			<div class="mt-3">
				<UButton v-if="!useMounts().repoIsMounted(repo.id)" color="yellow" variant="outline" icon="i-heroicons-folder" @click="isOpen = true">Mount</UButton>

				<UButtonGroup v-else>
					<UButton color="gray" disabled icon="i-heroicons-folder">{{ useMounts().repoIsMounted(repo.id)?.path }}</UButton>
					<UButton @click="unmount" color="indigo">Unmount</UButton>
				</UButtonGroup>
				<UButton icon="i-heroicons-trash" color="red" class="ml-2" @click="openDelete = true">Delete</UButton>
			</div>
		</div>

		<UDivider class="my-5" />
		<div></div>
		<div>
			<RepositorySnapshots />
		</div>
		<UDivider class="my-10" />
		<div><RepositoryPruneOptions @update="(val) => (prunes = val)" :prunes="prunes" /></div>
		<UModal v-model="isOpen">
			<UCard>
				<template #header> Select a mount point </template>
				<PathAutocomplete @selected="(p) => (shouldMountPath = p)" />
				<template #footer><UButton color="orange" icon="i-heroicons-folder-open" @click="mount">Mount</UButton> {{ shouldMountPath }}</template>
			</UCard>
		</UModal>
		<UModal v-model="openDelete">
			<UCard>
				<template #header><span class="text-red-500">Delete repository</span> </template>
				<p>Do you really want to delete this repository and its associated schedules? <br /><span class="opacity-50 text-sm">This will not delete actual files.</span></p>
				<template #footer><UButton color="red" icon="i-heroicons-trash" @click="deleteRepo">Yes, delete</UButton></template>
			</UCard>
		</UModal>
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import _ from 'lodash'
	const isOpen = ref(false)
	const openDelete = ref(false)
	const mountPath = ref('')
	const shouldMountPath = ref('')
	const prunes = ref<[]>([])
	const idx = ref(-1)
	const deleteRepo = async () => {
		useSettings().settings!.repositories = useSettings().settings!.repositories.filter((item: Repository) => item.id !== repo.value.id)
		useSettings().settings!.schedules = useSettings().settings!.schedules.filter(
			(item: Schedule) => item.to_repository_id !== repo.value.id && item.from_repository_id !== repo.value.id
		)
		useSettings().save()
		openDelete.value = false
		return navigateTo('/repositories')
	}
	const mount = async () => {
		mountPath.value = shouldMountPath.value

		useApi().mount(useRoute().params.id as string, mountPath.value)

		isOpen.value = false
	}
	const unmount = () => {
		useApi().unmount(useRoute().params.id as string, mountPath.value)
		mountPath.value = ''
	}

	const update = _.debounce(() => {
		repo.value.prune_params = prunes.value
		useSettings().settings!.repositories[idx.value] = repo.value
		useSettings().save()
	}, 300)

	const repo = ref()

	onMounted(async () => {
		repo.value = useSettings().settings?.repositories.find((r: Repository) => r.id === useRoute().params.id)
		prunes.value = repo.value.prune_params
		idx.value = useSettings().settings!.repositories.findIndex((r: Repository) => r.id === repo.value.id)
		mountPath.value = useSettings().settings.mounts.find((m: Mount) => m.id === useRoute().params.id)?.path ?? ''
		watch(
			() => [JSON.stringify(prunes.value)],
			() => {
				update()
			}
		)
	})
</script>
