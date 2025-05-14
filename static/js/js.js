document.addEventListener("DOMContentLoaded", function () {
    // API endpoints
    const API_URL = "/api/notes";
    const CATEGORIES_URL = "/api/notes/categories";

    // check if page is all notes
    const isAllNotesPage = window.location.pathname === "/all-notes";
    if (isAllNotesPage) {
        document.title = "All Notes";
        // add event listener to my notes button
        const mynotesbtn = document.getElementById("my-notes-btn");
        if (mynotesbtn) {
            mynotesbtn.addEventListener("click", () => {
                window.location.href = "/home";
            });
        }

        function fetchAllNotes(query = "", categoryId = null) {
            let url = `${API_URL}/all`;
            let params = [];

            if (query) {
                params.push(`search=${encodeURIComponent(query)}`);
            }

            if (categoryId && categoryId !== "all") {
                params.push(`category-id=${encodeURIComponent(categoryId)}`);
            }

            if (params.length > 0) {
                url += `?${params.join("&")}`;
            }

            fetch(url)
                .then((response) => {
                    if (!response.ok) {
                        throw new Error("Failed to fetch notes");
                    }
                    return response.json();
                })
                .then((notes) => {
                    if (!notes || notes.length === 0) {
                        showAlert("No notes found for your search.", "warning");

                        // Clear the notes container
                        notesContainer.innerHTML = "";

                        return;
                    }
                    renderNotes(notes);
                })
                .catch((error) => {
                    console.error("Error fetching notes:", error);
                    showAlert(
                        "Error fetching notes. Please try again.",
                        "danger"
                    );
                });
        }
        // Fetch all notes on page load
        fetchAllNotes();
    } else {
        document.title = "My Notes";
    }

    // DOM elements
    const notesContainer = document.getElementById("notes-container");
    const createNoteBtn = document.getElementById("create-note-btn");
    const allNotesBtn = document.getElementById("all-notes-btn");
    const logoutBtn = document.getElementById("logout-btn");
    const logoutBtnAllNotes = document.getElementById("logout-btn-all-notes");
    const searchInputMyNotes = document.getElementById("search-input-my-notes");
    let searchQueryMyNotes = "";
    const searchInputAllNotes = document.getElementById(
        "search-input-all-notes"
    );
    let searchQueryAllNotes = "";

    const createNoteModal = document.getElementById("create-note-modal")
        ? new bootstrap.Modal(document.getElementById("create-note-modal"))
        : null;
    const editNoteModal = document.getElementById("edit-note-modal")
        ? new bootstrap.Modal(document.getElementById("edit-note-modal"))
        : null;
    const deleteNoteModal = document.getElementById("delete-note-modal")
        ? new bootstrap.Modal(document.getElementById("delete-note-modal"))
        : null;

    // Form elements
    const createNoteForm = document.getElementById("create-note-form");
    const editNoteForm = document.getElementById("edit-note-form");
    const deleteNoteForm = document.getElementById("delete-note-form");
    const createNoteCategorySelect = document.getElementById(
        "create-note-category"
    );
    const editNoteCategorySelect =
        document.getElementById("edit-note-category");

    const categoryFilterSelect = document.getElementById(
        "category-filter-my-notes"
    );
    const categoryFilterSelectAllNotes = document.getElementById(
        "category-filter-all-notes"
    );

    // Load notes and categories when page loads
    fetchCategories().then(() => {
        if (!isAllNotesPage) {
            fetchNotes(); // only fetch user notes on /home
        }
    });

    if (categoryFilterSelect) {
        categoryFilterSelect.addEventListener("change", filterByCategory);
    }

    if (categoryFilterSelectAllNotes) {
        categoryFilterSelectAllNotes.addEventListener(
            "change",
            filterByCategoryAllNotes
        );
    }

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

    if (logoutBtnAllNotes) {
        logoutBtnAllNotes.addEventListener("click", () => {
            window.location.href = "/logout";
        });
    }

    if (allNotesBtn) {
        allNotesBtn.addEventListener("click", () => {
            window.location.href = "/all-notes";
        });
    }

    function handleNoteSearchChange() {
        searchQueryMyNotes = searchInputMyNotes.value.trim();
        selectedCategoryId = categoryFilterSelect.value;
        fetchNotes(searchQueryMyNotes, selectedCategoryId);
    }

    if (searchInputMyNotes) {
        searchInputMyNotes.addEventListener(
            "input",
            debounce(handleNoteSearchChange, 300)
        );
    }

    function handleNoteSearchChangeAllNotes() {
        searchQueryAllNotes = searchInputAllNotes.value.trim();
        selectedCategoryId = categoryFilterSelectAllNotes.value;
        fetchAllNotes(searchQueryAllNotes, selectedCategoryId);
    }

    if (searchInputAllNotes) {
        searchInputAllNotes.addEventListener(
            "input",
            debounce(handleNoteSearchChangeAllNotes, 300)
        );
    }

    // debounce function – stops function from running too often (e.g. while typing)
    function debounce(fn, delay) {
        let timer; // store timeout so we can clear it

        return function (...args) {
            clearTimeout(timer); // clear previous timer if user types again

            // wait `delay` ms before running the actual function
            timer = setTimeout(() => {
                fn.apply(this, args); // run fn with same args and context
            }, delay);
        };
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

            // If selectElement is categoryFilterSelect or categoryFilterSelectAllNotes, add "all" option
            if (
                selectElement === categoryFilterSelect ||
                selectElement === categoryFilterSelectAllNotes
            ) {
                const allOption = document.createElement("option");
                allOption.value = "all";
                allOption.textContent = "All Categories";
                selectElement.appendChild(allOption);
            } else {
                // Set the default value to "Select a category"
                const defaultOption = document.createElement("option");
                defaultOption.value = "";
                defaultOption.textContent = "Select a category";
                selectElement.appendChild(defaultOption);
            }

            // Add categories from the API
            categories.forEach((category) => {
                category.id = parseInt(category.id, 10); // Ensure ID is an integer
                const option = document.createElement("option");
                option.value = category.id;
                option.textContent = category.name;
                selectElement.appendChild(option);
            });
        }

        // Populate select elements
        populateSelect(createNoteCategorySelect);
        populateSelect(editNoteCategorySelect);
        populateSelect(categoryFilterSelect);
        populateSelect(categoryFilterSelectAllNotes);
    }

    // Fetch all notes from API
    function fetchNotes(query = "", categoryId = null) {
        let url = API_URL;
        let params = [];

        if (query) {
            params.push(`search=${encodeURIComponent(query)}`);
        }

        if (categoryId && categoryId !== "all") {
            params.push(`category-id=${encodeURIComponent(categoryId)}`);
        }

        if (params.length > 0) {
            url += `?${params.join("&")}`;
        }

        fetch(url)
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Failed to fetch notes");
                }
                return response.json();
            })
            .then((notes) => {
                if (!notes || notes.length === 0) {
                    showAlert("No notes found for your search.", "warning");

                    // Clear the notes container
                    notesContainer.innerHTML = "";

                    return;
                }
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

        // Check current page
        const currentPath = window.location.pathname;
        const isHomePage = currentPath === "/home";
        const showEditDeleteButtons = isHomePage;

        notes.forEach((note) => {
            const noteCard = document.createElement("div");
            noteCard.className = "col-md-4 mb-4";

            // Generate action buttons HTML only if we're on the home page
            const actionButtons = showEditDeleteButtons
                ? `
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
        `
                : "";

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
                ${actionButtons}
              </div>
            </div>
          </div>
        `;

            notesContainer.appendChild(noteCard);
        });

        // Only add event listeners if buttons are present (on home page)
        if (showEditDeleteButtons) {
            // Add event listeners to edit and delete buttons
            document.querySelectorAll(".edit-note-btn").forEach((btn) => {
                btn.addEventListener("click", prepareEditNote);
            });

            document.querySelectorAll(".delete-note-btn").forEach((btn) => {
                btn.addEventListener("click", prepareDeleteNote);
            });
        }
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

                deleteNoteModal.hide();
                showAlert("Note deleted successfully!", "success");
                fetchNotes();
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

    function filterByCategory() {
        var categoryId = document.getElementById(
            "category-filter-my-notes"
        ).value;

        // if categoryId is "all", dont add it to the url
        if (categoryId === "all") {
            categoryId = null;
        }

        const url = categoryId
            ? `/api/notes?category-id=${categoryId}`
            : `/api/notes`;

        fetch(url)
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
                console.error("Error:", error);
                showAlert("No notes found for this category.", "warning");
            });
    }

    function filterByCategoryAllNotes() {
        var categoryId = document.getElementById(
            "category-filter-all-notes"
        ).value;

        // if categoryId is "all", dont add it to the url
        if (categoryId === "all") {
            categoryId = null;
        }

        const url = categoryId
            ? `/api/notes/all?category-id=${categoryId}`
            : `/api/notes/all`;

        fetch(url)
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
                console.error("Error:", error);
                showAlert("No notes found for this category.", "warning");
            });
    }
});
