<script>
    import { createEventDispatcher } from 'svelte';
    import { fade } from 'svelte/transition';
    import FileIcon from './FileIcon.svelte';
    import MediaPreview from './MediaPreview.svelte';

    export let files = [];
    export let currentPath = '/';

    const dispatch = createEventDispatcher();

    let showPreview = false;
    let previewFile = null;
    let previewUrl = '';

    function navigateToDirectory(path) {
        dispatch('pathChange', { path });
    }

    function getFileUrl(file) {
        const baseApiUrl = '/api/files'; // Adjust this to your API endpoint
        const filePath = currentPath === '/'
            ? `${currentPath}${file.name}`
            : `${currentPath}/${file.name}`;

        return `${baseApiUrl}?path=${encodeURIComponent(filePath)}`;
    }

    function formatFileSize(sizeInBytes) {
        if (!sizeInBytes) return 'N/A';

        const units = ['B', 'KB', 'MB', 'GB', 'TB'];
        let size = sizeInBytes;
        let unitIndex = 0;

        while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024;
            unitIndex++;
        }

        return `${size.toFixed(1)} ${units[unitIndex]}`;
    }

    function formatDate(timestamp) {
        if (!timestamp) return 'N/A';
        return new Date(timestamp).toLocaleString();
    }

    function isPreviewable(file) {
        if (!file || file.type === 'directory') return false;

        const previewableExtensions = [
            // Images
            'jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg',
            // Videos
            'mp4', 'webm', 'ogg', 'mov',
            // Audio
            'mp3', 'wav', 'ogg', 'flac', 'm4a',
            // Documents
            'pdf',
            // Text/code
            'txt', 'md', 'js', 'css', 'html', 'json', 'xml', 'csv',
            'py', 'go', 'java', 'c', 'cpp'
        ];

        return file.extension && previewableExtensions.includes(String(file.extension).toLowerCase());
    }

    function openPreview(file) {
        if (isPreviewable(file)) {
            previewFile = file;
            previewUrl = getFileUrl(file);
            showPreview = true;
        }
    }

    function closePreview() {
        showPreview = false;
        setTimeout(() => {
            previewFile = null;
            previewUrl = '';
        }, 300);
    }
</script>

<div class="mt-4">
    {#if files.length === 0}
        <div class="p-4 text-center text-gray-500" transition:fade={{ duration: 200 }}>No files found in this directory</div>
    {:else}
        <div class="bg-white shadow rounded overflow-hidden" transition:fade={{ duration: 200 }}>
            <table class="min-w-full">
                <thead class="bg-gray-100">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Size</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Modified</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
                </thead>
                <tbody class="divide-y divide-gray-200">
                {#if currentPath !== '/'}
                    <tr class="hover:bg-gray-50">
                        <td class="px-6 py-4">
                            <button
                                    class="flex items-center text-blue-600 bg-gray-100 hover:underline"
                                    on:click={() => {
                    const parentPath = currentPath.split('/').slice(0, -1).join('/') || '/';
                    navigateToDirectory(parentPath);
                  }}
                            >
                                <FileIcon type="directory" />
                                <span class="ml-2 bg-gray-100">..</span>
                            </button>
                        </td>
                        <td class="px-6 py-4">--</td>
                        <td class="px-6 py-4">--</td>
                        <td class="px-6 py-4">--</td>
                    </tr>
                {/if}

                {#each files as file}
                    <tr class="hover:bg-gray-50">
                        <td class="px-6 py-4">
                            {#if file.type === 'directory'}
                                <button
                                        class="flex items-center text-blue-600 bg-gray-100 hover:underline"
                                        on:click={() => {
                      const newPath = currentPath === '/'
                        ? `/${file.name}`
                        : `${currentPath}/${file.name}`;
                      navigateToDirectory(newPath);
                    }}
                                >
                                    <FileIcon type="directory" />
                                    <span class="ml-2 text-blue-500 text-sm">{file.name}</span>
                                </button>
                            {:else}
                                <div class="flex items-center bg-gray-100">
                                    {#if isPreviewable(file)}
                                        <button
                                                class="flex items-center bg-gray-100 hover:text-blue-600"
                                                on:click={() => openPreview(file)}
                                        >
                                            <FileIcon type={file.extension} />
                                            <span class="ml-2 text-gray-500 text-sm">{file.name}</span>
                                        </button>
                                    {:else}
                                        <div class="flex items-center">
                                            <FileIcon type={file.extension} />
                                            <span class="ml-2 text-gray-500 text-sm">{file.name}</span>
                                        </div>
                                    {/if}
                                </div>
                            {/if}
                        </td>
                        <td class="px-6 py-4 text-gray-500 text-sm">{formatFileSize(file.size)}</td>
                        <td class="px-6 py-4 text-gray-500 text-sm">{formatDate(file.modifiedAt)}</td>
                        <td class="px-6 py-4 space-x-2">
                            {#if file.type !== 'directory'}
                                {#if isPreviewable(file)}
                                    <button
                                            on:click={() => openPreview(file)}
                                            class="text-gray-600 bg-transparent hover:text-blue-600"
                                    >
                                        <i class="fas fa-eye mr-1"></i>
                                        Preview
                                    </button>
                                {/if}
                                <a
                                        href={getFileUrl(file)}
                                        download={file.name}
                                        class="text-indigo-600 hover:text-indigo-900 ml-2"
                                >
                                    <i class="fas fa-download mr-1"></i>
                                    Download
                                </a>
                            {/if}
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

{#if showPreview}
    <MediaPreview
            file={previewFile}
            url={previewUrl}
            onClose={closePreview}
    />
{/if}