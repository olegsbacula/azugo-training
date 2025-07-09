package routes

import (
	"encoding/json"
	"azugo.io/azugo"
	"example.com/project/models"
	"example.com/project/repository"
	"example.com/project/services"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/url"
)


func init() {
	for i, user := range repository.Users.Users {
		bytes, err := services.HashPassword(user.Password)
		if err != nil {
			panic(err)
		}
    	repository.Users.Users[i].Password = string(bytes)
	}
}

func InitLogger(l *zap.Logger) {
  repository.Logger = l
}




// LoginByID godoc
// @Summary     Login by username or email
// @Description Authenticate user by passing either username or email and password in the URL path
// @Tags        Users
// @Produce     plain
// @Param       id        path      string  true  "Username or Email"
// @Param       password  path      string  true  "Password"
// @Success     200       {string}  string  "Login successful"
// @Failure     400       {string}  string  "Missing credentials or wrong password"
// @Failure     404       {string}  string  "User not found"
// @Router      /login/{id}/{password} [get]
func LoginByID(ctx *azugo.Context) {
    RawId := ctx.Params.String("id")
    password := ctx.Params.String("password")
	id,err := url.QueryUnescape(RawId)
	if err!=nil{
		ctx.StatusCode(fasthttp.StatusInternalServerError)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("Cannot decode your email")
		repository.Logger.Info("Login failed: undecodable email")
        return
	}
    
    if id == "" || password == "" {
        ctx.StatusCode(fasthttp.StatusBadRequest)
        ctx.ContentType("text/plain")
        ctx.Context().SetBodyString("Username/email and password required")
        repository.Logger.Info("Login failed: missing credentials")
        return
    }


    var foundUser *models.User
    for _, u := range repository.Users.Users {
        if u.Username == id || u.Email == id {
            foundUser = &u
            break
        }
    }
    if foundUser == nil {
        ctx.StatusCode(fasthttp.StatusNotFound)
        ctx.ContentType("text/plain")
        ctx.Context().SetBodyString("User not found")
        repository.Logger.Info("Login failed: user not found", zap.String("id", id))
        return
    }

 
    if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password)); err != nil {
        ctx.StatusCode(fasthttp.StatusBadRequest)
        ctx.ContentType("text/plain")
        ctx.Context().SetBodyString("Incorrect password")
        repository.Logger.Info("Login failed: wrong password", zap.String("id", id))
        return
    }


    ctx.StatusCode(fasthttp.StatusOK)
    ctx.ContentType("text/plain")
    ctx.Context().SetBodyString("Login successful")
    repository.Logger.Info("Login successful", zap.String("id", id))
}

// FindUser godoc
// @Summary     Find user by username
// @Description Looks up a user by the username provided in the URL path and returns public user data.
// @Tags        Users
// @Produce     json
// @Param       username  path      string  true  "The user's unique username"
// @Success     200       {string}  string  "Username returned in JSON"
// @Failure     404       {string}  string  "User not found"
// @Router /find/{username} [post]
func GetUser (ctx *azugo.Context){
 	username := ctx.Params.String("username")            
	found:=false
	for _, record := range repository.Users.Users{
		if record.Username == username{
			found = true
			break
		}

	}
	if !found{
		ctx.StatusCode(fasthttp.StatusNotFound)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("User with this username doesn't exist")
        repository.Logger.Warn("No user found with provided username", zap.String("username", username))

        return	
	}

	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("text/plain")
	 ctx.JSON(map[string]string{"username": username})
	repository.Logger.Info("User found successfully")

}

// CheckUsersExistence godoc
// @Summary Check if user exists
// @Description Checks whether a user exists and whether the password is correct
// @Tags Users
// @Accept json
// @Produce plain
// @Param user body models.User true "Username and Password"
// @Success 200 {string} string "All good"
// @Failure 400 {string} string "Invalid JSON, wrong password or username"
// @Router /check [post]
func CheckUsersExistence (ctx *azugo.Context){
	exist := false
	var user models.User
	if err := json.Unmarshal(ctx.Body.Bytes(), &user); err != nil{
        ctx.StatusCode(fasthttp.StatusBadRequest)
        repository.Logger.Warn("invalid json", zap.Error(err))
        return		
	}

	
	if err := services.ValidUser(user); err != nil {
	ctx.StatusCode(fasthttp.StatusBadRequest)
	repository.Logger.Info("Validation failed", zap.Error(err))
	return
	}
	foundUsernameOrEmail := false

	for _, record := range repository.Users.Users {
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
	repository.Logger.Info("Validation successful")
	} else if foundUsernameOrEmail {
	ctx.StatusCode(fasthttp.StatusBadRequest)
	ctx.ContentType("text/plain")
	ctx.Context().SetBodyString("Incorrect password")
	repository.Logger.Info("Wrong password")
	} else {
	ctx.StatusCode(fasthttp.StatusBadRequest)
	ctx.ContentType("text/plain")
	ctx.Context().SetBodyString("Incorrect username or password")
	repository.Logger.Info("Validation unsuccessful")
	}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Returns a list of all users (without passwords)
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} models.PublicUserS
// @Failure 500 {string} string "Internal server error"
// @Router /list [get]
func GetAllUsers (ctx *azugo.Context){
	var publicList  models. PublicUserS
	for _, user := range repository.Users.Users {
		publicUser := models.PublicUser{
			Username: user.Username,
			Email:    user.Email,
		}
		publicList.PublicUsers = append(publicList.PublicUsers, publicUser)
	}
	jsonBytes, err := json.Marshal(publicList)
	if err != nil {
		ctx.StatusCode(fasthttp.StatusInternalServerError)
		repository.Logger.Warn("invalid json", zap.Error(err))
		return
	}

	
	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(jsonBytes))
	repository.Logger.Info("Returned all users")

}

// AddUserToTheList godoc
// @Summary Add a new user
// @Description Adds a new user and returns the updated user list
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "New user data"
// @Success 200 {object} models.PublicUserS
// @Failure 400 {string} string "Invalid JSON or validation failed"
// @Failure 409 {string} string "User with same username or email already exists"
// @Router /add [post]
func AddUserToTheList (ctx *azugo.Context){
	
	var newUser models.User
	var addedPublicList  models. PublicUserS
	err := json.Unmarshal(ctx.Body.Bytes(),&newUser)
	if err != nil {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := services.ValidUser(newUser); err != nil {
	ctx.StatusCode(fasthttp.StatusBadRequest)
	repository.Logger.Info("Validation failed", zap.Error(err))
	return
	}
	
	for _, record := range repository.Users.Users{
		if record.Username == newUser.Username || record.Email == newUser.Email {
			ctx.StatusCode(fasthttp.StatusConflict)
			ctx.ContentType("text/plain")
			ctx.Context().SetBodyString("User with same Username or Email already exists")
			repository.Logger.Info("User already exists")
			return
		}
	}
	
	hashedPassword ,err:=services.HashPassword(newUser.Password)
	
	if err != nil {
		ctx.StatusCode(fasthttp.StatusInternalServerError)
		repository.Logger.Error("Error while encrypting password", zap.Error(err))
		return
	}

	newUser.Password=string(hashedPassword)

	repository.Users.Users = append(repository.Users.Users, newUser)

	for _, record := range repository.Users.Users{
		publicUser := models.PublicUser{
			Username: record.Username,
			Email:    record.Email,
		}
		addedPublicList.PublicUsers = append(addedPublicList.PublicUsers,publicUser)
	}

	usersJson, err := json.Marshal(addedPublicList)
	if err != nil {
    ctx.StatusCode(fasthttp.StatusInternalServerError)
    repository.Logger.Warn("invalid json", zap.Error(err))
    return
	}

	
	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(usersJson))
	repository.Logger.Info("User added successfully")

}

// PatchUser godoc
// @Summary Update user
// @Description Updates user email or password by username
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User data to update"
// @Success 200 {object} models.PublicUser
// @Failure 400 {string} string "Validation failed or nothing to update"
// @Failure 404 {string} string "User not found"
// @Router /update [put]
func PatchUser (ctx *azugo.Context){

	var PatchUser models.User
	var changedResponse  models. PublicUser
	found := false
	err := json.Unmarshal(ctx.Body.Bytes(), &PatchUser)
	
	if err := services.ValidUser(PatchUser); err != nil {
	ctx.StatusCode(fasthttp.StatusBadRequest)
	repository.Logger.Info("Validation failed", zap.Error(err))
	return
	}
	if err != nil {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		return
	}

	if PatchUser.Username == "" && PatchUser.Email == "" && PatchUser.Password == "" {
    ctx.StatusCode(fasthttp.StatusBadRequest)
    ctx.Context().SetBodyString("Nothing to update")
    repository.Logger.Info("Patch failed: no fields to update")
    return
	}
	
	var changed models.User
	for i, record := range repository.Users.Users {
    if record.Username == PatchUser.Username {
        if PatchUser.Email != "" {
            record.Email = PatchUser.Email
        }
        if PatchUser.Password != "" {
            hash, err := services.HashPassword(PatchUser.Password)
			if err != nil {
				ctx.StatusCode(fasthttp.StatusInternalServerError)
				repository.Logger.Error("Failed to hash password", zap.Error(err))
				return
			}
            record.Password = string(hash)
        }
        repository.Users.Users[i] = record
        changed = record
        found = true
        break
    }
}
	if !found{
		ctx.StatusCode(fasthttp.StatusNotFound)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("Nothing was found")
		repository.Logger.Info("Didn't find any records")
		return
	}
	
	changedResponse = models.PublicUser{
		Username: changed.Username,
		Email: changed.Email,
	}

	usersJson, err := json.Marshal(changedResponse)
	if err != nil {
    ctx.StatusCode(fasthttp.StatusInternalServerError)
    repository.Logger.Warn("invalid json", zap.Error(err))
    return
	}

	ctx.StatusCode(fasthttp.StatusOK)
	ctx.ContentType("application/json")
	ctx.JSON(string(usersJson))
	repository.Logger.Info("User updated successfully")
}

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes a user by username or email
// @Tags Users
// @Accept json
// @Produce plain
// @Param user body models.User true "Username or Email to delete"
// @Success 200 {string} string "User deleted"
// @Failure 400 {string} string "Validation failed or bad request"
// @Failure 404 {string} string "User not found"
// @Router /delete [delete]
func DeleteUser (ctx *azugo.Context){
	var UserToDelete models.User

	err := json.Unmarshal(ctx.Body.Bytes(), &UserToDelete)
	if err != nil {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := services.ValidUser(UserToDelete); err != nil {
	ctx.StatusCode(fasthttp.StatusBadRequest)
	repository.Logger.Info("Validation failed", zap.Error(err))
	return
	}

	if UserToDelete.Username == "" && UserToDelete.Email == "" {
		ctx.StatusCode(fasthttp.StatusBadRequest)
		ctx.Context().SetBodyString("Username or email required")
		return
	}
	found := false
	for i,record := range repository.Users.Users{
		if record.Username == UserToDelete.Username || record.Email == UserToDelete.Email{
					repository.Users.Users = append(repository.Users.Users[:i],repository.Users.Users[i+1:]...)
					found = true
					break
		}
	}
	
	if found {
		ctx.StatusCode(fasthttp.StatusOK)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("User deleted")
		repository.Logger.Info("User deleted")
	} else {
		ctx.StatusCode(fasthttp.StatusNotFound)
		ctx.ContentType("text/plain")
		ctx.Context().SetBodyString("User not found")
		repository.Logger.Info("User not found")
	}

}