<template>
	<div class="collapse bg-base-200 mb-5">
		<input type="radio" name="repository-accordion" />
		<h3 class="collapse-title m-0 text-warning"><FaIcon icon="hand-sparkles" class="mr-2" />Cleanup Policies</h3>

		<div class="collapse-content">
			<p>Define the policies for snapshot retentions.</p>
			<div class="grid grid-cols-2 gap-5">
				<div>
					<h4 class="text-primary">Keep tags</h4>
					<p>Specify the tags to keep. One per line</p>
					<textarea class="textarea textarea-bordered w-full" placeholder="tag"></textarea>
				</div>
				<div>
					<h4 class="text-primary">Keep</h4>
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
					<h4 class="text-primary">Keep within</h4>
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
		</div>
	</div>
</template>
