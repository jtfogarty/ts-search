<script lang="ts">
	import { onMount } from 'svelte';
	import * as Popover from "$lib/components/ui/popover";
	import * as Command from "$lib/components/ui/command";
	import { Button } from "$lib/components/ui/button";
	import { cn } from "$lib/utils";
	
	import { GetShakespeareWorks } from '$lib/wailsjs/go/main/App';
	import type  ShakespeareWorks  from '$lib/wailsjs/go/models';

	let open = false; 
	let value = "";
	let searchTerm = "";
	let searchInput = "";
	
	type Book = {
		value: string;
		label: string;
	};
	
	let books: Book[] = [];
	
	async function loadBooks() {
		try {
			const works = await GetShakespeareWorks();
			
			books = [
				{ value: 'all', label: 'All' },
				...works.map((work: ShakespeareWorks) => ({
					value: work.id.toString(),
					label: work.title
				})).sort((a: { label: string; }, b: { label: any; }) => a.label.localeCompare(b.label)) // Sort books by label
			];
		} catch (error) {
			console.error('Failed to load Shakespeare works:', error);
			books = [];
		}
	}
	
	onMount(() => {
		loadBooks();
	});

	function handleSearch() {
		const selectedBook = books.find(book => book.value === value);
		const bookTitle = selectedBook ? selectedBook.label : "No book selected";
		alert(`Selected Book: ${bookTitle}\nSearch Term: ${searchInput}`);
	}
</script>

<div class="flex max-w-[1700px] mx-auto w-full overflow-hidden">
	<div class="flex-1 text-white lm-10 mt-3 overflow-hidden lg:ml-[90px] ">
		<Popover.Root bind:open>
			<Popover.Trigger asChild let:builder>
				<Button
					variant="outline"
					role="combobox"
					aria-expanded={open}
					class="w-[210px] justify-between"
					builders={[builder]}
				>
					{value
						? books.find((book) => book.value === value)?.label
						: "Select Shakespeare Book..."}
				</Button>
			</Popover.Trigger>
			<Popover.Content class="w-[250px] max-h-[300px] overflow-y-auto p-0">
				<Command.Root>
					<Command.Input placeholder="Search books..." class="h-9" bind:value={searchTerm} />
					<Command.Empty>No book found.</Command.Empty>
					<Command.Group>
						{#each books.filter(book => book.label.toLowerCase().includes(searchTerm.toLowerCase())) as book}
						<Command.Item
							value={book.value}
							onSelect={(currentValue: string) => {
								value = currentValue === value ? "" : currentValue;
								open = false;
							}}
						>
							{book.label}
						</Command.Item>
						{/each}
					</Command.Group>
				</Command.Root>
			</Popover.Content>
		</Popover.Root>
	</div>
	<div class="flex-grow text-white lm-10 mt-3 overflow-hidden lg:ml-[90px]">
		<form class="hidden md:flex bg-[#1f2e36] border border-white w-[90%] text-white text-sm item-center px-3 rounded-xl" on:submit|preventDefault={handleSearch}>
			<input type="text" placeholder="search" class="flex-1 text-white outline-none bg-transparent" bind:value={searchInput} />
			<Button variant="ghost" size="icon" class="w-[35px] h-[35px] hover:bg-[#3b5369]" on:click={handleSearch}>
				<img src="/icons/Search.svg" alt="Search" class="opacity-100" />
			</Button>
		</form>
	</div>
</div>