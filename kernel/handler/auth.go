package handler

import (
	"errors"
	"labflow/dto"
	"labflow/service"
	"labflow/turnstile"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.Service
}

func NewAuthHandler(svc *service.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.Register(c.Request.Context(), service.RegisterParams{
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
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
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

func (h *AuthHandler) LoginWithCookie(c *gin.Context) {
	var req dto.LoginRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	if err := turnstile.Verify(req.TurnstileToken, c.ClientIP()); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "captcha_failed", Message: "Turnstile verification failed"})
		return
	}

	resp, token, err := h.svc.LoginWithCookie(c.Request.Context(), service.LoginParams{
		Identifier: req.Identifier,
		Password:   req.Password,
	})
	if err != nil {
		status, body := mapRegisterError(err)
		c.JSON(status, body)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"Authorization",
		token,
		3600*24,
		"/",
		"",
		false,
		true,
	)
	c.JSON(http.StatusAccepted, dto.LoginResponse{
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

func (h *AuthHandler) VerifyOtp(c *gin.Context) {
	var req dto.VerifyEmail
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.VerifyOtp(c.Request.Context(), service.VerifyOTPParams{
		Email: req.Email,
		Otp:   req.Otp,
	})
	if err != nil {
		status, body := mapRegisterError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusAccepted, dto.SuccessResponse{Message: resp.Message})
}

func (h *AuthHandler) ResendOtp(c *gin.Context) {
	var req dto.ResendOTP
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.ResendOtp(c.Request.Context(), service.ResendOTPParams{
		Email: req.Email,
	})
	if err != nil {
		status, body := mapRegisterError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusAccepted, dto.SuccessResponse{Message: resp.Message})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: "missing token"})
		return
	}

	if err := h.svc.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error", Message: "failed to logout"})
		return
	}
	c.SetCookie(
		"Authorization",
		token,
		-1,
		"/",
		"",
		false,
		true,
	)
	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Logout successful"})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized", Message: "user not found"})
		return
	}

	resp, err := h.svc.GetMe(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error", Message: "failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, dto.MeResponse{
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

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	err := h.svc.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUser) {
			writeErr(c, http.StatusNotFound, "not found", "User with this email does not exist.")
			return
		}
		writeErr(c, http.StatusInternalServerError, "internal error", "Failed to process forgot password request.")
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "OTP sent successfully"})
}

func (h *AuthHandler) VerifyResetOtp(c *gin.Context) {
	var req dto.VerifyResetOtpRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	token, err := h.svc.VerifyResetOtp(c.Request.Context(), req.Email, req.Otp)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOtp) {
			writeErr(c, http.StatusBadRequest, "invalid otp", "Invalid OTP provided.")
			return
		}
		writeErr(c, http.StatusInternalServerError, "internal error", "Failed to verify OTP.")
		return
	}

	c.JSON(http.StatusOK, dto.VerifyResetOtpResponse{ResetToken: token})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := dto.BindStrictJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	resp, err := h.svc.ResetPassword(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		status, body := mapResetPasswordError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusAccepted, dto.SuccessResponse{Message: resp.Message})
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
