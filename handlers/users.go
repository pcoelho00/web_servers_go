package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/pcoelho00/server_go/jsondecoders"
)

type PostUser struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ExpireSecs int    `json:"expires_in_seconds"`
}

type ResponseUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *ApiConfig) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := PostUser{}
	err := decoder.Decode(&user)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	NewUser, err := cfg.DB.CreateUser(user.Email, user.Password)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Couldn't Create the User")
	}

	new_token, err := createToken(cfg.JwtSecret, user.ExpireSecs, NewUser.Id)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Couldn't create the token")
	}

	jsondecoders.RespondWithJson(w, http.StatusCreated, ResponseUser{
		Id:    NewUser.Id,
		Email: NewUser.Email,
		Token: new_token,
	})

}

func (cfg *ApiConfig) GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	Users, err := cfg.DB.GetUsers()
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Couldn't return msg from the database")
	}

	jsondecoders.RespondWithJson(w, http.StatusOK, Users)
}

func (cfg *ApiConfig) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	UserId, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Error getting ID")
	}
	User, err := cfg.DB.GetPublicUser(UserId)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusNotFound, "Chirp ID doesn't exists")
	} else {
		jsondecoders.RespondWithJson(w, http.StatusOK, User)
	}

}
