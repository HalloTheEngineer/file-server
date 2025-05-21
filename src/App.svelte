<script lang="ts">
    import {onMount} from 'svelte';
    import FileList from "./components/FileList.svelte"
    import Breadcrumbs from './components/Breadcrumbs.svelte';
    import {fetchDirectoryContents} from "./lib/fileService.js"
    import {fade} from "svelte/transition";

    let currentPath = '/';
    let files: any[] = [];
    let loading = true;
    let error: string | null = null;

    onMount(async () => {
        await loadDirectory(currentPath);
    });

    async function loadDirectory(path: string) {
        loading = true;
        error = null;

        try {
            files = await fetchDirectoryContents(path);
            currentPath = path;
        } catch (err: any) {
            error = `Error loading directory: ${err.message}`;
            console.error(err);
        }

        loading = false;
    }

    function handlePathChange(event) {
        const path = event.detail.path;
        loadDirectory(path);
    }
</script>

<main class="container mx-auto p-4">
    <h1 class="text-3xl font-bold mb-6">File Server</h1>

    <Breadcrumbs path={currentPath} on:pathChange={handlePathChange}/>

    {#if loading}
        <div class="mt-4 p-4 text-center" transition:fade={{ duration: 150 }}>Loading...</div>
    {:else if error}
        <div class="mt-4 p-4 bg-red-100 text-red-700 rounded" transition:fade={{ duration: 200 }}>{error}</div>
    {:else}
        <FileList
                files={files}
                currentPath={currentPath}
                on:pathChange={handlePathChange}
        />
    {/if}
</main>

<style>
    main {
        max-width: 1200px;
    }
</style>