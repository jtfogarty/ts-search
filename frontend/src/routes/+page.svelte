<script lang="ts">
	import SideBar from '$lib/components/SideBar.svelte';
	import Collections from '$lib/components/Collections.svelte';
</script>

<div class="flex max-w-[1700px] mx-auto w-full">
	<SideBar />
	<div class="lg:ml-[80px] flex-1 overflow-clip">
		<Collections />
	</div>
</div>