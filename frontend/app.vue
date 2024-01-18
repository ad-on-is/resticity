<template>
	<div class="min-h-screen overflow-hidden bg-cool pb-40" data-theme="resticity">
		<div v-if="loading"><Logo class="h-8 w-auto fill-orange-500 stroke-orange-500" /></div>
		<NuxtLayout v-else>
			<NuxtPage />
		</NuxtLayout>
		<UNotifications />
	</div>
</template>

<script setup lang="ts">
	const loading = ref(true)
	onMounted(async () => {
		await useSettings().init()
		await useSocket().init()
		loading.value = false
	})
</script>

<style>
	/* * {
		font-family: 'Hack';
	} */
	.page-enter-active,
	.page-leave-active {
		transition: all 0.2s;
	}
	.page-enter-from,
	.page-leave-to {
		opacity: 0;
		filter: blur(1rem);
	}
</style>
