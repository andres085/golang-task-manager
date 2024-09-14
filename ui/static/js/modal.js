let formToSubmit;
let entityToDelete;

function confirmDelete(button, entity) {
    formToSubmit = button.closest('form');
    entityToDelete = entity;
    const modalMessage = document.getElementById('modalMessage');
    modalMessage.innerText = `Are you sure that you want to delete this ${entityToDelete}?`
}

document.addEventListener('DOMContentLoaded', () => {
    const confirmDeleteButton = document.getElementById('confirmDeleteButton');
    confirmDeleteButton.addEventListener('click', () => {
        if (formToSubmit) {
            formToSubmit.submit();
        }
    });
});
