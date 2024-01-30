package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/upload-media-auth/config"
	db "github.com/upload-media-auth/database"
	"github.com/upload-media-auth/types"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserExists (email string, fromLogin bool) (*types.User, error) {
	var user types.User;
	fmt.Println(email)
	row := db.GetDBInstance().QueryRow("SELECT id, email, password, name FROM users WHERE email = $1;", email)

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Name)

	fmt.Println("new err 1", &user, err.Error())

	if err == sql.ErrNoRows {
		if fromLogin {
			return nil, errors.New ("user does not exist. please sign up")
		}

		fmt.Println("triggered")
		return nil, nil
	} else if err != sql.ErrNoRows {
		if (!fromLogin) {
			return nil, errors.New ("user already exists. please log in")
		}
	}

	// For Registration
	
	return &user, nil
}

func PersistUser (reqUser *types.User) (*types.User, error) {
	// Generate hashed password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("server error, unable to create your account")
	}

	// Convert hashed password to string
	hashedPassword := string(hashedPasswordBytes)

	// Generate a new UUID for the user ID
	userID := uuid.NewString()

	// Insert user into the database
	_, err = db.GetDBInstance().Exec(
		"INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)",
		userID, reqUser.Name, reqUser.Email, hashedPassword,
	)
	if err != nil {
		return nil, fmt.Errorf("error persisting user: %v", err)
	}

	// Create a Users struct for the saved user
	saveUser := &types.User{
		Id:       userID,
		Name:     reqUser.Name,
		Email:    reqUser.Email,
		Password: hashedPassword,
	}

	return saveUser, nil
}

func GenerateToken(key string, expiration time.Duration, user_id string) (string, error) {
	claims := &types.Claims{
		UserId: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration*time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key[:]))
	if err != nil {
		return "", errors.New("error occured while generating token")
	}

	return signedToken, nil
}

func ParseToken(tokenString string, key string) (*types.Claims, error) {
	claims := &types.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, errors.New("error occurred while parsing token")
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}

func ExtractUserId(c echo.Context) (string, error) {
	// Extract the Authorization header value
	authHeader := c.Request().Header.Get("Authorization")

	if authHeader == "" {
		// Authorization header is missing
		return "", errors.New("Unauthorized, missing token");
	}

	// Check if the Authorization header has the Bearer scheme
	const bearerScheme = "Bearer "
	var token = "";
	if len(authHeader) > len(bearerScheme) && authHeader[:len(bearerScheme)] == bearerScheme {
		// Extract the Bearer token
		token = authHeader[len(bearerScheme):]
	}

	if token == "" {
		return "", errors.New("Unauthorized, missing token");
	}

	tokenDerivedPayload, tokenErr := ParseToken(token, config.GetConfig().JWT_SECRET_KEY)

	if tokenErr != nil {
		return "", tokenErr
	}

	return tokenDerivedPayload.UserId, nil;
}

func PersistPicture (userId string, uploadURL string) (*types.Picture, error) {
	_, createPictureErr := db.GetDBInstance().Exec("INSERT INTO pictures (url, user_id) VALUES ($1, $2);", uploadURL, userId)

	if createPictureErr != nil {
		log.Println("Error executing query:", createPictureErr)
		return nil, errors.New(createPictureErr.Error())
	}

	fmt.Println("User created successfully")

	return &types.Picture{
		UserId: userId,
		PictureURL: uploadURL,
	}, nil
}

func FetchUserInfoByID (userId string) (*types.User, error) {
	rows, err := db.GetDBInstance().Query("SELECT id, name, email, password FROM users WHERE id=$1;", userId)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	var users []types.User

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			log.Println("Error scanning row :", err)
			return nil, errors.New(err.Error())
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, errors.New(err.Error())
	}

	if len(users) == 0 {
		return nil, errors.New("no users with this user ID exist")
	}

	result := users[0];
	return &result, nil;
}