let formToSubmit;
let entityToDelete;

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.delete-btn').forEach(button => {
        button.addEventListener('click', function() {
            formToSubmit = this.closest('form');
            entityToDelete = this.getAttribute('data-entity');

            const modalMessage = document.getElementById('modalMessage');
            modalMessage.innerText = `Are you sure that you want to delete this ${entityToDelete}?`;
        });
    });

    const confirmDeleteButton = document.getElementById('confirmDeleteButton');
    confirmDeleteButton.addEventListener('click', () => {
        if (formToSubmit) {
            formToSubmit.submit();
        }
    });
});
