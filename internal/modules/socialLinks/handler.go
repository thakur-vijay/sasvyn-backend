package sociallinks

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sasvyn/backend/internal/modules/auth"
	"github.com/sasvyn/backend/internal/modules/idempotency"
	"github.com/sasvyn/backend/internal/response"
)

type Handler struct {
	repository         *Repository
	idempotencyService *idempotency.Service
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")

	var request CreateSocialLinkDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			now := time.Now().UTC()

			link := SocialLink{
				ID:        request.ID,
				UserID:    userID,
				Type:      request.Type,
				Url:       request.Url,
				CreatedAt: now,
				UpdatedAt: now,
			}

			createdLink, err := h.repository.Create(
				r.Context(),
				link,
			)
			if err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusCreated,
				Message:    "Social Link added successfully",
				Data:       createdLink,
			}, nil
		},
	)

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
				"social link was not found",
			)
			return
		}

		log.Printf("Create social link error: %v", err)

		response.Write(
			w,
			http.StatusInternalServerError,
			"social link could not be created",
		)
		return
	}

	response.WriteRaw(
		w,
		result.StatusCode,
		result.Body,
	)
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserID(r.Context())
	socialLinks, err := h.repository.Fetch(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, "failed to fetch social links")
		return
	}

	response.WriteList(w, http.StatusOK, "Social Links fetched successfully", socialLinks)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	socialLinkID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")

	var request UpdateSocialLinkDTO

	if err := response.DecodeJSONAndValidate(r, &request); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			updatedLink, err := h.repository.Update(
				r.Context(),
				socialLinkID,
				userID,
				request,
			)
			if err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusOK,
				Message:    "social link updated successfully",
				Data:       updatedLink,
			}, nil
		},
	)

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
				"social link was not found",
			)
			return
		}

		if strings.Contains(err.Error(), "23505") {
			response.Write(
				w,
				http.StatusConflict,
				"social link already exists",
			)
			return
		}

		log.Printf("Update social link error: %v", err)

		response.Write(
			w,
			http.StatusInternalServerError,
			"social link could not be updated",
		)
		return
	}

	response.WriteRaw(
		w,
		result.StatusCode,
		result.Body,
	)
}

func (h *Handler) FetchByID(w http.ResponseWriter, r *http.Request) {
	socialLinkID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())

	socialLink, err := h.repository.FetchByID(
		r.Context(),
		socialLinkID,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Write(w, http.StatusNotFound, "social link was not found")
			return
		}

		log.Printf("Fetch social link error: %v", err)
		response.Write(w, http.StatusInternalServerError, "failed to fetch social link")
		return
	}

	response.WriteItem(
		w,
		http.StatusOK,
		"Social Link fetched successfully",
		socialLink,
	)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	socialLinkID := r.PathValue("id")
	userID, _ := auth.UserID(r.Context())
	idempotencyKey := r.Header.Get("Idempotency-Key")

	result, err := h.idempotencyService.Execute(
		r.Context(),
		userID,
		idempotencyKey,
		func() (idempotency.Result, error) {
			if err := h.repository.Delete(
				r.Context(),
				socialLinkID,
				userID,
			); err != nil {
				return idempotency.Result{}, err
			}

			return idempotency.Result{
				StatusCode: http.StatusOK,
				Message:    "social link deleted successfully",
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
				"social link was not found",
			)
			return
		}

		log.Printf("Delete social link error: %v", err)

		response.Write(
			w,
			http.StatusInternalServerError,
			"social link could not be deleted",
		)
		return
	}

	response.WriteRaw(
		w,
		result.StatusCode,
		result.Body,
	)
}
