<script lang="ts">
    import { fade } from 'svelte/transition';
    import { onMount, onDestroy } from 'svelte';
    import { loadTextPreview } from '../lib/textPreviewService';

    export let file = null;
    export let url = '';
    export let onClose = () => {};

    $: fileType = getFileType(file?.extension);

    function getFileType(extension) {
        if (!extension) return 'unknown';

        const lowerExt = String(extension).toLowerCase();

        // Image types
        if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg'].includes(lowerExt)) {
            return 'image';
        }

        // Video types
        if (['mp4', 'webm', 'ogg', 'mov', 'avi', 'mkv'].includes(lowerExt)) {
            return 'video';
        }

        // Audio types
        if (['mp3', 'wav', 'ogg', 'flac', 'm4a'].includes(lowerExt)) {
            return 'audio';
        }

        // PDF
        if (lowerExt === 'pdf') {
            return 'pdf';
        }

        // Text/code
        if (['txt', 'md', 'js', 'css', 'html', 'json', 'xml', 'csv', 'py', 'go', 'java', 'c', 'cpp'].includes(lowerExt)) {
            return 'text';
        }

        return 'unknown';
    }

    function handleKeydown(event) {
        if (event.key === 'Escape') {
            onClose();
        }
    }

    onMount(() => {
        document.addEventListener('keydown', handleKeydown);

        if (fileType === 'text' && url) {
            loadTextPreview(url);
        }
    });

    onDestroy(() => {
        document.removeEventListener('keydown', handleKeydown);
    });
</script>

<svelte:window on:keydown={handleKeydown}/>

<div class="fixed inset-0 bg-black bg-opacity-80 z-50 flex items-center justify-center p-4"
     transition:fade={{ duration: 200 }}
     on:click={onClose} role="none">

    <div class="relative bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] flex flex-col"
         on:click|stopPropagation role="none">

        <!-- Header -->
        <div class="flex items-center justify-between p-4 border-b">
            <h3 class="text-xl font-semibold text-gray-800 truncate">
                {file?.name || 'Preview'}
                {#if fileType !== 'unknown'}
                    <span class="text-sm text-gray-500 ml-2">{fileType}</span>
                {/if}
            </h3>
            <button class="text-gray-500 bg-gray-100 hover:text-gray-700" on:click={onClose} aria-label="close">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </button>
        </div>

        <!-- Content -->
        <div class="flex-1 overflow-auto p-4 flex items-center justify-center bg-gray-100">
            {#if fileType === 'image'}
                <img
                        src={url}
                        alt={file?.name || 'Image preview'}
                        class="max-w-full max-h-[70vh] object-contain"
                />
            {:else if fileType === 'video'}
                <video
                        src={url}
                        controls
                        autoplay
                        class="max-w-full max-h-[70vh]"
                        controlsList="nodownload"
                >
                    Your browser does not support the video tag.
                </video>
            {:else if fileType === 'audio'}
                <div class="w-full max-w-lg bg-white p-6 rounded-lg shadow">
                    <div class="text-center mb-4">
                        <i class="fas fa-music text-5xl text-blue-500"></i>
                    </div>
                    <audio
                            src={url}
                            controls
                            autoplay
                            class="w-full"
                            controlsList="nodownload"
                    >
                        Your browser does not support the audio tag.
                    </audio>
                    <div class="text-center mt-2 text-gray-600 text-sm">
                        {file?.name}
                    </div>
                </div>
            {:else if fileType === 'pdf'}
                <iframe
                        src={`${url}#toolbar=0`}
                        class="w-full h-[70vh]"
                        title={file?.name || 'PDF preview'}
                ></iframe>
            {:else if fileType === 'text'}
                <div class="w-full h-[70vh] bg-white p-4 rounded shadow overflow-auto">
                    <div class="flex justify-between items-center mb-2 pb-2 border-b">
                        <span class="text-sm text-gray-500">Text Preview</span>
                        <span class="text-xs bg-gray-200 px-2 py-1 rounded">{file?.extension?.toUpperCase()}</span>
                    </div>
                    <div id="text-content" class="font-mono text-sm whitespace-pre-wrap">
                        Loading text content...
                    </div>
                </div>
            {:else}
                <div class="text-center">
                    <i class="fas fa-file text-4xl text-gray-400 mb-2"></i>
                    <p class="text-lg text-gray-600">Preview not available for this file type</p>
                    <p class="text-sm text-gray-500 mt-1">({file?.extension || 'unknown'} format)</p>
                </div>
            {/if}
        </div>

        <!-- Footer -->
        <div class="border-t p-4 flex justify-between items-center">
            <div class="text-sm text-gray-600">
                {#if file?.size}
                    <span>{formatFileSize(file.size)}</span>
                {/if}
                {#if file?.modifiedAt}
                    <span class="ml-4">Modified: {new Date(file.modifiedAt).toLocaleString()}</span>
                {/if}
            </div>
            <a
                    href={url}
                    download={file?.name}
                    class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg flex items-center"
            >
                <i class="fas fa-download mr-2"></i>
                Download
            </a>
        </div>
    </div>
</div>

<script context="module">
    export function formatFileSize(sizeInBytes) {
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
</script>