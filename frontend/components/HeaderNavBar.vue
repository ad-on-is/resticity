<template>
	<div class="bg-base-300 relative">
		<div class="navbar md:container md:mx-auto">
			<div class="navbar-start"><Logo class="h-10 w-auto fill-white stroke-white" /><span class="ml-3 font-bold text-2xl">resticity</span></div>
			<div class="navbar-center hidden lg:flex">
				<ul class="menu menu-horizontal px-1 gap-3">
					<li>
						<NuxtLink to="/"><FaIcon icon="fire" />Dashboard</NuxtLink>
					</li>
					<li>
						<NuxtLink to="/repositories" :class="useRoute().path.includes('repositories') ? 'active' : ''"><FaIcon icon="fa-server" />Repositories</NuxtLink>
					</li>
					<li>
						<NuxtLink to="/backups" :class="useRoute().path.includes('backups') ? 'active' : ''"><FaIcon icon="fa-upload" />Backups</NuxtLink>
					</li>
					<li>
						<NuxtLink to="/settings"><FaIcon icon="fa-gears" />Settings</NuxtLink>
					</li>
				</ul>
			</div>
			<div class="navbar-end"></div>
		</div>
	</div>
</template>

<script lang="ts" setup></script>
