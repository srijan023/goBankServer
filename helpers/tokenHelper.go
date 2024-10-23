package helpers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/srijan023/jwtauth/models"
	"gorm.io/gorm"
)

type SignedDetails struct {
	Email string
	Id    string
	jwt.RegisteredClaims
}

var sercretKey string = os.Getenv("JWT_SECRET")

func validateRefreshToken(ctx *fiber.Ctx) (*SignedDetails, bool) {
	refreshToken := ctx.Cookies("refreshToken")
	if refreshToken == "" {
		return nil, false
	}

	token, err := jwt.ParseWithClaims(
		refreshToken,
		&SignedDetails{},
		func(refreshToken *jwt.Token) (interface{}, error) {
			return []byte(sercretKey), nil
		})

	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, false
	}

	if claims.ExpiresAt.Time.Before(time.Now().Local()) {
		return nil, false
	}

	return claims, true
}

// saving the refresh token in http only cookie securly to prevent js interception
func saveRefreshToken(ctx *fiber.Ctx, refreshToken string) {
	cookie := new(fiber.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.Expires = time.Now().Local().Add(time.Hour * 30 * 24) // expires in a month
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.SameSite = "Strict"

	ctx.Cookie(cookie)
}

func generateAlltokens(uuid string, email string) (string, string, error) {
	claims := &SignedDetails{
		Email: email,
		Id:    uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Minute * 15)),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(30*24))),
		},
	}

	token, tokenErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(sercretKey))

	if tokenErr != nil {
		return "", "", tokenErr
	}

	refreshToken, rfrErr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(sercretKey))

	if rfrErr != nil {
		return "", "", rfrErr
	}

	return token, refreshToken, nil
}

func UpdateAllTokens(ctx *fiber.Ctx, userInfo *models.User, db *gorm.DB) (string, error) {
	authToken, refreshToken, tokenizingErr := generateAlltokens(userInfo.ID.String(), userInfo.Email)

	if tokenizingErr != nil {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error while genering jwt token",
		})
		return "", tokenizingErr
	}

	saveRefreshToken(ctx, refreshToken)

	// NOTE: not saving the information into the database

	// err := db.Model(&userInfo).
	// 	Where("id = ?", userInfo.ID).
	// 	Select("Token", "RefreshToken").
	// 	Updates(models.User{
	// 		Token:        &authToken,
	// 		RefreshToken: &refreshToken,
	// 	})
	//
	// if err.Error != nil {
	// 	ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
	// 		"message": "Error while inserting data into the database",
	// 	})
	// 	return "", err.Error
	// }

	return authToken, nil
}

func generateAuthToken(id string, email string) (string, error) {
	claims := &SignedDetails{
		Email: email,
		Id:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Minute * 15)),
		},
	}

	token, tokenErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(sercretKey))

	if tokenErr != nil {
		return "", tokenErr
	}

	return token, nil
}

func ValidateToken(ctx *fiber.Ctx, signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(signedToken *jwt.Token) (interface{}, error) {
			return []byte(sercretKey), nil
		},
	)

	if err != nil {

		if errors.Is(err, jwt.ErrTokenExpired) {
			refClaim, isValid := validateRefreshToken(ctx)
			if isValid {
				token, err := generateAuthToken(refClaim.Id, refClaim.Email)
				if err != nil {
					return nil, err
				}
				ctx.Set("Authorization", token)
				return refClaim, nil
			}
			return nil, errors.New("token has expierd")
		}
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
