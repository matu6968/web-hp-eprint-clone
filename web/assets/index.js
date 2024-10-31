const uploadForm = document.getElementById('uploadForm');
const uploadStatus = document.getElementById('uploadStatus');
const commandOutput = document.getElementById('commandOutput');

uploadForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    uploadStatus.textContent = 'Uploading and printing...';
    commandOutput.textContent = '';

    const formData = new FormData(uploadForm);

    try {
        const response = await fetch('/upload', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        if (result.error) {
            throw new Error(result.error);
        } else {
            uploadStatus.textContent = result.message;
            commandOutput.textContent = result.output;
        }
    } catch (error) {
        uploadStatus.textContent = `Error: ${error.message}`;
    }
});
