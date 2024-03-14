<template>
	<div class="py-3" :class="colorClass">
		<div class="md:container md:mx-auto flex justify-between">
			<div class=""><Logo class="h-10 w-10" /></div>
			<div class="">
				<ul class="flex">
					<li>
						<UButton to="/" variant="ghost" :color="useRoute().path === '/' ? 'orange' : 'gray'" icon="i-heroicons-fire">Dashboard</UButton>
					</li>
					<li>
						<UButton to="/repositories" variant="ghost" :color="useRoute().path.includes('repositories') ? 'purple' : 'gray'" icon="i-heroicons-server-stack">Repositories</UButton>
					</li>
					<li>
						<UButton to="/backups" variant="ghost" :color="useRoute().path.includes('backups') ? 'sky' : 'gray'" icon="i-heroicons-arrow-up-tray">Backups</UButton>
					</li>
					<li>
						<UButton to="/schedules" variant="ghost" :color="useRoute().path.includes('schedules') ? 'yellow' : 'gray'" icon="i-heroicons-clock"> Schedules </UButton>
					</li>
					<li>
						<UButton to="/logs" variant="ghost" :color="useRoute().path.includes('logs') ? 'teal' : 'gray'" icon="i-heroicons-document-text">Logs</UButton>
					</li>
					<li>
						<UButton to="/settings" variant="ghost" :color="useRoute().path.includes('settings') ? 'green' : 'gray'" icon="i-heroicons-cog-6-tooth">Settings</UButton>
					</li>
				</ul>
			</div>
		</div>
		<div class="bg-purple-800"></div>
	</div>
</template>

<script lang="ts" setup>
	const colorClass = computed(() => {
		return useColorMode().value === 'dark' ? 'bg-gray-950' : 'bg-gray-950 bg-opacity-10'
	})
</script>
