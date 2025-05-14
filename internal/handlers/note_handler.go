package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/AlexGithub777/notes-rest-app/internal/db"
	"github.com/AlexGithub777/notes-rest-app/internal/models"
	"github.com/AlexGithub777/notes-rest-app/internal/utils"
)

func JSONError(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]string{"error": msg})
}

// GET /notes/categories
func GetAllCategoriesHandler(c echo.Context) error {

	categories, err := db.GetAllCategories()
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch categories")
	}
	return c.JSON(http.StatusOK, categories)
}

// GET /notes (for logged in user)
func GetAllNotesHandler(c echo.Context) error {
	// get user id from the cookie
	cookie, err := c.Cookie("user_id")
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "User not logged in")
	}
	userID, err := strconv.Atoi(cookie.Value)

	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid user ID")
	}

	// Check for query params (category-id, search)
	categoryID := c.QueryParam("category-id")
	search := c.QueryParam("search")

	// log the query params
	fmt.Printf("categoryID: %s, search: %s\n", categoryID, search)

	// If categoryID is not empty, convert it to an integer
	var catID int
	if categoryID != "" {
		catID, err = strconv.Atoi(categoryID)
		if err != nil {
			return utils.JSONError(c, http.StatusBadRequest, "Invalid category ID")
		}
	}

	// If search is not empty, check if it is a valid string
	if search != "" {
		if len(search) < 1 {
			return utils.JSONError(c, http.StatusBadRequest, "Search term must be at least 1 characters long")
		}
	}

	// Handle different query parameter combinations
	if categoryID != "" && search != "" {
		// Both categoryID and search are provided
		notes, err := db.GetNotesForUserByCategoryAndSearch(userID, catID, search)
		if err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch notes")
		}
		return c.JSON(http.StatusOK, notes)
	} else if categoryID != "" {
		// Only categoryID is provided
		notes, err := db.GetNotesForUserByCategory(userID, catID)
		if err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch notes")
		}
		return c.JSON(http.StatusOK, notes)
	} else if search != "" {
		// Only search is provided
		notes, err := db.GetNotesForUserBySearch(userID, search)
		if err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch notes")
		}
		return c.JSON(http.StatusOK, notes)
	} else {
		// No query parameters, return all notes
		notes, err := db.GetNotesForUser(userID)
		if err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch notes")
		}
		return c.JSON(http.StatusOK, notes)
	}
}

// GET /notes/:id
func GetNoteByIDHandler(c echo.Context) error {
	// get user id from the cookie
	cookie, err := c.Cookie("user_id")
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "User not logged in")
	}
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid user ID")
	}
	// Fetch the note by ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid note ID")
	}

	note, err := db.GetNoteByID(userID, id)
	if err == sql.ErrNoRows {
		return utils.JSONError(c, http.StatusNotFound, "Note not found")
	} else if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch note")
	}

	return c.JSON(http.StatusOK, note)
}

// / POST /notes
func CreateNoteHandler(c echo.Context) error {
	// USE models CREATENOTEREQUEST TYPE
	noteReq := models.CreateNoteRequest{}

	if err := c.Bind(&noteReq); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid request payload")
	}

	// Check if all required fields are present
	if noteReq.Title == "" || noteReq.Content == "" || noteReq.CategoryID == 0 {
		return utils.JSONError(c, http.StatusBadRequest, "Missing required fields")
	}

	// Retrieve user_id from the cookie
	cookie, err := c.Cookie("user_id")
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "User not authenticated")
	}

	// Convert user_id to an integer
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "Invalid user ID")
	}

	now := time.Now()
	// Use the CreateNote function to insert into DB
	createdNote, err := db.CreateNote(userID, noteReq.Title, noteReq.Content, noteReq.CategoryID, now)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to create note")
	}

	return c.JSON(http.StatusCreated, createdNote)
}

// PUT /notes/:id
func UpdateNoteHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid note ID")
	}

	// USE models CREATENOTEREQUEST TYPE
	editNoteReq := models.UpdateNoteRequest{}

	if err := c.Bind(&editNoteReq); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid request payload")
	}

	// Check if all required fields are present
	if editNoteReq.Title == "" || editNoteReq.Content == "" || editNoteReq.CategoryID == 0 {
		return utils.JSONError(c, http.StatusBadRequest, "Missing required fields")
	}

	// Retrieve user_id from the cookie
	cookie, err := c.Cookie("user_id")
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "User not authenticated")
	}

	// Convert user_id to an integer
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "Invalid user ID")
	}

	// Check if the note exists and belongs to the user
	updatedNote, err := db.UpdateNote(userID, id, editNoteReq.Title, editNoteReq.Content, editNoteReq.CategoryID, time.Now())
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to update note")
	}

	return c.JSON(http.StatusOK, updatedNote)
}

// DELETE /notes/:id
func DeleteNoteHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "Invalid note ID")
	}

	// Retrieve user_id from the cookie
	cookie, err := c.Cookie("user_id")
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "User not authenticated")
	}

	// Convert user_id to an integer
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return utils.JSONError(c, http.StatusUnauthorized, "Invalid user ID")
	}

	// Delete the note if it exists and belongs to the user
	if err := db.DeleteNote(userID, id); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to delete note")
	}

	return c.NoContent(http.StatusNoContent)
}

// GET /notes/all
func GetAllNotesForAllUsersHandler(c echo.Context) error {
	// Check for query params (category-id, search)
	categoryID := c.QueryParam("category-id")
	search := c.QueryParam("search")

	// log the query params
	fmt.Printf("categoryID: %s, search: %s\n", categoryID, search)

	// If categoryID is not empty, convert it to an integer
	var catID int
	var err error
	if categoryID != "" {
		catID, err = strconv.Atoi(categoryID)
		if err != nil {
			return utils.JSONError(c, http.StatusBadRequest, "Invalid category ID")
		}
	}

	// If search is not empty, check if it is a valid string
	if search != "" {
		if len(search) < 1 {
			return utils.JSONError(c, http.StatusBadRequest, "Search term must be at least 1 characters long")
		}
	}

	// Handle different query parameter combinations
	if categoryID != "" && search != "" {
		// Both categoryID and search are provided
		notes, err := db.GetAllNotesForAllUsersByCategoryAndSearch(catID, search)
		if err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch notes")
		}
		return c.JSON(http.StatusOK, notes)
	} else if categoryID != "" {
		// Only categoryID is provided
		notes, err := db.GetAllNotesForAllUsersByCategory(catID)
		if err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch notes")
		}
		return c.JSON(http.StatusOK, notes)
	} else if search != "" {
		// Only search is provided
		notes, err := db.GetAllNotesForAllUsersBySearch(search)
		if err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch notes")
		}
		return c.JSON(http.StatusOK, notes)
	} else {
		// No query parameters, return all notes
		notes, err := db.GetAllNotesForAllUsers()
		if err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch notes")
		}
		return c.JSON(http.StatusOK, notes)
	}
}
