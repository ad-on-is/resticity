<template>
	<h1 class="text-primary"><FaIcon icon="fa-server" class="mr-3" />Repositories</h1>
	<div class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
		<NuxtLink
			to="#"
			class="opacity-40 border border-dashed border-success border-opacity-40 hover:opacity-100 shadow-lg bg-base-300 rounded-lg no-underline hover:bg-success transition-all hover:bg-opacity-10"
			@click="dialog?.showModal()"
		>
			<div class="p-5">
				<h3 class="m-0 text-success"><FaIcon icon="fa-plus-circle" class="mr-2" />Add Repository</h3>
				<p class="text-sm opacity-40">Initialize a new or connect an existing repository</p>
			</div>
		</NuxtLink>
		<NuxtLink
			:to="`/repositories/${repo.id}`"
			v-for="repo in useSettings().settings?.repositories"
			class="shadow-lg bg-base-300 rounded-lg no-underline hover:bg-primary transition-all hover:bg-opacity-10"
		>
			<div class="p-5">
				<h3 class="m-0 text-info"><FaIcon icon="fa-hard-drive" class="mr-2" />{{ repo.name }}</h3>
				<p class="text-sm break-words">{{ repo.path }}</p>
			</div>
		</NuxtLink>
	</div>

	<dialog class="modal" ref="dialog">
		<div class="modal-box">
			<h2 class="m-0 mb-3">New repository</h2>
			<p v-if="checkStatus === ''" class="text-sm opacity-40">Initialize a new or connect an existing repository</p>
			<p v-if="checkStatus !== '' && checkStatus.includes('REPO_OK')" class="text-sm text-success"><FaIcon icon="check" class="mr-2" />All checks passed</p>
			<p v-if="checkStatus !== '' && !checkStatus.includes('REPO_OK')" class="text-sm text-error"><FaIcon icon="warning" class="mr-2" />{{ checkStatus }}</p>
			<div class="">
				<input class="input input-bordered input-sm mb-5 min-w-full" placeholder="Name" v-model="newRepository.name" />

				<div class="join min-w-full">
					<input :type="pwType" class="flex-grow input input-bordered join-item input-sm mb-5" placeholder="Password" v-model="newRepository.password" />
					<button class="btn join-item btn-sm input-bordered" @click="togglePw"><FaIcon icon="fa-eye" /></button>
				</div>
				<div class="join w-full">
					<input class="input input-bordered join-item input-sm flex-grow" disabled placeholder="Location" v-model="newRepository.path" />
					<button class="btn join-item btn-sm input-bordered" @click="openDir()"><FaIcon icon="folder" /></button>
				</div>
			</div>
			<div class="modal-action">
				<button v-if="checkStatus.includes('REPO_OK')" class="btn btn-sm btn-success" @click="save"><FaIcon icon="plus-circle" />Add Repository</button>
				<button v-else class="btn btn-sm btn-warning btn-outline" @click="check">Check Repository</button>
			</div>
		</div>
	</dialog>
</template>

<script setup lang="ts">
	import { generateUUID } from '~/utils'

	const dialog = ref<HTMLDialogElement>()
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
		if (dir !== '') {
			checkStatus.value = ''
			newRepository.value.path = dir
		}
	}

	const check = async () => {
		checkStatus.value = await CheckRepository(newRepository.value)
	}

	const save = async () => {
		if (!checkStatus.value.includes('REPO_OK')) {
			return
		}
		if (checkStatus.value == 'REPO_OK_EMPTY') {
			const res = await InitializeRepository(newRepository.value)
			if (res === null) {
				useSettings().settings?.repositories.push(newRepository.value)
				await useSettings().save()
				dialog.value?.close()
			} else {
			}
		} else {
			useSettings().settings?.repositories.push(newRepository.value)
			await useSettings().save()
			dialog.value?.close()
		}
	}
</script>
