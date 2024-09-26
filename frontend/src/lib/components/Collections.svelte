<script lang="ts">
    import { onMount } from 'svelte';
    import { GetCollections } from '$lib/wailsjs/go/main/App';
    import Card from '$lib/components/ui/card/card.svelte';
    import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "$lib/components/ui/dialog";
    import { cn } from '$lib/utils';
    import type { Collection } from '$lib/wailsjs/go/models';

    let collections: Collection[] = [];
    let selectedCollection: Collection | null = null;
    let dialogOpen = false;

    onMount(async () => {
        try {
            collections = await GetCollections();
            console.log('Collections:', collections);
        } catch (error) {
            console.error('Failed to fetch collections:', error);
        }
    });

    function openDialog(collection: Collection) {
        selectedCollection = collection;
        dialogOpen = true;
    }
</script>

<div class="p-6 bg-[#1e3040] min-h-screen">
    <h1 class="text-2xl font-bold mb-6 text-white">Collections</h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {#each collections as collection}
            <Card 
                class={cn("bg-[#243a4d] text-white p-4 hover:bg-[#2a4358] transition-colors cursor-pointer relative")}
                on:click={() => openDialog(collection)}
            >
                <h2 class="text-lg font-semibold mb-2">{collection.name}</h2>
                <p class="text-sm text-gray-300 absolute bottom-2 right-4">
                    {collection.num_documents} documents
                </p>
            </Card>
        {/each}
    </div>
</div>

<Dialog bind:open={dialogOpen}>
        <DialogContent class="bg-[#243a4d] text-white border-gray-600">
            <DialogHeader>
                <DialogTitle class="text-white">{selectedCollection?.name}</DialogTitle>
                <DialogDescription class="text-gray-300">
                    This collection contains {selectedCollection?.num_documents} documents.
                </DialogDescription>
        </DialogHeader>
            <div>
                <p>Additional details about the collection can be displayed here.</p>
            </div>
        </DialogContent>
</Dialog>