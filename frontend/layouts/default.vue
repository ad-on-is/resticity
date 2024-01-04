<template>
	<div>
		<HeaderNavBar />
		<div class="md:container md:mx-auto prose pt-5">
			<slot />
		</div>
	</div>
</template>

<script lang="ts" setup>
	useJobs().init()
</script>
