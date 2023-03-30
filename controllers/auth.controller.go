package controllers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shahriarsohan/new_blog/initializers"
	"github.com/shahriarsohan/new_blog/models"
	"github.com/shahriarsohan/new_blog/utils"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func valiDateEmail(email string) bool {
	Re := regexp.MustCompile(`[a-z0-9. %+/-]+@[a-z0-9. %+/-]+\.[a-z0-9. %+/-]`)
	return Re.MatchString(email)
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func SignUpUser(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}
	fmt.Println(data)
	fmt.Println("data...", data["email"])

	if len(data["password"].(string)) <= 8 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"msg": "Password must have more than 8 character",
		})
	}

	if !valiDateEmail(strings.TrimSpace(data["email"].(string))) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"msg": "Email is not valid",
		})
	}

	initializers.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)

	if userData.ID != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"msg": "Email already exixts",
		})
	}

	hasedPassword, err := utils.HashPassword(data["password"].(string))
	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"msg": "Unable to hash password",
		})
	}

	now := time.Now()
	user := models.User{
		Name:      data["name"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		Password:  hasedPassword,
		Role:      "user",
		Provider:  data["provider"].(string),
		Verified:  false,
		CreateAt:  now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		c.Status(400)
		return c.JSON(fiber.Map{
			"msg": "Email already exixts",
		})
	}

	config, _ := initializers.LoadConfig(".")
	code := randstr.String(20)
	verificationCode := utils.Encode(code)

	user.VerificationCode = verificationCode
	initializers.DB.Save(user)

	var firstName = user.Name
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "You accont verification code",
	}

	utils.SendMail(&user, &emailData)

	c.Status(200)
	return c.JSON(fiber.Map{
		"msg": "We sent an email with a verification code to " + user.Email,
	})
}

func VerifyEmail(c *fiber.Ctx) error {
	code := c.Params("verificationCode")
	verificationCode := utils.Encode(code)

	var updateUser models.User

	result := initializers.DB.First(&updateUser, "verification_code=?", verificationCode)

	if result.Error != nil {
		c.Status(404)
		c.JSON(fiber.Map{
			"msg": "something went wrong", // basically a 404 error
		})
	}

	if updateUser.Verified {
		c.Status(409)
		return c.JSON(fiber.Map{
			"msg": "User already verified",
		})
	}

	updateUser.VerificationCode = ""
	updateUser.Verified = true
	initializers.DB.Save(&updateUser)

	c.Status(200)
	return c.JSON(fiber.Map{
		"msg": "User verified successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userDara models.User

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Cannot parse data")
	}

	initializers.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userDara)
	if !userDara.Verified {
		c.Status(400)
		return c.JSON(fiber.Map{
			"msg": "User is not verified",
		})
	}
	if err := utils.VerifyPassword(userDara.Password, data["password"].(string)); err != nil {
		c.Status(400)
		c.JSON(fiber.Map{
			"msg": "Invalid credentials",
		})
	}

	config, _ := initializers.LoadConfig(".")

	token, err := utils.GenerateToken(config.TokenExpiration, userDara.ID, config.TokenSecret)
	if err != nil {
		if !userDara.Verified {
			c.Status(400)
			c.JSON(fiber.Map{
				"msg": err.Error(),
			})
		}
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	c.Status(200)
	return c.JSON(fiber.Map{
		"msg":   "Success",
		"token": token,
		"user":  userDara,
	})
}
