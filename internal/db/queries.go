package db

import (
	"database/sql"
	"time"

	"github.com/AlexGithub777/notes-rest-app/internal/models"
)

func GetAllCategories() ([]models.Category, error) {
	rows, err := DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetAllNotes returns all notes for the logged-in user
func GetNotesForUser(userID int) ([]models.Note, error) {
	rows, err := DB.Query(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		WHERE notes.user_id = $1
		ORDER BY notes.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// GetNotesForUserByCategory returns all notes for the logged-in user filtered by category
func GetNotesForUserByCategory(userID, categoryID int) ([]models.Note, error) {
	rows, err := DB.Query(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		WHERE notes.user_id = $1 AND notes.category = $2
		ORDER BY notes.created_at DESC
	`, userID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// GetNotesForUserBySearch returns all notes for the logged-in user filtered by search term
func GetNotesForUserBySearch(userID int, search string) ([]models.Note, error) {
	rows, err := DB.Query(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		WHERE notes.user_id = $1 AND (notes.title ILIKE '%' || $2 || '%' OR notes.content ILIKE '%' || $2 || '%')
		ORDER BY notes.created_at DESC
	`, userID, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// GetNotesForUserByCategoryAndSearch returns all notes for the logged-in user filtered by category and search term
func GetNotesForUserByCategoryAndSearch(userID, categoryID int, search string) ([]models.Note, error) {
	rows, err := DB.Query(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		WHERE notes.user_id = $1 AND notes.category = $2 AND (notes.title ILIKE '%' || $3 || '%' OR notes.content ILIKE '%' || $3 || '%')
		ORDER BY notes.created_at DESC
	`, userID, categoryID, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// GetNoteByID returns a specific note for the logged-in user
func GetNoteByID(userID, noteID int) (models.Note, error) {
	var note models.Note
	err := DB.QueryRow(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		WHERE notes.id = $1 AND notes.user_id = $2
	`, noteID, userID).
		Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username)
	return note, err
}

func NoteExists(id int) (bool, error) {
	var exists int
	err := DB.QueryRow("SELECT id FROM notes WHERE id = $1", id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

// CreateNote creates a new note and associates it with the logged-in user, including the category.
func CreateNote(userID int, title, content string, categoryID int, now time.Time) (models.Note, error) {
	var note models.Note
	err := DB.QueryRow(
		"INSERT INTO notes (title, content, created_at, updated_at, user_id, category) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, content, category, user_id, created_at, updated_at",
		title, content, now, now, userID, categoryID,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CategoryID, &note.UserID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return note, err
	}
	// Fetch category name
	categoryName, err := getCategoryNameByID(categoryID)
	if err != nil {
		return note, err
	}
	note.CategoryName = categoryName
	return note, nil
}

// UpdateNote updates an existing note, including the category, and ensures it belongs to the logged-in user.
func UpdateNote(userID, id int, title, content string, categoryID int, now time.Time) (models.Note, error) {
	var note models.Note
	_, err := DB.Exec(
		"UPDATE notes SET title = $1, content = $2, updated_at = $3, category = $4 WHERE id = $5 AND user_id = $6",
		title, content, now, categoryID, id, userID,
	)
	if err != nil {
		return note, err
	}
	// Fetch updated note
	err = DB.QueryRow(
		"SELECT id, title, content, category, user_id, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2",
		id, userID,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CategoryID, &note.UserID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return note, err
	}
	// Fetch category name
	categoryName, err := getCategoryNameByID(categoryID)
	if err != nil {
		return note, err
	}
	note.CategoryName = categoryName
	return note, nil
}

// DeleteNote deletes a note if it belongs to the logged-in user.
func DeleteNote(userID, id int) error {
	_, err := DB.Exec(
		"DELETE FROM notes WHERE id = $1 AND user_id = $2",
		id, userID,
	)
	return err
}

// Helper function to fetch category name by ID
func getCategoryNameByID(categoryID int) (string, error) {
	var categoryName string
	err := DB.QueryRow("SELECT name FROM categories WHERE id = $1", categoryID).Scan(&categoryName)
	return categoryName, err
}

// GetAllNotesForAllUsers returns all notes for all users
func GetAllNotesForAllUsers() ([]models.Note, error) {
	rows, err := DB.Query(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		ORDER BY notes.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// GetAllNotesForAllUsersByCategory returns all notes for all users filtered by category
func GetAllNotesForAllUsersByCategory(categoryID int) ([]models.Note, error) {
	rows, err := DB.Query(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		WHERE notes.category = $1
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// GetAllNotesForAllUsersBySearch returns all notes for all users filtered by search term
func GetAllNotesForAllUsersBySearch(search string) ([]models.Note, error) {
	rows, err := DB.Query(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		WHERE (notes.title ILIKE '%' || $1 || '%' OR notes.content ILIKE '%' || $1 || '%')
	`, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// GetAllNotesForAllUsersByCategoryAndSearch returns all notes for all users filtered by category and search term
func GetAllNotesForAllUsersByCategoryAndSearch(categoryID int, search string) ([]models.Note, error) {
	rows, err := DB.Query(`
		SELECT notes.id, notes.title, notes.content, notes.created_at, notes.updated_at, categories.name, notes.category, users.username
		FROM notes
		LEFT JOIN categories ON notes.category = categories.id
		LEFT JOIN users ON notes.user_id = users.id
		WHERE notes.category = $1 AND (notes.title ILIKE '%' || $2 || '%' OR notes.content ILIKE '%' || $2 || '%')
	`, categoryID, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.CategoryName, &note.CategoryID, &note.Username); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}
