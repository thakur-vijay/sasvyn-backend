package skills

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/idempotency"
	"github.com/sasvyn/backend/internal/realtime"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	repository         *Repository
	idempotencyService *idempotency.Service
	publisher          *realtime.Publisher
}

func NewHandler(repository *Repository, idempotencyService *idempotency.Service, publisher *realtime.Publisher) *Handler {
	return &Handler{repository: repository, idempotencyService: idempotencyService, publisher: publisher}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	var request CreateSkillDTO

	if err := response.DecodeJSON(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			now := time.Now().UTC()

			skill := Skill{
				ID:          request.ID,
				UserID:      userID,
				Skill:       request.Skill,
				Category:    request.Category,
				SyncVersion: 1,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			if err := h.repository.Create(r.Context(), skill); err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusCreated,
				Message:    "Skill added successfully",
				Data:       skill,
			}, nil
		},
	)

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(
				w,
				http.StatusNotFound,
				"Skill was not found",
			)
			return
		}

		log.Printf("Create Skill error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Skill could not be created",
		)
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventSkillCreated,
				Data: result.Data,
			},
		); err != nil {
			log.Printf("failed to publish skill.created event: %v", err)
		}
	}
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserID(r.Context())
	skills, err := h.repository.Fetch(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to fetch skills")
		return
	}

	response.WriteList(w, http.StatusOK, "Skill fetched successfully", skills)
}

func (h *Handler) FetchByID(w http.ResponseWriter, r *http.Request) {
	skillID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())

	skill, err := h.repository.FetchByID(r.Context(), skillID, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "Skill was not found")
			return
		}

		log.Printf("Fetch Skill error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Failed to fetch Skill")
		return
	}

	response.WriteItem(
		w,
		http.StatusOK,
		"Skill fetched successfully",
		skill,
	)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	skillID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")

	var request UpdateSkillDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		updatedSkill, err := h.repository.Update(r.Context(), skillID, userID, request)
		if err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Skill updated successfully",
			Data:       updatedSkill,
		}, nil
	})

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(
				w,
				http.StatusConflict,
				err.Error(),
			)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(
				w,
				http.StatusNotFound,
				"Skill was not found",
			)
			return
		}

		if strings.Contains(err.Error(), "23505") {
			response.Write(
				w,
				http.StatusConflict,
				"Skill already exists",
			)
			return
		}

		log.Printf("Update Skill error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Skill could not be updated",
		)

		return
	}

	response.WriteRaw(
		w,
		result.StatusCode,
		result.Body,
	)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventSkillUpdated,
				Data: result.Data,
			},
		); err != nil {
			log.Printf("failed to publish skill.created event: %v", err)
		}
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	skillID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		if err := h.repository.Delete(r.Context(), skillID, userID); err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Skill deleted successfully",
		}, nil
	})

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(
				w,
				http.StatusNotFound,
				"Skill was not found",
			)
			return
		}

		log.Printf("Delete Skill error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Skill could not be deleted",
		)
		return
	}

	response.WriteRaw(
		w,
		result.StatusCode,
		result.Body,
	)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventSkillDeleted,
				Data: result.Data,
			},
		); err != nil {
			log.Printf("failed to publish skill.created event: %v", err)
		}
	}
}
