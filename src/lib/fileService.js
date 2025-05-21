export async function fetchDirectoryContents(path) {
    try {
        const response = await fetch(`/api/files/list?path=${encodeURIComponent(path)}`);

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        return data.files;
    } catch (error) {
        console.error('Error fetching directory contents:', error);
        throw error;
    }
}