<script>
    export let type = 'file';

    function getIconClass() {
        switch (type) {
            case 'directory':
                return 'fas fa-folder text-yellow-400';
            case 'pdf':
                return 'fas fa-file-pdf text-red-500';
            case 'doc':
            case 'docx':
                return 'fas fa-file-word text-blue-500';
            case 'xls':
            case 'xlsx':
                return 'fas fa-file-excel text-green-500';
            case 'ppt':
            case 'pptx':
                return 'fas fa-file-powerpoint text-orange-500';
            case 'jpg':
            case 'jpeg':
            case 'png':
            case 'gif':
                return 'fas fa-file-image text-purple-500';
            case 'mp3':
            case 'wav':
                return 'fas fa-file-audio text-blue-400';
            case 'mp4':
            case 'avi':
            case 'mkv':
                return 'fas fa-file-video text-pink-500';
            case 'zip':
            case 'rar':
            case 'tar':
            case 'gz':
                return 'fas fa-file-archive text-amber-500';
            case 'txt':
                return 'fas fa-file-alt text-gray-500';
            case 'html':
            case 'css':
            case 'js':
                return 'fas fa-file-code text-teal-500';
            default:
                return 'fas fa-file text-gray-400';
        }
    }
</script>

<i class={getIconClass()}></i>