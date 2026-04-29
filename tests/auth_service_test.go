package tests

import (
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository adalah "objek tiruan" (Mock).
// Dalam testing, kita tidak ingin memakai database beneran, jadi kita buat Mock ini.
type MockUserRepository struct {
	mock.Mock
}

// Create mensimulasikan fungsi simpan user.
func (m *MockUserRepository) Create(user *model.User) error {
	// m.Called(user) mencatat bahwa fungsi ini dipanggil dengan input 'user'.
	// 'args' di sini berisi nilai kembalian (return values) yang sudah kita atur di unit test.
	args := m.Called(user)
	// args.Error(0) artinya kita mengambil return value pertama (index 0) sebagai tipe 'error'.
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	// args.Get(0) mengambil return value pertama.
	// Kita harus melakukan "Type Assertion" (*model.User) karena Get() mengembalikan tipe 'any'.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	// Inisialisasi Mock
	mockRepo := new(MockUserRepository)
	// Inisialisasi Service dengan memasukkan Mock tadi (Dependency Injection).
	authService := service.NewAuthService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// "Setting" Mock: Jika fungsi Create dipanggil dengan user apapun, kembalikan nil (tidak error).
	mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	res, err := authService.Register(req)

	// Assertions: Mengecek apakah hasilnya sesuai ekspektasi kita.
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Test User", res.User.Name)
	assert.Equal(t, "test@example.com", res.User.Email)
	assert.Equal(t, "user", res.User.Role)
	assert.NotEmpty(t, res.Token)

	// Memastikan semua fungsi mock yang kita 'On' tadi benar-benar dipanggil.
	mockRepo.AssertExpectations(t)
}
