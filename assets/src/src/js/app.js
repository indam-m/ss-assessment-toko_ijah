function updateItem(e) {
    var editForm = document.getElementById('update-form');
    // display update form
    editForm.style.display = 'block';
    var id = e.parentNode.parentNode.childNodes[3].innerText;
    document.getElementById('id-to-update').value = id;
}

function deleteItem(e) {
    var deleteForm = document.getElementById('delete-form');
    deleteForm.style.display = 'block';
    var id = e.parentNode.parentNode.childNodes[3].innerText;
    document.getElementById('id-to-delete').value = id;
}
