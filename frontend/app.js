const API_BASE = '/api/dogs';
let dogs = [];
let dogToDelete = null;

// Load dogs on page load
document.addEventListener('DOMContentLoaded', function() {
    loadDogs();
    
    // Add button event listener
    document.getElementById('addBtn').addEventListener('click', openAddModal);
    
    // Cancel button event listener
    document.getElementById('cancelBtn').addEventListener('click', closeModal);
    
    // Cancel delete button event listener
    document.getElementById('cancelDeleteBtn').addEventListener('click', closeDeleteModal);
    
    // Event delegation for edit and delete buttons
    document.getElementById('dogsContainer').addEventListener('click', function(e) {
        if (e.target.classList.contains('edit-btn')) {
            const id = parseInt(e.target.dataset.id);
            openEditModal(id);
        } else if (e.target.classList.contains('delete-btn')) {
            const id = parseInt(e.target.dataset.id);
            const breed = e.target.dataset.breed;
            openDeleteModal(id, breed);
        }
    });
});

// Form submission
document.getElementById('dogForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const dogId = document.getElementById('dogId').value;
    const breed = document.getElementById('breed').value.trim();
    const subBreeds = document.getElementById('subBreeds').value.trim();

    const data = {
        breed: breed,
        subBreeds: subBreeds || null
    };

    try {
        let response;
        if (dogId) {
            response = await fetch(`${API_BASE}/${dogId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
        } else {
            response = await fetch(API_BASE, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
        }

        const result = await response.json();
        if (result.success) {
            showNotification(result.message, 'success');
            closeModal();
            loadDogs();
        } else {
            showNotification(result.message, 'error');
        }
    } catch (error) {
        showNotification('Error saving dog', 'error');
    }
});

// Delete confirmation
document.getElementById('confirmDeleteBtn').addEventListener('click', async function() {
    if (dogToDelete) {
        try {
            const response = await fetch(`${API_BASE}/${dogToDelete}`, {
                method: 'DELETE'
            });
            const result = await response.json();
            if (result.success) {
                showNotification(result.message, 'success');
                closeDeleteModal();
                loadDogs();
            } else {
                showNotification(result.message, 'error');
            }
        } catch (error) {
            showNotification('Error deleting dog', 'error');
        }
    }
});

async function loadDogs() {
    try {
        const response = await fetch(API_BASE);
        const result = await response.json();
        if (result.success) {
            dogs = result.data;
            renderDogs(dogs);
            updateStats(dogs);
        }
    } catch (error) {
        showNotification('Error loading dogs', 'error');
        document.getElementById('dogsContainer').innerHTML = '<div class="error">Error loading dogs</div>';
    }
}

function renderDogs(dogsToRender) {
    const container = document.getElementById('dogsContainer');
    
    if (dogsToRender.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <p>No dogs found</p>
            </div>
        `;
        return;
    }

    container.innerHTML = dogsToRender.map(dog => `
        <div class="dog-item">
            <div class="dog-info">
                <div class="breed">${dog.breed}</div>
                ${dog.subBreeds ? `<div class="sub-breeds">${dog.subBreeds}</div>` : ''}
                ${dog.updatedAt ? `<div class="timestamp">Updated: ${formatDate(dog.updatedAt)}</div>` : ''}
            </div>
            <div class="dog-actions">
                <button class="btn btn-secondary edit-btn" data-id="${dog.id}">Edit</button>
                <button class="btn btn-danger delete-btn" data-id="${dog.id}" data-breed="${dog.breed}">Delete</button>
            </div>
        </div>
    `).join('');
}

function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString();
}

function updateStats(dogsList) {
    // Stats removed for minimalist design
}

function openAddModal() {
    document.getElementById('modalTitle').textContent = 'Add Dog';
    document.getElementById('dogId').value = '';
    document.getElementById('breed').value = '';
    document.getElementById('subBreeds').value = '';
    document.getElementById('dogModal').classList.add('active');
}

function openEditModal(id) {
    const dog = dogs.find(d => d.id === id);
    if (dog) {
        document.getElementById('modalTitle').textContent = 'Edit Dog';
        document.getElementById('dogId').value = dog.id;
        document.getElementById('breed').value = dog.breed;
        document.getElementById('subBreeds').value = dog.subBreeds || '';
        document.getElementById('dogModal').classList.add('active');
    }
}

function closeModal() {
    document.getElementById('dogModal').classList.remove('active');
}

function openDeleteModal(id, breed) {
    dogToDelete = id;
    document.getElementById('deleteDogName').textContent = breed;
    document.getElementById('deleteModal').classList.add('active');
}

function closeDeleteModal() {
    document.getElementById('deleteModal').classList.remove('active');
    dogToDelete = null;
}

function showNotification(message, type) {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `notification ${type} active`;
    
    setTimeout(() => {
        notification.classList.remove('active');
    }, 3000);
}

