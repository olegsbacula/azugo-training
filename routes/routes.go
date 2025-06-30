package routes

import (
	"encoding/json"
	"azugo.io/core/validation"
	"azugo.io/azugo"
	"example.com/project/models"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
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

func CheckUsersExistence (ctx *azugo.Context){
	exist := false
	var user models.User
	if err := json.Unmarshal(ctx.Body.Bytes(), &user); err != nil{
        ctx.StatusCode(fasthttp.StatusConflict)
        logger.Warn("invalid json", zap.Error(err))
        return		
	}

	
	validate := validation.New()
	if err := user.Validate(validate); err != nil {
        ctx.StatusCode(fasthttp.StatusBadRequest)
		logger.Info("validation failed")
    }

	for _, record := range Users.Users {
    if record == user {
			exist = true
			break
		}
	}

	if  exist {
		ctx.StatusCode(fasthttp.StatusOK)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("All good")
		logger.Info("Validation successful")
	} else {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("Wrong")
		logger.Info("Validation unsuccessful")
	}
}

func GetAllUsers (ctx *azugo.Context){

	jsonBytes, err := json.Marshal(Users)
	if err != nil {
    ctx.StatusCode(fasthttp.StatusInternalServerError)
    logger.Warn("invalid json", zap.Error(err))
    return
	}

	
	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(jsonBytes))
	logger.Info("Validation unsuccessful")

}


func AddUserToTheList (ctx *azugo.Context){
	
	var newUser models.User

	err := json.Unmarshal(ctx.Body.Bytes(),&newUser)
	if err != nil {
		ctx.StatusCode(fasthttp.StatusBadRequest)
	}
	Users.Users = append(Users.Users, newUser)

	usersJson, err := json.Marshal(Users.Users)
	if err != nil {
    ctx.StatusCode(fasthttp.StatusInternalServerError)
    logger.Warn("invalid json", zap.Error(err))
    return
	}
	
	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(usersJson))
	logger.Info("Validation unsuccessful")

}

func PatchUser (ctx *azugo.Context){

	var PatchUser models.User
	found := false
	err := json.Unmarshal(ctx.Body.Bytes(), &PatchUser)

	if err != nil {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		return
	}
	var changed models.User
	for i, record := range Users.Users {
    if record.Username == PatchUser.Username || record.Email == PatchUser.Email {
        Users.Users[i] = PatchUser
		changed = Users.Users[i]
		found = true
        break
   	 }
  }
	if !found{
		ctx.StatusCode(fasthttp.StatusBadRequest)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("Nothing was found")
		logger.Info("Didn't find any records")
		return
	}
	usersJson, err := json.Marshal(changed)
	if err != nil {
    ctx.StatusCode(fasthttp.StatusInternalServerError)
    logger.Warn("invalid json", zap.Error(err))
    return
	}

	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(usersJson))
	logger.Info("Validation unsuccessful")
}