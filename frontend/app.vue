<template>
	<div class="min-h-screen overflow-hidden bg-cool" data-theme="resticity">
		<NuxtLayout>
			<NuxtPage />
		</NuxtLayout>
	</div>
</template>

<script setup lang="ts">
	onMounted(() => {
		useSettings().init()
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
