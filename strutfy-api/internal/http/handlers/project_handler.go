package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/masaru/strutfy-api/internal/project"
	"github.com/masaru/strutfy-api/internal/shared/response"
	"github.com/masaru/strutfy-api/internal/workflow"
)

type ProjectHandler struct {
	service *project.Service
	engine  *workflow.Engine
}

func NewProjectHandler(service *project.Service, engine *workflow.Engine) *ProjectHandler {
	return &ProjectHandler{
		service: service,
		engine:  engine,
	}
}

func getOrgID(r *http.Request) (uuid.UUID, error) {
	orgIDHeader := r.Header.Get("X-Org-ID")
	return uuid.Parse(orgIDHeader)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	orgID, err := getOrgID(r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid org id")
		return
	}

	var input project.CreateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.service.Create(r.Context(), orgID, input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, created)
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	orgID, err := getOrgID(r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid org id")
		return
	}

	projects, err := h.service.List(r.Context(), orgID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) Find(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orgID, err := getOrgID(r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid org id")
		return
	}

	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	found, err := h.service.Find(r.Context(), orgID, id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, found)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orgID, err := getOrgID(r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid org id")
		return
	}

	id, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	if err := h.service.Delete(r.Context(), orgID, id); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
	})
}
