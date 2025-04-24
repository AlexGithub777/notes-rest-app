// notes.js - Vanilla JS implementation for Notes App with category support

document.addEventListener("DOMContentLoaded", function () {
    // API endpoints
    const API_URL = "/api/notes";
    const CATEGORIES_URL = "/api/notes/categories";

    // DOM elements
    const notesContainer = document.getElementById("notes-container");
    const createNoteBtn = document.getElementById("create-note-btn");
    const logoutBtn = document.getElementById("logout-btn");
    const createNoteModal = new bootstrap.Modal(
        document.getElementById("create-note-modal")
    );
    const editNoteModal = new bootstrap.Modal(
        document.getElementById("edit-note-modal")
    );
    const deleteNoteModal = new bootstrap.Modal(
        document.getElementById("delete-note-modal")
    );

    // Form elements
    const createNoteForm = document.getElementById("create-note-form");
    const editNoteForm = document.getElementById("edit-note-form");
    const deleteNoteForm = document.getElementById("delete-note-form");
    const createNoteCategorySelect = document.getElementById(
        "create-note-category"
    );
    const editNoteCategorySelect =
        document.getElementById("edit-note-category");

    // Load notes and categories when page loads
    fetchCategories().then(() => {
        fetchNotes();
    });

    // Event listeners
    if (createNoteBtn) {
        createNoteBtn.addEventListener("click", () => createNoteModal.show());
    }

    if (createNoteForm) {
        createNoteForm.addEventListener("submit", handleCreateNote);
    }

    if (editNoteForm) {
        editNoteForm.addEventListener("submit", handleEditNote);
    }

    if (deleteNoteForm) {
        deleteNoteForm.addEventListener("submit", handleDeleteNote);
    }

    if (logoutBtn) {
        logoutBtn.addEventListener("click", () => {
            window.location.href = "/logout";
        });
    }

    // Map category_id to Bootstrap badge classes
    function getCategoryBadgeClass(categoryId) {
        switch (categoryId) {
            case 1:
                return "bg-primary"; // General
            case 2:
                return "bg-success"; // Work
            case 3:
                return "bg-warning"; // Personal
            case 4:
                return "bg-secondary"; // Other
            default:
                return "bg-light text-dark border";
        }
    }

    // Fetch all categories from API
    function fetchCategories() {
        return fetch(CATEGORIES_URL)
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Failed to fetch categories");
                }
                return response.json();
            })
            .then((categories) => {
                populateCategorySelects(categories);
            })
            .catch((error) => {
                console.error("Error fetching categories:", error);
                showAlert(
                    "Error fetching categories. Please try again.",
                    "danger"
                );
            });
    }

    // Populate category select dropdowns
    function populateCategorySelects(categories) {
        // Function to populate a select element with categories
        function populateSelect(selectElement) {
            if (!selectElement) return;

            // Keep the default "Select a category" option
            selectElement.innerHTML =
                '<option value="" disabled selected>Select a category</option>';

            // Add categories from the API
            categories.forEach((category) => {
                // convert to integer
                category.id = parseInt(category.id, 10);
                const option = document.createElement("option");
                option.value = category.id;
                option.textContent = category.name;
                selectElement.appendChild(option);
            });
        }

        // Populate both create and edit form selects
        populateSelect(createNoteCategorySelect);
        populateSelect(editNoteCategorySelect);
    }

    // Fetch all notes from API
    function fetchNotes() {
        fetch(API_URL)
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Failed to fetch notes");
                }
                return response.json();
            })
            .then((notes) => {
                renderNotes(notes);
            })
            .catch((error) => {
                console.error("Error fetching notes:", error);
                showAlert("Error fetching notes. Please try again.", "danger");
            });
    }

    // Render notes as Bootstrap cards
    function renderNotes(notes) {
        if (!notesContainer) return;

        notesContainer.innerHTML = "";

        if (notes.length === 0) {
            notesContainer.innerHTML = `
          <div class="col-12 text-center py-5">
            <p class="text-muted">No notes found. Create your first note!</p>
          </div>
        `;
            return;
        }

        notes.forEach((note) => {
            console.log("Note:", note); // Debugging line
            const noteCard = document.createElement("div");
            noteCard.className = "col-md-4 mb-4";
            noteCard.innerHTML = `
            <div class="card h-100">
            <div class="card-header pb-2">
              <div class="d-flex justify-content-between align-items-start">
                <h5 class="card-title mb-0">${escapeHtml(note.title)}</h5>
                <span class="badge ${getCategoryBadgeClass(
                    note.category_id
                )}">${escapeHtml(note.category_name)}</span>
              </div>
            </div>
            <div class="card-body">
              <p class="card-text">${escapeHtml(note.content)}</p>
            </div>
            <div class="card-footer bg-transparent">
              <div class="d-flex justify-content-between align-items-center">
                <div>
                  <small class="text-muted">${formatDate(
                      note.created_at
                  )}</small>
                  <span class="badge bg-info text-dark ms-2">@${escapeHtml(
                      note.username
                  )}</span>
                </div>
                <div>
                  <button class="btn btn-sm btn-outline-primary me-2 edit-note-btn" 
                    data-id="${note.id}" 
                    data-title="${escapeHtml(note.title)}" 
                    data-content="${escapeHtml(note.content)}"
                    data-category="${note.category_id}"
                    aria-label="Edit note">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-edit">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                    </svg>
                  </button>
                  <button class="btn btn-sm btn-outline-danger delete-note-btn" 
                    data-id="${note.id}" 
                    data-title="${escapeHtml(note.title)}"
                    aria-label="Delete note">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-trash-2">
                      <polyline points="3 6 5 6 21 6"></polyline>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                      <line x1="10" y1="11" x2="10" y2="17"></line>
                      <line x1="14" y1="11" x2="14" y2="17"></line>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
          
        `;

            notesContainer.appendChild(noteCard);
        });

        // Add event listeners to edit and delete buttons
        document.querySelectorAll(".edit-note-btn").forEach((btn) => {
            btn.addEventListener("click", prepareEditNote);
        });

        document.querySelectorAll(".delete-note-btn").forEach((btn) => {
            btn.addEventListener("click", prepareDeleteNote);
        });
    }

    // Prepare edit note modal
    function prepareEditNote(event) {
        const button = event.currentTarget;
        const id = button.getAttribute("data-id");
        const title = button.getAttribute("data-title");
        const content = button.getAttribute("data-content");
        const categoryId = button.getAttribute("data-category");

        // set the form values
        document.getElementById("edit-note-id").value = id;
        document.getElementById("edit-note-title").value = title;
        document.getElementById("edit-note-content").value = content;
        document.getElementById("edit-note-category").value = categoryId;

        console.log("Category ID:", categoryId);
        console.log("Selected Category:", editNoteCategorySelect.value);

        // Set the category if it exists
        if (editNoteCategorySelect && categoryId) {
            editNoteCategorySelect.value = categoryId;
        }

        editNoteModal.show();
    }

    // Prepare delete note modal
    function prepareDeleteNote(event) {
        const button = event.currentTarget;
        const id = button.getAttribute("data-id");
        const title = button.getAttribute("data-title");

        document.getElementById("delete-note-id").value = id;
        document.getElementById("delete-note-title").textContent = title;

        deleteNoteModal.show();
    }

    // Handle create note form submission
    function handleCreateNote(event) {
        event.preventDefault();

        const title = document.getElementById("create-note-title").value.trim();
        const content = document
            .getElementById("create-note-content")
            .value.trim();
        const categoryId = createNoteCategorySelect
            ? createNoteCategorySelect.value
            : null;

        if (!title || !content) {
            showAlert("Please fill in all required fields", "warning");
            return;
        }

        const noteData = {
            title,
            content,
        };

        // Add category_id if selected
        if (categoryId) {
            noteData.category_id = parseInt(categoryId, 10);
        }

        fetch(API_URL, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(noteData),
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Failed to create note");
                }
                return response.json();
            })
            .then(() => {
                createNoteModal.hide();
                createNoteForm.reset();
                showAlert("Note created successfully!", "success");
                fetchNotes();
            })
            .catch((error) => {
                console.error("Error creating note:", error);
                showAlert("Error creating note. Please try again.", "danger");
            });
    }

    // Handle edit note form submission
    function handleEditNote(event) {
        event.preventDefault();

        const id = document.getElementById("edit-note-id").value;
        const title = document.getElementById("edit-note-title").value.trim();
        const content = document
            .getElementById("edit-note-content")
            .value.trim();
        const categoryId = editNoteCategorySelect
            ? editNoteCategorySelect.value
            : null;

        if (!title || !content) {
            showAlert("Please fill in all required fields", "warning");
            return;
        }

        const noteData = {
            title,
            content,
        };

        // Add category_id if selected
        if (categoryId) {
            noteData.category_id = parseInt(categoryId, 10);
        }

        fetch(`${API_URL}/${id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(noteData),
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Failed to update note");
                }
                return response.json();
            })
            .then(() => {
                editNoteModal.hide();
                showAlert("Note updated successfully!", "success");
                fetchNotes();
            })
            .catch((error) => {
                console.error("Error updating note:", error);
                showAlert("Error updating note. Please try again.", "danger");
            });
    }

    function handleDeleteNote(event) {
        event.preventDefault();

        const id = document.getElementById("delete-note-id").value;

        fetch(`${API_URL}/${id}`, {
            method: "DELETE",
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Failed to delete note");
                }
                // Don't try to parse JSON when status is 204
                deleteNoteModal.hide();
                showAlert("Note deleted successfully!", "success");
                fetchNotes(); // Reloads the page content (via JS)
            })
            .catch((error) => {
                console.error("Error deleting note:", error);
                showAlert("Error deleting note. Please try again.", "danger");
            });
    }

    // Helper function to display alerts
    function showAlert(message, type = "info") {
        const alertContainer = document.getElementById("alert-container");
        if (!alertContainer) return;

        const alert = document.createElement("div");
        alert.className = `alert alert-${type} alert-dismissible fade show`;
        alert.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
      `;

        alertContainer.appendChild(alert);

        // Auto dismiss after 3 seconds
        setTimeout(() => {
            const bsAlert = new bootstrap.Alert(alert);
            bsAlert.close();
        }, 3000);
    }

    // Helper function to escape HTML to prevent XSS
    function escapeHtml(unsafe) {
        if (!unsafe) return "";
        return unsafe
            .toString()
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    // Helper function to format date
    function formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString();
    }
});
