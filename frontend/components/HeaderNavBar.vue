<template>
	<div class="bg-base-300 py-3">
		<div class="md:container md:mx-auto flex justify-between">
			<div class="flex"><Logo class="h-8 w-auto fill-orange-500 stroke-orange-500" /><span class="ml-2 font-bold text-2xl text-orange-500">resticity</span></div>
			<div class="">
				<ul class="flex">
					<li>
						<!-- <NuxtLink to="/"><FaIcon icon="fire" />Dashboard</NuxtLink> -->
						<UButton to="/" variant="ghost" :color="useRoute().path === '/' ? 'orange' : 'gray'" icon="i-heroicons-fire">Dashboard</UButton>
					</li>
					<li>
						<UButton to="/repositories" variant="ghost" :color="useRoute().path.includes('repositories') ? 'purple' : 'gray'" icon="i-heroicons-server-stack">Repositories</UButton>
						<!-- <NuxtLink to="/repositories" :class="useRoute().path.includes('repositories') ? 'active' : ''"><FaIcon icon="fa-server" />Repositories</NuxtLink> -->
					</li>
					<li>
						<UButton to="/backups" variant="ghost" :color="useRoute().path.includes('backups') ? 'sky' : 'gray'" icon="i-heroicons-arrow-up-tray">Backups</UButton>

						<!-- <NuxtLink to="/backups" :class="useRoute().path.includes('backups') ? 'active' : ''"><FaIcon icon="fa-upload" />Backups</NuxtLink> -->
					</li>
					<li>
						<UButton to="/schedules" variant="ghost" :color="useRoute().path.includes('schedules') ? 'yellow' : 'gray'" icon="i-heroicons-clock"> Schedules </UButton>

						<!-- <NuxtLink to="/schedules" :class="useRoute().path.includes('schedules') ? 'active' : ''"><FaIcon icon="clock" />Schedules</NuxtLink> -->
					</li>
					<li>
						<UButton to="/logs" variant="ghost" :color="useRoute().path.includes('logs') ? 'blue' : 'gray'" icon="i-heroicons-clock">Logs</UButton>

						<!-- <NuxtLink to="/schedules" :class="useRoute().path.includes('schedules') ? 'active' : ''"><FaIcon icon="clock" />Schedules</NuxtLink> -->
					</li>
				</ul>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup></script>
