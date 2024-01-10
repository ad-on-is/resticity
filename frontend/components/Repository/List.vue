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
		<UCard>
			<template #header>
				<h2 class="text-purple-500 font-bold">New Repository</h2>
				<p v-if="checkStatus === ''" class="text-sm opacity-40">Initialize a new or connect an existing repository</p>

				<p v-if="checkStatus !== '' && checkStatus.includes('REPO_OK')" class="text-sm text-success">
					<FaIcon icon="check" class="mr-2" />All checks passed: {{ checkStatus === 'REPO_OK_EMPTY' ? 'Folder is empty' : 'Repository is valid' }}
				</p>
				<p v-if="checkStatus !== '' && !checkStatus.includes('REPO_OK')" class="text-sm text-error"><FaIcon icon="warning" class="mr-2" />Cannot use this repository</p>
			</template>
			<template #footer>
				<UButton
					v-if="checkStatus.includes('REPO_OK')"
					@click="save"
					color="purple"
					icon="i-heroicons-plus-circle"
					:disabled="newRepository.path === '' || newRepository.name === '' || newRepository.password === ''"
					>{{ checkStatus === 'REPO_OK_EMPTY' ? 'Initialize' : 'Add' }} Repository</UButton
				>
				<UButton v-else @click="check" color="yellow" variant="outline" icon="i-heroicons-plus-circle" :disabled="newRepository.path === '' || newRepository.password === ''"
					>Check Repository</UButton
				>
			</template>
			<UInput variant="outline" v-model="newRepository.name" placeholder="Name" class="mb-5" />
			<UButtonGroup class="flex">
				<UInput v-model="newRepository.path" placeholder="/path/to/repository" class="flex-grow" />
				<UButton icon="i-heroicons-folder-open" color="indigo" @click="openDir()" />
			</UButtonGroup>
			<p class="text-xs opacity-70">Path must be either an empty folder or an existing repository</p>
			<UButtonGroup class="flex mt-5">
				<UInput v-model="newRepository.password" :type="pwType" placeholder="Password" class="flex-grow" />
				<UButton icon="i-heroicons-eye" color="gray" @click="togglePw" />
			</UButtonGroup>
		</UCard>
	</UModal>
</template>

<script setup lang="ts">
	import { generateUUID } from '~/utils'
	const isOpen = ref(false)
	const pwType = ref('password')
	const newRepository = ref({
		id: generateUUID(),
		name: '',
		password: '',
		path: '',
		type: 0,
		prune_params: [],
		options: {},
	})
	const checkStatus = ref('')
	const togglePw = () => {
		pwType.value = pwType.value === 'password' ? 'text' : 'password'
	}
	const openDir = async () => {
		const dir = await SelectDirectory('Select a repository')
		if (useSettings().settings?.repositories.find((repo: any) => repo.path === dir)) {
			checkStatus.value = `${dir} is already a repository`
			return
		}
		newRepository.value.path = dir
	}

	watch(
		() => JSON.stringify(newRepository.value),
		(_, o) => {
			const old = JSON.parse(o)
			if (old.path !== newRepository.value.path) {
				checkStatus.value = ''
			}
		}
	)

	const check = async () => {
		checkStatus.value = await useApi().checkRepository(newRepository.value)
	}

	const save = async () => {
		if (!checkStatus.value.includes('REPO_OK')) {
			return
		}
		if (checkStatus.value === 'REPO_OK_EMPTY') {
			const res = await useApi().initRepository(newRepository.value)
			if (res === 'OK') {
				useSettings().settings?.repositories.push(newRepository.value)
				await useSettings().save()
				isOpen.value = false
			}
		} else {
			useSettings().settings?.repositories.push(newRepository.value)
			await useSettings().save()
			isOpen.value = false
		}
	}
</script>
