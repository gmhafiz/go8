package authentication

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gmhafiz/scs/v2"

	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/internal/utility/param"
	"github.com/gmhafiz/go8/internal/utility/request"
	"github.com/gmhafiz/go8/internal/utility/respond"
)

const (
	minPasswordLength = 13
)

var (
	ErrEmailRequired  = errors.New("email is required")
	ErrPasswordLength = fmt.Errorf("password must be at least %d characters", minPasswordLength)
)

type Handler struct {
	repo    Repo
	session *scs.SessionManager
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	if req.Email == "" {
		respond.Error(w, http.StatusBadRequest, ErrEmailRequired)
		return
	}

	if len(req.Password) < minPasswordLength {
		respond.Error(w, http.StatusBadRequest, ErrPasswordLength)
		return
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}

	if err := h.repo.Register(r.Context(), req.FirstName, req.LastName, req.Email, hashedPassword); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.Status(w, http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	ctx := r.Context()

	user, _, err := h.repo.Login(ctx, req)
	if err != nil {
		respond.Status(w, http.StatusUnauthorized)
		return
	}

	if err := h.session.RenewToken(ctx); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	h.session.Put(ctx, string(middleware.KeyID), user.ID)

	respond.Status(w, http.StatusOK)
}

func (h *Handler) Protected(w http.ResponseWriter, _ *http.Request) {
	respond.Json(w, http.StatusOK, map[string]string{"success": "yup!"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID := h.session.Get(r.Context(), string(middleware.KeyID))

	respond.Json(w, http.StatusOK, map[string]any{"user_id": userID})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.session.Destroy(r.Context())
	if err != nil {
		respond.Status(w, http.StatusBadRequest)
		return
	}
}

func (h *Handler) ForceLogout(w http.ResponseWriter, r *http.Request) {
	// Authorization is needed to ensure that only super admin can force delete
	currUser := h.session.Get(r.Context(), string(middleware.KeyID))
	// A more robust authorization is needed for real-world implementation.
	// For now, we naively check if user id is equal to 1.
	if currUser.(uint64) != 1 {
		respond.Status(w, http.StatusInternalServerError)
		return
	}

	userID, err := param.UInt64(r, "userID")
	if err != nil {
		respond.Status(w, http.StatusInternalServerError)
		return
	}

	ok, err := h.repo.Logout(r.Context(), userID)
	if err != nil {
		respond.Status(w, http.StatusInternalServerError)
		return
	}

	if !ok {
		respond.Json(w, http.StatusInternalServerError, map[string]string{"message": "unable to log out"})
	}
}

// Csrf stores a new csrf token in the database.
// For a Data modifying requests in <form action="" method="POST"> including PUT and PATCH,
// this csrf token needs to be attached along in the HTML along.
// Then check in this API for its existence.
func (h *Handler) Csrf(w http.ResponseWriter, r *http.Request) {
	_, ok := h.session.Get(r.Context(), string(middleware.KeyID)).(uint64)
	if !ok {
		respond.Error(w, http.StatusBadRequest, errors.New("you need to be logged in"))
		return
	}

	token, err := h.repo.Csrf(r.Context())
	if err != nil {
		respond.Status(w, http.StatusInternalServerError)
		return
	}

	respond.Json(w, http.StatusOK, &RespondCsrf{CsrfToken: token})
}

func NewHandler(session *scs.SessionManager, repo Repo) *Handler {
	return &Handler{
		repo:    repo,
		session: session,
	}
}
