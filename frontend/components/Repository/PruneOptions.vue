<template>
	<h3 class="text-purple-500 mb-3"><UIcon name="i-heroicons-trash" class="mr-2" />Prune options</h3>
	<div class="grid grid-cols-2 gap-5 bg-gray-950 p-10 bg-opacity-70 rounded-lg shadow-lg">
		<div>
			<h4 class="text-indigo-500 font-medium">Keep tags</h4>
			<p class="text-xs mb-3">Specify the tags to keep. One per line</p>
			<UTextarea placeholder="tags"></UTextarea>
		</div>
		<div>
			<h4 class="text-indigo-500 font-medium">Keep</h4>
			<p class="text-xs mb-3">Keep range</p>
			<div class="grid grid-cols-3 gap-3">
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Last</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Hourly</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Daily</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Weekly</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Monthly</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Yearly</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
			</div>
			<h4 class="text-indigo-500 mt-5">Keep within</h4>
			<p class="text-xs mb-3">Keep within a range</p>
			<div class="grid grid-cols-3 gap-5">
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Within</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Hourly</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Daily</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Weekly</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Monthly</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
				<div class="join">
					<label class="bg-base-300 input input-bordered input-sm disabled join-item w-32">Yearly</label>
					<input class="input input-bordered join-item input-sm w-16" placeholder="0" />
				</div>
			</div>
		</div>
	</div>
</template>
