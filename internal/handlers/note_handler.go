package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/AlexGithub777/notes-rest-app/internal/db"
)

// Helper to return consistent error responses
func echoError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"error": message})
}

// GET /notes
func GetAllNotesHandler(c echo.Context) error {
	notes, err := db.GetAllNotes()
	if err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to fetch notes")
	}
	return c.JSON(http.StatusOK, notes)
}

// GET /notes/:id
func GetNoteByIDHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echoError(c, http.StatusBadRequest, "Invalid note ID")
	}

	note, err := db.GetNoteByID(id)
	if err == sql.ErrNoRows {
		return echoError(c, http.StatusNotFound, "Note not found")
	} else if err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to fetch note")
	}

	return c.JSON(http.StatusOK, note)
}

// POST /notes
func CreateNoteHandler(c echo.Context) error {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.Bind(&input); err != nil {
		return echoError(c, http.StatusBadRequest, "Invalid input")
	}

	now := time.Now()
	noteID, err := db.CreateNote(input.Title, input.Content, now)
	if err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to create note")
	}

	note, err := db.GetNoteByID(noteID)
	if err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to fetch created note")
	}
	return c.JSON(http.StatusCreated, note)
}

// PUT /notes/:id
func UpdateNoteHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echoError(c, http.StatusBadRequest, "Invalid note ID")
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.Bind(&input); err != nil {
		return echoError(c, http.StatusBadRequest, "Invalid input")
	}

	exists, err := db.NoteExists(id)
	if err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to check note existence")
	}
	if !exists {
		return echoError(c, http.StatusNotFound, "Note not found")
	}

	now := time.Now()
	if err := db.UpdateNote(id, input.Title, input.Content, now); err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to update note")
	}

	note, err := db.GetNoteByID(id)
	if err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to fetch updated note")
	}
	return c.JSON(http.StatusOK, note)
}

// DELETE /notes/:id
func DeleteNoteHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echoError(c, http.StatusBadRequest, "Invalid note ID")
	}

	exists, err := db.NoteExists(id)
	if err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to check note existence")
	}
	if !exists {
		return echoError(c, http.StatusNotFound, "Note not found")
	}

	if err := db.DeleteNote(id); err != nil {
		return echoError(c, http.StatusInternalServerError, "Failed to delete note")
	}

	return c.NoContent(http.StatusNoContent)
}
