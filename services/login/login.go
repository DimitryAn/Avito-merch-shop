package login

import (
	"errors"
	"log"
	"os"
	"root/lib/errs"
	"root/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

const (
	issuer        = "avitoShop"
	defautBalance = 1000
)

var Secret string

type repository interface {
	CreateUser(user *models.User) error
	GetUserByName(u *models.User) error
}

type LoginService struct {
	repo repository
}

func NewLoginService(repo repository) *LoginService {

	if err := godotenv.Load(); err != nil {
		log.Fatal("can't load env file for service login")
	}

	Secret = os.Getenv("SECRET")

	if Secret == "" {
		log.Fatal("empty secret key for signed jwt")
	}

	return &LoginService{repo}
}

// создание пользователя в бд + создание jwt токена + хэш пароля
func (ls *LoginService) Login(name string, password string) (string, error) {

	u := &models.User{
		Username: name,
		Password: hashPassword(password),
		Balance:  defautBalance,
	}

	err := ls.repo.GetUserByName(u) // попытка получить пользователя

	if errors.Is(err, pgx.ErrNoRows) { // если нет пользователя в бд -> создаем
		err := ls.repo.CreateUser(u)
		if err != nil {
			return "", err
		}
	} else if err != nil || checkHashPassword(u.Password, password) != nil { // был, но пароль не совпал
		return "", errs.WrongPassword
	}

	return createJwt(u)
}

func createJwt(u *models.User) (string, error) {
	var claims = jwt.MapClaims{
		"iss":       issuer,
		"sub":       u.Username,
		"iat":       time.Now().Unix(),
		"user_id":   u.ID,
		"user_name": u.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	return signedToken, err
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func checkHashPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
