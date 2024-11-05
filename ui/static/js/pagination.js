document.addEventListener('DOMContentLoaded', function () {
  const limitSelect = document.getElementById('limitSelect');
  if (limitSelect) {
    limitSelect.addEventListener('change', function () {
      document.getElementById('limitForm').submit();
    });
  }
});
