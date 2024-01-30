package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/upload-media-auth/config"
	db "github.com/upload-media-auth/database"
	auth_helpers "github.com/upload-media-auth/helpers"
	"github.com/upload-media-auth/producers"
	"github.com/upload-media-auth/types"

	"github.com/labstack/echo/v4"
)

func authSignup(c echo.Context) error {
	user := new(types.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{Status: false, ErrMsg: "Invalid request body"})
	}

	fmt.Println(user)

	if err := validator.New().Struct(user); err != nil {
		// Handle validation errors
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{Status: false, ErrMsg: "Validation failed", ErrDetails: validationErrors})
	}

	_, userExistsError := auth_helpers.CheckUserExists(user.Email, false)

	if userExistsError != nil {
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{Status: false, ErrMsg: userExistsError.Error()})
	}

	fmt.Println(user)

	// Encrypting Provided Password and Storing User info in DB.
	savedUser, persistUserError := auth_helpers.PersistUser(user);

	if persistUserError != nil || savedUser == nil {
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{Status: false, ErrMsg: persistUserError.Error()})
	}

	// Generate Access and Refresh Tokens 
	access_token, accessTokenErr := auth_helpers.GenerateToken(
		config.GetConfig().JWT_SECRET_KEY, 
		time.Duration(config.GetConfig().ACCESS_TOKEN_EXPIRY),
		savedUser.Id,
	)

	refresh_token, refreshTokenErr := auth_helpers.GenerateToken(
		config.GetConfig().JWT_SECRET_KEY, 
		time.Duration(config.GetConfig().REFRESH_TOKEN_EXPIRY),
		savedUser.Id,
	)

	if accessTokenErr != nil || refreshTokenErr != nil {
		errStr := ""
		if accessTokenErr != nil && refreshTokenErr != nil {
			errStr = accessTokenErr.Error()
		} else if accessTokenErr != nil {
			errStr = accessTokenErr.Error()
		} else {
			errStr = refreshTokenErr.Error()
		}
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse {
			Status: false,
			ErrMsg: errStr,
		})
	}

	// Successful Sign Up Returns Access and Refresh Tokens :

	return c.JSON(http.StatusCreated, types.LoginUserResponse{
		Status: true,
		AccessToken: access_token,
		RefreshToken: refresh_token,
	})
}

func uploadPicture (c echo.Context) error {
	extractedUserId, extractionError := auth_helpers.ExtractUserId(c)
	if extractionError != nil {
		return c.JSON(http.StatusForbidden, types.LoginUserResponse{
			Status: false,
			ErrMsg: extractionError.Error(),
		})
	}

	// Get logged in user Data : 
	loggedInUser, fetchLoggedInUserError := auth_helpers.FetchUserInfoByID(extractedUserId)

	if fetchLoggedInUserError != nil {
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{
			Status: false,
			ErrMsg: fetchLoggedInUserError.Error(),
		})
	}

	cloudinaryClient, _ := cloudinary.NewFromParams(config.GetConfig().CLOUDINARY_CLOUD_NAME, config.GetConfig().CLOUDINARY_API_KEY, config.GetConfig().CLOUDINARY_API_SECRET)

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{
			Status: false,
			ErrMsg: "Failed to parse form data",
		})
	}

	// Retrieve the file
	file, err := form.File["file"][0].Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{
			Status: false,
			ErrMsg: "Failed to open file",
		})
	}
	defer file.Close()

	// Upload file to Cloudinary
	uploadResult, uploadError := cloudinaryClient.Upload.Upload(context.Background(), file, uploader.UploadParams{})

	if uploadError != nil {
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{
			Status: false,
			ErrMsg: uploadError.Error(),
		})
	}

	// Store picture Link in DB 
	_, persistingPictureErr := auth_helpers.PersistPicture(extractedUserId, uploadResult.SecureURL)

	if persistingPictureErr != nil {
		return c.JSON(http.StatusBadRequest, types.LoginUserResponse{
			Status: false,
			ErrMsg: persistingPictureErr.Error(),
		})
	}
	
	// Call to email service
	producers.PublishToQueue("email-queue", types.ProducedOrConsumedMessage{
		Email: loggedInUser.Email,
		PictureURL: uploadResult.SecureURL,
	})

	return c.JSON(http.StatusAccepted, types.UploadPicResponse{
		Status: true,
		PictureURL: uploadResult.SecureURL,
	})
}

func main () {
	e := echo.New()

	db.CreateDatabase();
	
	e.POST("/signup", authSignup);
	e.POST("/upload-picture", uploadPicture)

	e.Start(":8085")
}
