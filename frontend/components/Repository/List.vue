<template>
	<h1 class="text-purple-500 font-bold mb-3"><UIcon name="i-heroicons-server-stack" class="mr-2" />Repositories</h1>
	<div class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
		<div
			v-if="showNew"
			class="cursor-pointer opacity-40 border border-dashed border-purple-500 border-opacity-40 hover:opacity-100 shadow-lg rounded-lg no-underline hover:bg-purple-500 transition-all hover:bg-opacity-10"
			:class="colorClass"
			@click="isOpen = true"
		>
			<div class="p-5">
				<h3 class="m-0 text-purple-500 p-0"><FaIcon icon="fa-plus-circle" class="mr-2" />Add Repository</h3>
				<p class="text-sm" :class="textColorClass">Initialize a new or connect an existing repository</p>
			</div>
		</div>
		<NuxtLink
			:to="`/repositories/${repo.id}`"
			v-for="repo in useSettings().settings?.repositories"
			class="shadow-lg rounded-lg no-underline hover:bg-purple-500 transition-all hover:bg-opacity-10"
			:class="colorClass"
		>
			<div class="p-5 pb-0">
				<h3 class="m-0 font-medium text-purple-500 p-0"><UIcon :name="getRepoIcon(repo)" class="mr-2" />{{ repo.name }}</h3>
				<p class="text-xs break-words p-0 m-0" :class="textColorClass">{{ repo.path }}</p>
				<div :class="useJobs().repoIsSynching(repo.id) || useJobs().repoIsRunning(repo.id) ? 'opacity-100' : 'opacity-0'">
					<span class="loading loading-infinity loading-sm text-orange-500"></span>
				</div>
			</div>
		</NuxtLink>
	</div>

	<UModal v-model="isOpen">
		<RepositoryNew @finish="isOpen = false" />
	</UModal>
</template>

<script setup lang="ts">
	const props = defineProps({
		showNew: {
			type: Boolean,
			default: true,
		},
	})

	import { getRepoIcon } from '~/utils'

	const isOpen = ref(false)

	const colorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'bg-base-300' : 'bg-base-300 bg-opacity-10'
	})
	const textColorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'opacity-50' : 'text-black'
	})
</script>
