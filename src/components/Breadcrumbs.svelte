<script>
    import { createEventDispatcher } from 'svelte';

    export let path = '/';

    const dispatch = createEventDispatcher();

    function navigateTo(pathSegment) {
        dispatch('pathChange', { path: pathSegment });
    }

    $: pathSegments = path.split('/').filter(segment => segment);
    $: breadcrumbs = pathSegments.map((segment, index) => {
        const path = '/' + pathSegments.slice(0, index + 1).join('/');
        return { name: segment, path };
    });
</script>

<div class="flex items-center space-x-2 text-gray-700 mb-4">
    <button
            class="hover:text-blue-600 bg-gray-100 font-medium"
            on:click={() => navigateTo('/')}
    >
        Home
    </button>

    {#each breadcrumbs as crumb, i}
        <span>/</span>
        <button
                class="hover:text-blue-600 bg-gray-100 font-medium"
                on:click={() => navigateTo(crumb.path)}
        >
            {crumb.name}
        </button>
    {/each}
</div>