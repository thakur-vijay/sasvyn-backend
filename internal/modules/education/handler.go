package education

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/sasvyn/backend/internal/middleware"
	"github.com/sasvyn/backend/internal/modules/idempotency"
	"github.com/sasvyn/backend/internal/realtime"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	service            *Service
	idempotencyService *idempotency.Service
	publisher          *realtime.Publisher
}

func NewHandler(service *Service, idempotencyService *idempotency.Service, publisher *realtime.Publisher) *Handler {
	return &Handler{service: service, idempotencyService: idempotencyService, publisher: publisher}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	var request CreateEducationDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			education, err := h.service.Create(r.Context(), request, userID)
			if err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusCreated,
				Message:    "Education added successfully",
				Data:       education,
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
				"Education was not found",
			)
			return
		}

		log.Printf("Create Education error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Education could not be created",
		)
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventEducationCreated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish education.created event: %v", err)
		}
	}
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	educations, err := h.service.Fetch(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to fetch educations")
		return
	}

	response.WriteList(w, http.StatusOK, "Educations fetched successfully", educations)
}

func (h *Handler) FetchByID(w http.ResponseWriter, r *http.Request) {
	educationID := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())

	education, err := h.service.FetchByID(r.Context(), educationID, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "Education was not found")
			return
		}

		log.Printf("Fetch Education error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Failed to fetch Education")
		return
	}

	response.WriteItem(
		w,
		http.StatusOK,
		"Education fetched successfully",
		education,
	)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	educationID := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	var request UpdateEducationDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		updatedEducation, err := h.service.Update(r.Context(), request, educationID, userID)
		if err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Education updated successfully",
			Data:       updatedEducation,
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
				"Education was not found",
			)
			return
		}

		if strings.Contains(err.Error(), "23505") {
			response.Write(
				w,
				http.StatusConflict,
				"Education already exists",
			)
			return
		}

		log.Printf("Update Education error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Education could not be updated",
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
				Type: realtime.EventEducationUpdated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish education.updated event: %v", err)
		}
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	educationID := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		if err := h.service.Delete(r.Context(), educationID, userID); err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Education deleted successfully",
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
				"Education was not found",
			)
			return
		}

		log.Printf("Delete Education error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Education could not be deleted",
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
				Type: realtime.EventEducationDeleted,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish education.deleted event: %v", err)
		}
	}
}
