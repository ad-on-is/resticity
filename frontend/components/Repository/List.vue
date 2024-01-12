<template>
	<h1 class="text-purple-500 font-bold mb-3"><UIcon name="i-heroicons-server-stack" class="mr-2" dynamic />Repositories</h1>
	<div class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
		<div
			class="cursor-pointer opacity-40 border border-dashed border-purple-500 border-opacity-40 hover:opacity-100 shadow-lg bg-base-300 rounded-lg no-underline hover:bg-purple-500 transition-all hover:bg-opacity-10"
			@click="isOpen = true"
		>
			<div class="p-5">
				<h3 class="m-0 text-purple-500 p-0"><FaIcon icon="fa-plus-circle" class="mr-2" />Add Repository</h3>
				<p class="text-sm opacity-40">Initialize a new or connect an existing repository</p>
			</div>
		</div>
		<NuxtLink
			:to="`/repositories/${repo.id}`"
			v-for="repo in useSettings().settings?.repositories"
			class="shadow-lg bg-base-300 rounded-lg no-underline hover:bg-purple-500 transition-all hover:bg-opacity-10"
		>
			<div class="p-5 pb-0">
				<h3 class="m-0 font-medium text-purple-500 p-0"><FaIcon icon="fa-hard-drive" class="mr-2" />{{ repo.name }}</h3>
				<p class="text-xs break-words opacity-50 p-0 m-0">{{ repo.path }}</p>
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
	const isOpen = ref(false)
</script>
