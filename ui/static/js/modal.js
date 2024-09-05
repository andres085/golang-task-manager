let formToSubmit;

function confirmDelete(button) {
    formToSubmit = button.closest('form');
}

document.addEventListener('DOMContentLoaded', () => {
    const confirmDeleteButton = document.getElementById('confirmDeleteButton');
    confirmDeleteButton.addEventListener('click', () => {
        if (formToSubmit) {
            formToSubmit.submit();
        }
    });
});
