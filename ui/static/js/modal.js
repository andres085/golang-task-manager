let formToSubmit;

function confirmDelete(button) {
    formToSubmit = button.closest('form');
    console.log(formToSubmit)
}

document.addEventListener('DOMContentLoaded', () => {
    const confirmDeleteButton = document.getElementById('confirmDeleteButton');
    confirmDeleteButton.addEventListener('click', () => {
        if (formToSubmit) {
            formToSubmit.submit();
        }
    });
});
