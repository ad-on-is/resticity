<template>
	<div class="bg-base-300 relative">
		<div class="navbar md:container md:mx-auto">
			<div class="navbar-start"><Logo class="h-10 w-auto fill-white stroke-white" /><span class="ml-3 font-bold text-2xl">resticity</span></div>
			<div class="navbar-center hidden lg:flex">
				<ul class="menu menu-horizontal px-1 gap-3">
					<li>
						<!-- <NuxtLink to="/"><FaIcon icon="fire" />Dashboard</NuxtLink> -->
						<UButton to="/" variant="ghost" color="gray" icon="i-heroicons-fire">Dashboard</UButton>
					</li>
					<li>
						<UButton
							to="/repositories"
							:variant="useRoute().path.includes('repositories') ? 'solid' : 'ghost'"
							:color="useRoute().path.includes('repositories') ? 'blue' : 'gray'"
							icon="i-heroicons-server-stack"
							>Repositories</UButton
						>
						<!-- <NuxtLink to="/repositories" :class="useRoute().path.includes('repositories') ? 'active' : ''"><FaIcon icon="fa-server" />Repositories</NuxtLink> -->
					</li>
					<li>
						<UButton to="/backups" variant="ghost" :color="useRoute().path.includes('backups') ? 'blue' : 'gray'" icon="i-heroicons-arrow-up-tray">Backups</UButton>

						<!-- <NuxtLink to="/backups" :class="useRoute().path.includes('backups') ? 'active' : ''"><FaIcon icon="fa-upload" />Backups</NuxtLink> -->
					</li>
					<li>
						<UButton to="/schedules" variant="ghost" :color="useRoute().path.includes('schedules') ? 'blue' : 'gray'" icon="i-heroicons-clock">Schedules</UButton>

						<!-- <NuxtLink to="/schedules" :class="useRoute().path.includes('schedules') ? 'active' : ''"><FaIcon icon="clock" />Schedules</NuxtLink> -->
					</li>
					<li>
						<UButton to="/logs" variant="ghost" :color="useRoute().path.includes('logs') ? 'blue' : 'gray'" icon="i-heroicons-clock">Logs</UButton>

						<!-- <NuxtLink to="/schedules" :class="useRoute().path.includes('schedules') ? 'active' : ''"><FaIcon icon="clock" />Schedules</NuxtLink> -->
					</li>
				</ul>
			</div>
			<div class="navbar-end"></div>
		</div>
	</div>
</template>

<script lang="ts" setup></script>
