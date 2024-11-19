// internal/handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
)

type UserHandler struct {
	UserService services.UserServiceInterface
}

func NewUserHandler(us services.UserServiceInterface) *UserHandler {
	return &UserHandler{UserService: us}
}

func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uh.UserService.Signup(r.Context(), &user); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Signup successful. Please verify your email."})
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := uh.UserService.Login(r.Context(), &loginData)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.WriteJSON(w, map[string]string{"token": token})
}

func (uh *UserHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uh.UserService.ResendOTP(r.Context(), requestData.Email); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "A new OTP has been sent to your email address."})
}

func (uh *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := uh.UserService.VerifyEmail(r.Context(), requestData.Email, requestData.OTP)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Email verified successfully", "token": token})
}

func (uh *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uh.UserService.ForgotPassword(r.Context(), requestData.Email); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "If the email exists, an OTP has been sent."})
}

func (uh *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uh.UserService.ResetPassword(r.Context(), requestData.Email, requestData.OTP, requestData.NewPassword); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Password has been reset successfully."})
}

func (uh *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	userInfo, err := uh.UserService.GetUserInfo(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.WriteJSON(w, userInfo)
}

func (uh *UserHandler) SearchUsersByUsername(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		utils.WriteJSONError(w, "Missing search query", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)

	results, err := uh.UserService.SearchUsersByUsername(r.Context(), userEmail, query)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, results)
}
