<template>
	<div>
		<div class="mt-0"><RepositoryList :show-new="false" /></div>
		<UDivider class="my-10" />
		<div><BackupList :show-new="false" /></div>
		<UDivider class="my-10" />
		<div><ScheduleList /></div>
	</div>
</template>
