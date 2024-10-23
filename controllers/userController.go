package controllers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/srijan023/jwtauth/helpers"
	"github.com/srijan023/jwtauth/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var validate = validator.New()

type UserRouter struct {
	DB *gorm.DB
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func verifyPassword(actual string, stored string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(actual))

	if err != nil {
		return false, err
	}

	return true, nil
}

func updateAllTokens(userInfo *models.User, ctx *fiber.Ctx) {

}

func (ur *UserRouter) Signup() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		userInfo := &models.User{}

		err := ctx.BodyParser(&userInfo)

		if err != nil {
			ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"message": "Invalid form of data",
			})
			return err
		}

		userInfo.UserType = "CLIENT"

		validationError := validate.Struct(userInfo)

		if validationError != nil {
			ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"message": "Data validation failed",
			})
			return validationError
		}

		hashedPassword, hashErr := hashPassword(userInfo.Password)

		if hashErr != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"message": "failed hashing the password",
			})
			return hashErr
		}

		userInfo.Password = hashedPassword

		dataErr := ur.DB.Create(&userInfo).Error
		if dataErr != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"message": "Could not store data into the database",
			})
			return dataErr
		}

		token, tokenUpdateErr := helpers.UpdateAllTokens(ctx, userInfo, ur.DB)

		if tokenUpdateErr != nil {
			return tokenUpdateErr
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "Inserted data into the database",
			"token":   token,
		})

		return nil
	}
}

func (ur *UserRouter) Login() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		user := &models.User{}
		foundUser := &models.User{}

		err := ctx.BodyParser(&user)

		if err != nil {
			ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"message": "Body does not contain the required information",
			})
			return err
		}

		// use Find instead of First because it does not handle no data found error properly
		fetchErr := ur.DB.Where("email = ?", user.Email).Find(&foundUser)

		if fetchErr.RowsAffected == 0 {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"message": "there is no user with that data",
			})
			return fetchErr.Error
		}

		if fetchErr != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"message": "could not get the information from database",
			})
			return fetchErr.Error
		}

		if foundUser.Email == "" {
			ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
				"message": "There is no such user",
			})
			return nil
		}

		isValidPassword, _ := verifyPassword(user.Password, foundUser.Password)

		if !isValidPassword {
			ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"message": "Password does not match",
			})
			return nil
		}

		token, err := helpers.UpdateAllTokens(ctx, foundUser, ur.DB)

		if err != nil {
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"authToken": token,
			"message":   "Login successful",
		})

		return nil

	}
}

func (ur *UserRouter) GetUsers() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		users := &[]models.User{}

		err := ur.DB.Find(&users).Error

		if err != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"message": "Could not get all the users",
			})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "Users fetched successfully",
			"data":    users,
		})
		return nil
	}
}

func (ur *UserRouter) GetUser() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Params("id")

		if userId == "" {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "user id can't be empty"})
			return nil
		}

		// validate if the user is requesting the data of different user or himself
		// err := helper.MatchUserById(ctx, userId)
		// if err != nil {
		// 	ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "No such user found"})
		// 	return nil
		// }

		user := &models.User{}

		err := ur.DB.Where("id = ?", userId).First(user).Error
		if err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "No such user exists"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message": "Successfully fetched user",
				"data":    user,
			})

		return nil

	}
}
