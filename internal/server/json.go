package server

import (
	"bytes"
	"cinema-recommender/internal/models"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

func (s *Server) decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if ct := r.Header.Get("Content-Type"); ct != "" {
		media, _, err := mime.ParseMediaType(ct)
		if err != nil || media != "application/json" {
			s.writeJSON(w, 415, map[string]string{"error": "use application/json"})
			return false
		}
	}
	data, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		s.writeJSON(w, 413, map[string]string{"error": "request body exceeds 64 KiB"})
		return false
	}
	if len(bytes.TrimSpace(data)) == 0 || bytes.TrimSpace(data)[0] != '{' {
		s.writeJSON(w, 400, map[string]string{"error": "send a JSON object"})
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dst)
	if err == nil {
		var extra any
		err = decoder.Decode(&extra)
		if err == io.EOF {
			return true
		}
	}
	status := http.StatusBadRequest
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		status = http.StatusRequestEntityTooLarge
	}
	s.writeJSON(w, status, map[string]string{"error": "send one valid JSON object with supported fields (maximum 64 KiB)"})
	return false
}
func validTier(t models.PrivilegeTier) bool {
	switch t {
	case "", "Standard", "Student", "Silver", "Gold", "Platinum":
		return true
	}
	return false
}
