<template>
	<div class="collapse bg-base-200 mb-5">
		<input type="radio" name="backup-accordion" />
		<h3 class="collapse-title m-0 text-warning"><FaIcon icon="clock" class="mr-2" />Schedule</h3>

		<div class="collapse-content">
			<p>Set a cronjob for this backup.</p>
			<div class="join">
				<select class="select select-bordered select-sm join-item w-32">
					<option value="* * * * *" selected>Never</option>
					<option value="* * * * *">Every hour</option>
					<option value="* * * * *">Every 2 hours</option>
					<option value="* * * * *">Every day</option>
					<option value="custom">Custom</option>
				</select>
				<input class="input input-bordered join-item input-sm w-48 disabled:input-bordered" placeholder="0" disabled />
			</div>
		</div>
	</div>
</template>
