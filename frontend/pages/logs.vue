<template>
	<div>
		<h1>Out</h1>
		<UTextarea v-model="useLogs().out" :rows="20" disabled></UTextarea>

		<h1>Error</h1>
		<UTextarea v-model="useLogs().err" :rows="20" color="red" disabled></UTextarea>
	</div>
</template>
