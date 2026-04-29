package service

import (
	"errors"

	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/repository"
	"learn/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

// authService adalah struct yang mengimplementasikan interface AuthService.
// Struct ini memiliki dependensi ke UserRepository (ini disebut Dependency Injection).
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService menginisialisasi service autentikasi.
// Fungsi ini mengembalikan 'interface' AuthService, bukan struct 'authService' secara langsung.
// Tujuannya agar layer lain tidak bergantung pada implementasi detail, tapi pada kontrak (interface).
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

// Register menangani pendaftaran user baru.
// Parameter 'req' (dto.RegisterRequest) berisi data dari HTTP request.
// Mengembalikan '*dto.AuthResponse' (pointer ke struct response) dan 'error'.
func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// 1. Hash password menggunakan bcrypt (tidak boleh simpan plain text di DB).
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 2. Tentukan role (default: user).
	role := req.Role
	if role == "" {
		role = "user"
	}

	// 3. Mapping dari DTO (Data Transfer Object) ke Model Database.
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	// 4. Panggil repository untuk menyimpan data ke database.
	// s.userRepo adalah interface, jadi kita tidak peduli database apa yang dipakai di baliknya.
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 5. Generate JWT Token untuk autentikasi selanjutnya.
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// 6. Siapkan response untuk dikirim balik ke client.
	res := &dto.AuthResponse{
		Token: token,
	}
	res.User.ID = user.ID
	res.User.Name = user.Name
	res.User.Email = user.Email
	res.User.Role = user.Role

	return res, nil
}

// Login memvalidasi email dan password user.
// Jika valid, fungsi ini mengembalikan data user dan JWT token.
func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	res := &dto.AuthResponse{
		Token: token,
	}
	res.User.ID = user.ID
	res.User.Name = user.Name
	res.User.Email = user.Email
	res.User.Role = user.Role

	return res, nil
}
