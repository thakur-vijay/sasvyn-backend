package experiences

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

	var request CreateExperienceDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			experience, err := h.service.Create(r.Context(), request, userID)
			if err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusCreated,
				Message:    "Experience added successfully",
				Data:       experience,
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
				"Experience was not found",
			)
			return
		}

		log.Printf("Create Experience error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Experience could not be created",
		)
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventExperienceCreated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish experience.created event: %v", err)
		}
	}
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	experiences, err := h.service.Fetch(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to fetch experiences")
		return
	}

	response.WriteList(w, http.StatusOK, "Experiences fetched successfully", experiences)
}

func (h *Handler) FetchByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())

	experience, err := h.service.FetchByID(r.Context(), id, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "Experience was not found")
			return
		}

		log.Printf("Fetch Experience error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Failed to fetch Experience")
		return
	}

	response.WriteItem(
		w,
		http.StatusOK,
		"Experience fetched successfully",
		experience,
	)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	var request UpdateExperienceDTO
	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		experience, err := h.service.Update(r.Context(), request, id, userID)
		if err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Experience updated successfully",
			Data:       experience,
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
				"Experience was not found",
			)
			return
		}

		if strings.Contains(err.Error(), "23505") {
			response.Write(
				w,
				http.StatusConflict,
				"Experience already exists",
			)
			return
		}

		log.Printf("Update Experience error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Experience could not be updated",
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
				Type: realtime.EventExperienceUpdated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish experience.updated event: %v", err)
		}
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		if err := h.service.Delete(r.Context(), id, userID); err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Experience deleted successfully",
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
				"Experience was not found",
			)
			return
		}

		log.Printf("Delete Experience error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Experience could not be deleted",
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
				Type: realtime.EventExperienceDeleted,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish experience.deleted event: %v", err)
		}
	}
}

func (h *Handler) CreateResponsibility(w http.ResponseWriter, r *http.Request) {
	experienceID := r.PathValue("experience_id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	var request CreateExperienceResponsibilityDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			responsibility, err := h.service.CreateResponsibility(
				r.Context(),
				request,
				experienceID,
				userID,
			)
			if err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusCreated,
				Message:    "Experience responsibility added successfully",
				Data:       responsibility,
			}, nil
		},
	)

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(w, http.StatusNotFound, "Experience was not found")
			return
		}

		log.Printf("Create Experience Responsibility error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Experience responsibility could not be created")
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventExperienceResponsibilityCreated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish experience_responsibility.created event: %v", err)
		}
	}
}

func (h *Handler) FetchResponsibilities(w http.ResponseWriter, r *http.Request) {
	experienceID := r.PathValue("experience_id")
	userID, _ := middleware.UserID(r.Context())

	responsibilities, err := h.service.FetchResponsibilities(
		r.Context(),
		experienceID,
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Write(w, http.StatusNotFound, "Experience was not found")
			return
		}

		log.Printf("Fetch Experience Responsibilities error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Failed to fetch experience responsibilities")
		return
	}

	response.WriteList(
		w,
		http.StatusOK,
		"Experience responsibilities fetched successfully",
		responsibilities,
	)
}

func (h *Handler) UpdateResponsibility(w http.ResponseWriter, r *http.Request) {
	experienceID := r.PathValue("experience_id")
	responsibilityID := r.PathValue("responsibility_id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	var request UpdateExperienceResponsibilityDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			responsibility, err := h.service.UpdateResponsibility(
				r.Context(),
				request,
				userID,
				experienceID,
				responsibilityID,
			)
			if err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusOK,
				Message:    "Experience responsibility updated successfully",
				Data:       responsibility,
			}, nil
		},
	)

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(w, http.StatusNotFound, "Experience responsibility was not found")
			return
		}

		log.Printf("Update Experience Responsibility error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Experience responsibility could not be updated")
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventExperienceResponsibilityUpdated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish experience_responsibility.updated event: %v", err)
		}
	}
}

func (h *Handler) DeleteResponsibility(w http.ResponseWriter, r *http.Request) {
	experienceID := r.PathValue("experience_id")
	responsibilityID := r.PathValue("responsibility_id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			if err := h.service.DeleteResponsibility(
				r.Context(),
				userID,
				experienceID,
				responsibilityID,
			); err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusOK,
				Message:    "Experience responsibility deleted successfully",
				Data: map[string]string{
					"id":            responsibilityID,
					"experience_id": experienceID,
				},
			}, nil
		},
	)

	if err != nil {
		if errors.Is(err, idempotency.ErrAlreadyProcessing) {
			response.Write(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			response.Write(w, http.StatusNotFound, "Experience responsibility was not found")
			return
		}

		log.Printf("Delete Experience Responsibility error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Experience responsibility could not be deleted")
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventExperienceResponsibilityDeleted,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish experience_responsibility.deleted event: %v", err)
		}
	}
}
