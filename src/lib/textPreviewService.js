/**
 * Fetches and displays the content of a text file
 * @param {string} url - URL to the text file
 * @param {string} elementId - ID of the element to insert the text into
 * @param {number} maxLength - Maximum length of text to display (in characters)
 */
export async function loadTextPreview(url, elementId = 'text-content', maxLength = 100000) {
    try {
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`Failed to load text: ${response.status} ${response.statusText}`);
        }

        let text = await response.text();

        const isTruncated = text.length > maxLength;
        if (isTruncated) {
            text = text.substring(0, maxLength) + '\n\n[File truncated due to size. Download the file to see the full content.]';
        }

        // Find the element and inject the text
        const element = document.getElementById(elementId);
        if (element) {
            element.textContent = text;

            if (isTruncated) {
                const notice = document.createElement('div');
                notice.className = 'bg-yellow-100 p-2 mt-4 rounded text-yellow-800 text-sm';
                notice.textContent = 'This file is too large to preview completely. Download the file to see all content.';
                element.parentNode.insertBefore(notice, element.nextSibling);
            }
        }
    } catch (error) {
        console.error('Error loading text preview:', error);

        const element = document.getElementById(elementId);
        if (element) {
            element.innerHTML = `<div class="text-red-500 p-4">
        <i class="fas fa-exclamation-triangle mr-2"></i>
        Error loading text preview: ${error.message}
      </div>`;
        }
    }
}