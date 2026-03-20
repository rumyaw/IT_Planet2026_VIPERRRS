package handlers

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}

