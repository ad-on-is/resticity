<template>
	<UCard>
		<template #header>
			<div v-if="shouldInit">
				<h2 class="text-green-500 font-bold">Initializing Repository</h2>
			</div>
			<div v-else>
				<h2 class="text-purple-500 font-bold">New Repository</h2>
				<p v-if="checkStatus === ''" class="text-sm opacity-40">Initialize a new or connect an existing repository</p>

				<p v-if="checkStatus !== '' && checkStatus.includes('OK_REPO')" class="text-sm text-success">
					<FaIcon icon="check" class="mr-2" />All checks passed: {{ checkStatus === 'OK_REPO_EMPTY' ? 'Folder is empty' : 'Repository is valid' }}
				</p>
				<p v-if="checkStatus !== '' && !checkStatus.includes('OK_REPO')" class="text-sm text-error"><FaIcon icon="warning" class="mr-2" />Cannot use this repository</p>
			</div>
		</template>
		<template #footer>
			<div v-if="!shouldInit">
				<UButton
					v-if="checkStatus.includes('OK_REPO')"
					@click="save"
					color="purple"
					icon="i-heroicons-plus-circle"
					:disabled="newRepository.path === '' || newRepository.name === '' || newRepository.password === '' || initializing"
					>{{ checkStatus === 'OK_REPO_EMPTY' ? 'Initialize' : 'Add' }} Repository</UButton
				>
				<UButton v-else @click="check" color="yellow" variant="outline" icon="i-heroicons-plus-circle" :disabled="newRepository.path === '' || newRepository.password === ''"
					>Check Repository</UButton
				>
			</div>
			<div v-else>
				<UButton @click="finish" color="green" :disabled="initializing">Finish</UButton>
			</div>
		</template>
		<div v-if="shouldInit">
			<div v-if="initializing">
				<p>
					<UIcon name="i-heroicons-server" class="text-purple-500" /> <span class="text-purple-500">{{ newRepository.path }} /home/test</span> is being prepared to be used as a
					repository.
				</p>
			</div>
			<div v-else>
				<p>
					<UIcon name="i-heroicons-server" class="text-purple-500" /> <span class="text-purple-500">{{ newRepository.path }} </span>
					{{ initStatus === 'OK' ? 'is ready!' : 'could not be initialized' }}
				</p>
			</div>
		</div>
		<div v-else>
			<UTabs :items="items" v-model="selectedTab">
				<template #local="{ item }">
					<div class="mt-5">
						<UInput variant="outline" v-model="newRepository.name" placeholder="Name" class="mb-5" />
						<PathAutocomplete @selected="(p) => (newRepository.path = p)" />
						<p class="text-xs opacity-70">Path must be either an empty folder or an existing repository</p>
						<UButtonGroup class="flex mt-5">
							<UInput v-model="newRepository.password" :type="pwType" placeholder="Password" class="flex-grow" />
							<UButton icon="i-heroicons-eye" color="gray" @click="togglePw" />
						</UButtonGroup>
					</div>
				</template>
				<template #s3="{ item }">
					<div class="mt-5">
						<UAlert icon="i-heroicons-exclamation-circle" title="Attention" description="Please make sure the bucket is empty." class="mb-5" color="yellow" />
						<UInput variant="outline" v-model="newRepository.name" placeholder="Name" class="mb-5" />
						<UInput v-model="newRepository.password" :type="pwType" placeholder="Password" class="flex-grow mb-5" />
						<UInput variant="outline" v-model="newRepository.path" placeholder="s3.example.com/bucket" class="mb-5" />

						<UInput v-model="newRepository.options.s3_key" placeholder="Access key" class="flex-grow" />
						<UButtonGroup class="flex mt-5">
							<UInput v-model="newRepository.options.s3_secret" :type="pwType" placeholder="Access secret" class="flex-grow" />
							<UButton icon="i-heroicons-eye" color="gray" @click="togglePw" />
						</UButtonGroup>
					</div>
				</template>
			</UTabs>
		</div>
	</UCard>
</template>

<script setup lang="ts">
	const emit = defineEmits(['finish'])
	import { generateUUID } from '~/utils'
	const shouldInit = ref(false)
	const initializing = ref(false)
	const pwType = ref('password')
	const items = [
		{
			slot: 'local',
			label: 'Local',
			icon: 'i-heroicons-server',
		},
		{
			slot: 's3',
			label: 'S3/B2',
			icon: 'i-heroicons-server',
		},
	]

	const emptyRepo = () => ({
		id: generateUUID(),
		type: 'local',
		name: '',
		password: '',
		path: '',
		prune_params: [],
		options: {
			s3_key: '',
			s3_secret: '',
		},
	})
	const selectedTab = computed({
		get() {
			return 0
		},
		set(value) {
			newRepository.value.type = items[value].slot
		},
	})
	const newRepository = ref(emptyRepo())
	const checkStatus = ref('')
	const initStatus = ref('')
	const togglePw = () => {
		pwType.value = pwType.value === 'password' ? 'text' : 'password'
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
		if (newRepository.value.type === 's3' && !newRepository.value.path.startsWith('s3:')) {
			newRepository.value.path = `s3:${newRepository.value.path}`
		}
		checkStatus.value = await useApi().checkRepository(newRepository.value)
	}

	const finish = () => {
		emit('finish')
		newRepository.value = emptyRepo()
		checkStatus.value = ''
		shouldInit.value = false
		initializing.value = false
	}

	const save = async () => {
		if (!checkStatus.value.includes('OK_REPO')) {
			return
		}
		if (checkStatus.value === 'OK_REPO_EMPTY') {
			shouldInit.value = true
			initializing.value = true
			initStatus.value = await useApi().initRepository(newRepository.value)
			if (initStatus.value === 'OK') {
				useSettings().settings?.repositories.push(newRepository.value)
				await useSettings().save()
			}
			initializing.value = false
		} else {
			useSettings().settings?.repositories.push(newRepository.value)
			await useSettings().save()
		}
		finish()
	}
</script>
