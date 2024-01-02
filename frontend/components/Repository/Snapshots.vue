<template>
	<div class="collapse bg-base-200 mb-5">
		<input type="radio" name="repository-accordion" checked />
		<h3 class="collapse-title m-0 text-primary"><FaIcon icon="table-list" class="mr-2" />Snapshots</h3>
		<div class="collapse-content">
			<div v-if="loading" class="text-center text-primary"><span class="loading loading-bars loading-lg"></span></div>
			<table v-else class="table">
				<thead>
					<th width="100">ID</th>
					<th width="100">Host</th>
					<th width="160">Time</th>
					<th>Paths</th>
					<th width="100">Tags</th>
				</thead>
				<tr class="hover transition-all" v-for="snapshot in snapshots" :key="snapshot.short_id">
					<td>
						<span class="">{{ snapshot.id.slice(0, 8) }}</span>
					</td>
					<td>{{ snapshot.hostname }}</td>
					<td>{{ format(new Date(snapshot.time), 'dd.MM.yyyy H:I:s') }}</td>
					<td>
						<span>{{ snapshot.paths.join(',') }}</span>
					</td>
					<td>
						<span v-for="tag in snapshot.tags" class="badge badge-info badge-outline">{{ tag }}</span>
					</td>
				</tr>
			</table>
		</div>
	</div>
</template>

<script setup lang="ts">
	import { onMounted } from 'vue'
	import { format } from 'date-fns'

	const snapshots = ref<Array<Snapshot>>([])
	const loading = ref(true)

	onMounted(async () => {
		snapshots.value = await Snapshots(useRoute().params.id as string)
		loading.value = false
	})
</script>
