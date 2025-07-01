package routes

import (
	"encoding/json"

	"azugo.io/azugo"
	"azugo.io/core/validation"
	"example.com/project/models"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)
var(
 	logger   *zap.Logger

	Users = models.UserList{
		Users: []models.User{
			{
				Username: "Oleg",
				Password: "123456",
				Email:    "olegs@gmail.com",
			},
			{
				Username: "Anna",
				Password: "654321",
				Email:    "anna@example.com",
			},
		},
	}
)

func init() {
	for i, user := range Users.Users {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			panic(err)
		}
		Users.Users[i].Password = string(hash)
	}
}

func InitLogger(l *zap.Logger) {
  logger = l
}

func GetUser (ctx *azugo.Context){
	var request models.User
	var response models.PublicUser
	found:=false

	if err := json.Unmarshal(ctx.Body.Bytes(), &request); err != nil{
        ctx.StatusCode(fasthttp.StatusBadRequest)
        logger.Warn("invalid json", zap.Error(err))
        return		
	}

	for _, record := range Users.Users{
		if record.Username == request.Username{
			response = models.PublicUser{
				Username: request.Username,
				Email:record.Email,
			}
			found = true
			break
		}

	}
	if !found{
		ctx.StatusCode(fasthttp.StatusBadRequest)
        logger.Info("No user found with provided username or bad json")
        return	
	}
	jsonBytes,err := json.Marshal(response)
	if err != nil {
		ctx.StatusCode(fasthttp.StatusNotFound)
		logger.Warn("invalid json", zap.Error(err))
		ctx.Context().SetBodyString("User not found")
		return
	}

	
	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(jsonBytes))
	logger.Info("User found successfully")

}

func CheckUsersExistence (ctx *azugo.Context){
	exist := false
	var user models.User
	if err := json.Unmarshal(ctx.Body.Bytes(), &user); err != nil{
        ctx.StatusCode(fasthttp.StatusBadRequest)
        logger.Warn("invalid json", zap.Error(err))
        return		
	}

	
	validate := validation.New()
	if err := user.Validate(validate); err != nil {
        ctx.StatusCode(fasthttp.StatusBadRequest)
		logger.Info("validation failed")
		return 
    }
	foundUsernameOrEmail := false

	for _, record := range Users.Users {
	if record.Username == user.Username || record.Email == user.Email {
		foundUsernameOrEmail = true
		err := bcrypt.CompareHashAndPassword([]byte(record.Password), []byte(user.Password))
		if err == nil {
		exist = true
		break
		}
	}
	}

	if exist {
	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("text/plain")
	ctx.Context().SetBodyString("All good")
	logger.Info("Validation successful")
	} else if foundUsernameOrEmail {
	ctx.StatusCode(fasthttp.StatusBadRequest)
	ctx.ContentType("text/plain")
	ctx.Context().SetBodyString("Incorrect password")
	logger.Info("Wrong password")
	} else {
	ctx.StatusCode(fasthttp.StatusBadRequest)
	ctx.ContentType("text/plain")
	ctx.Context().SetBodyString("Incorrect username or password")
	logger.Info("Validation unsuccessful")
	}
}

func GetAllUsers (ctx *azugo.Context){
	var publicList  models. PublicUserS
	for _, user := range Users.Users {
		publicUser := models.PublicUser{
			Username: user.Username,
			Email:    user.Email,
		}
		publicList.PublicUsers = append(publicList.PublicUsers, publicUser)
	}
	jsonBytes, err := json.Marshal(publicList)
	if err != nil {
		ctx.StatusCode(fasthttp.StatusInternalServerError)
		logger.Warn("invalid json", zap.Error(err))
		return
	}

	
	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(jsonBytes))
	logger.Info("Returned all users")

}


func AddUserToTheList (ctx *azugo.Context){
	
	var newUser models.User
	var addedPublicList  models. PublicUserS
	err := json.Unmarshal(ctx.Body.Bytes(),&newUser)
	if err != nil {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		return
	}

	validate := validation.New()
	if err := newUser.Validate(validate); err != nil {
        ctx.StatusCode(fasthttp.StatusBadRequest)
		logger.Info("Validation failed")
		return 
    }
	
	for _, record := range Users.Users{
		if record.Username == newUser.Username || record.Email == newUser.Email {
			ctx.StatusCode(fasthttp.StatusConflict)
			ctx.ContentType("text/plain")
			ctx.Context().SetBodyString("User with same Username or Email already exists")
			logger.Info("User already exists")
			return
		}
	}


	NewPasswordBytes := []byte(newUser.Password)
	
	hashNewPasswordBytes, err := bcrypt.GenerateFromPassword(NewPasswordBytes, bcrypt.MinCost)
	
	if err != nil {
		ctx.StatusCode(fasthttp.StatusInternalServerError)
		logger.Error("Error while encrypting password", zap.Error(err))
		return
	}

	newUser.Password=string(hashNewPasswordBytes)

	Users.Users = append(Users.Users, newUser)

	for _, record := range Users.Users{
		publicUser := models.PublicUser{
			Username: record.Username,
			Email:    record.Email,
		}
		addedPublicList.PublicUsers = append(addedPublicList.PublicUsers,publicUser)
	}

	usersJson, err := json.Marshal(addedPublicList)
	if err != nil {
    ctx.StatusCode(fasthttp.StatusInternalServerError)
    logger.Warn("invalid json", zap.Error(err))
    return
	}

	
	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(usersJson))
	logger.Info("User added successfully")

}

func PatchUser (ctx *azugo.Context){

	var PatchUser models.User
	var changedResponse  models. PublicUser
	found := false
	err := json.Unmarshal(ctx.Body.Bytes(), &PatchUser)

	if err != nil {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		return
	}

	if PatchUser.Username == "" && PatchUser.Email == "" && PatchUser.Password == "" {
    ctx.StatusCode(fasthttp.StatusBadRequest)
    ctx.Context().SetBodyString("Nothing to update")
    logger.Info("Patch failed: no fields to update")
    return
	}
	
	var changed models.User
	for i, record := range Users.Users {
    if record.Username == PatchUser.Username {
        if PatchUser.Email != "" {
            record.Email = PatchUser.Email
        }
        if PatchUser.Password != "" {
            hash, err := bcrypt.GenerateFromPassword([]byte(PatchUser.Password), bcrypt.MinCost)
			if err != nil {
				ctx.StatusCode(fasthttp.StatusInternalServerError)
				logger.Error("Failed to hash password", zap.Error(err))
				return
			}
            record.Password = string(hash)
        }
        Users.Users[i] = record
        changed = record
        found = true
        break
    }
}
	if !found{
		ctx.StatusCode(fasthttp.StatusNotFound)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("Nothing was found")
		logger.Info("Didn't find any records")
		return
	}
	
	changedResponse = models.PublicUser{
		Username: changed.Username,
		Email: changed.Email,
	}

	usersJson, err := json.Marshal(changedResponse)
	if err != nil {
    ctx.StatusCode(fasthttp.StatusInternalServerError)
    logger.Warn("invalid json", zap.Error(err))
    return
	}

	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(usersJson))
	logger.Info("User updated successfully")
}


func DeleteUser (ctx *azugo.Context){
	var UserToDelete models.User

	err := json.Unmarshal(ctx.Body.Bytes(), &UserToDelete)
	if err != nil {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		return
	}
	if UserToDelete.Username == "" && UserToDelete.Email == "" {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		ctx.Context().SetBodyString("Username or email required")
		return
	}
	found := false
	for i,record := range Users.Users{
		if record.Username == UserToDelete.Username || record.Email == UserToDelete.Email{
					Users.Users = append(Users.Users[:i],Users.Users[i+1:]...)
					found = true
					break
		}
	}
	
	if found {
		ctx.StatusCode(fasthttp.StatusOK)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("User deleted")
		logger.Info("User deleted")
	} else {
		ctx.StatusCode(fasthttp.StatusNotFound)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("User not found")
		logger.Info("User not found")
	}

}