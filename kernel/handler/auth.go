package handler

import (
	"errors"
	"labflow/dto"
	"labflow/httpx"
	"labflow/service"
	"labflow/turnstile"
	"net/http"
)

type AuthHandler struct {
	svc *service.Service
}

func NewAuthHandler(svc *service.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.Register(r.Context(), service.RegisterParams{
		Username:    req.Username,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Company:     req.Company,
		Address:     req.Address,
	})
	if err != nil {
		status, body := mapRegisterError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, dto.RegisterResponse{
		ID:          resp.ID,
		Xid:         resp.Xid,
		Username:    resp.Username,
		Email:       resp.Email,
		PhoneNumber: resp.PhoneNumber,
		Name:        resp.Name,
		Company:     resp.Company,
		Address:     resp.Address,
		CreatedAt:   resp.CreatedAt,
		Message:     resp.Message,
	})
}

func (h *AuthHandler) LoginWithCookie(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	if err := turnstile.Verify(req.TurnstileToken, httpx.ClientIP(r)); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "captcha_failed", Message: "Turnstile verification failed"})
		return
	}

	resp, token, err := h.svc.LoginWithCookie(r.Context(), service.LoginParams{
		Identifier: req.Identifier,
		Password:   req.Password,
	})
	if err != nil {
		status, body := mapRegisterError(err)
		httpx.WriteJSON(w, status, body)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		MaxAge:   3600 * 24,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	httpx.WriteJSON(w, http.StatusAccepted, dto.LoginResponse{
		Message: resp.Message,
		User: dto.UserDTO{
			Xid:         resp.User.Xid,
			Username:    resp.User.Username,
			Email:       resp.User.Email,
			Name:        resp.User.Name,
			PhoneNumber: resp.User.PhoneNumber,
			Company:     resp.User.Company,
		},
	})
}

func (h *AuthHandler) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyEmail
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.VerifyOtp(r.Context(), service.VerifyOTPParams{
		Email: req.Email,
		Otp:   req.Otp,
	})
	if err != nil {
		status, body := mapRegisterError(err)
		httpx.WriteJSON(w, status, body)
		return
	}
	httpx.WriteJSON(w, http.StatusAccepted, dto.SuccessResponse{Message: resp.Message})
}

func (h *AuthHandler) ResendOtp(w http.ResponseWriter, r *http.Request) {
	var req dto.ResendOTP
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.ResendOtp(r.Context(), service.ResendOTPParams{
		Email: req.Email,
	})
	if err != nil {
		status, body := mapRegisterError(err)
		httpx.WriteJSON(w, status, body)
		return
	}
	httpx.WriteJSON(w, http.StatusAccepted, dto.SuccessResponse{Message: resp.Message})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: "missing token"})
		return
	}

	token := cookie.Value
	if err := h.svc.Logout(r.Context(), token); err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error", Message: "failed to logout"})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})
	httpx.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Logout successful"})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, exists := httpx.UserID(r.Context())
	if !exists {
		httpx.WriteJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user not found"})
		return
	}

	resp, err := h.svc.GetMe(r.Context(), userID)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error", Message: "failed to fetch user"})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.MeResponse{
		ID:          int32(resp.ID),
		Xid:         resp.Xid,
		Username:    resp.Username,
		Email:       resp.Email,
		Name:        resp.Name,
		PhoneNumber: resp.PhoneNumber,
		Company:     resp.Company,
		Address:     resp.Address,
		Roles:       resp.Roles,
		Permissions: resp.Permissions,
		Avatar:      resp.Avatar,
	})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	err := h.svc.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUser) {
			writeErr(w, http.StatusNotFound, "not found", "User with this email does not exist.")
			return
		}
		writeErr(w, http.StatusInternalServerError, "internal error", "Failed to process forgot password request.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.SuccessResponse{Message: "OTP sent successfully"})
}

func (h *AuthHandler) VerifyResetOtp(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyResetOtpRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	token, err := h.svc.VerifyResetOtp(r.Context(), req.Email, req.Otp)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOtp) {
			writeErr(w, http.StatusBadRequest, "invalid otp", "Invalid OTP provided.")
			return
		}
		writeErr(w, http.StatusInternalServerError, "internal error", "Failed to verify OTP.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, dto.VerifyResetOtpResponse{ResetToken: token})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := dto.BindStrictJSON(r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.ResetPassword(r.Context(), req.Token, req.Password)
	if err != nil {
		status, body := mapResetPasswordError(err)
		httpx.WriteJSON(w, status, body)
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, dto.SuccessResponse{Message: resp.Message})
}

func mapResetPasswordError(err error) (int, dto.ErrorResponse) {
	if errors.Is(err, service.ErrInvalidResetToken) {
		return errResp(http.StatusUnauthorized, "unauthorized", "Invalid or expired reset token.")
	}
	return errResp(http.StatusInternalServerError, "internal_error", "Failed to reset password.")
}

func mapRegisterError(err error) (int, dto.ErrorResponse) {
	if errors.Is(err, service.ErrEmailTaken) {
		return errResp(http.StatusConflict, "conflict", "email already exists")
	}

	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		return errResp(http.StatusUnauthorized, "unauthorized", err.Error())
	case errors.Is(err, service.ErrInvalidVerification):
		return errResp(http.StatusBadRequest, "bad_request", err.Error())
	case errors.Is(err, service.ErrInvalidUser):
		return errResp(http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, service.ErrAlreadyVerified):
		return errResp(http.StatusBadRequest, "bad_request", err.Error())
	case errors.Is(err, service.ErrSendingEmail):
		return errResp(http.StatusBadRequest, "bad_request", err.Error())
	case errors.Is(err, service.ErrRedisFailed):
		return errResp(http.StatusBadRequest, "bad_request", err.Error())
	}

	switch {
	case errors.Is(err, service.ErrInvalidOtp):
		return errResp(http.StatusBadRequest, "bad_request", err.Error())
	case errors.Is(err, service.ErrOtpCooldown):
		return errResp(http.StatusTooManyRequests, "too_many_requests", err.Error())
	}
	return errResp(http.StatusInternalServerError, "internal_error", "failed to process request")
}
