package documents

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

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

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	var request CreateDocumentDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			document, err := h.service.Create(
				r.Context(),
				request,
				userID,
			)

			if err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusCreated,
				Message:    "Document added successfully",
				Data:       document,
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
				"Document was not found",
			)
			return
		}

		log.Printf("Create Document error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Document could not be created",
		)
		return
	}

	response.WriteRaw(w, result.StatusCode, result.Body)

	if !result.Replayed {
		if err := h.publisher.Publish(
			r.Context(),
			userID,
			realtime.Event{
				Type: realtime.EventDocumentCreated,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish document.created event: %v", err)
		}
	}
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())
	documents, err := h.service.Fetch(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to fetch documents")
		return
	}

	response.WriteList(w, http.StatusOK, "Documents fetched successfully", documents)
}

func (h *Handler) FetchByID(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())

	document, err := h.service.FetchByID(r.Context(), documentID, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "Document was not found")
			return
		}

		log.Printf("Fetch Document error: %v", err)
		response.Write(w, http.StatusInternalServerError, "Failed to fetch Document")
		return
	}

	response.WriteItem(
		w,
		http.StatusOK,
		"Document fetched successfully",
		document,
	)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	userID, _ := middleware.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")
	clientID := r.Header.Get("X-Client-ID")

	result, err := h.idempotencyService.Execute(r.Context(), userID, idempotencyKey, func() (idempotency.Result, error) {
		if err := h.service.Delete(r.Context(), documentID, userID); err != nil {
			return idempotency.Result{}, err
		}

		return idempotency.Result{
			StatusCode: http.StatusOK,
			Message:    "Document deleted successfully",
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
				"Document was not found",
			)
			return
		}

		log.Printf("Delete Document error: %v", err)
		response.Write(
			w,
			http.StatusInternalServerError,
			"Document could not be deleted",
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
				Type: realtime.EventDocumentDeleted,
				Data: result.Data,
			},
			clientID,
		); err != nil {
			log.Printf("failed to publish document.deleted event: %v", err)
		}
	}
}
