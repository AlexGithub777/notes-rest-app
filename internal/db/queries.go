package db

import (
	"database/sql"
	"github.com/AlexGithub777/notes-rest-app/internal/models"
	"time"
)

func GetAllNotes() ([]models.Note, error) {
	rows, err := DB.Query("SELECT id, title, content, created_at, updated_at FROM notes ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func GetNoteByID(id int) (models.Note, error) {
	var note models.Note
	err := DB.QueryRow("SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1", id).
		Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
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

func CreateNote(title, content string, now time.Time) (int, error) {
	var noteID int
	err := DB.QueryRow(
		"INSERT INTO notes (title, content, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		title, content, now, now,
	).Scan(&noteID)
	return noteID, err
}

func UpdateNote(id int, title, content string, now time.Time) error {
	_, err := DB.Exec("UPDATE notes SET title = $1, content = $2, updated_at = $3 WHERE id = $4", title, content, now, id)
	return err
}

func DeleteNote(id int) error {
	_, err := DB.Exec("DELETE FROM notes WHERE id = $1", id)
	return err
}
