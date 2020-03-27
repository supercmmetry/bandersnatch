package handlers

import (
	"bandersnatch/pkg"
	"bandersnatch/pkg/entities"
	"bandersnatch/pkg/player"
	"bandersnatch/utils"
	"encoding/json"
	"github.com/badoux/checkmail"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
)

func SignUp(playerSvc *player.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &entities.Player{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			utils.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := checkmail.ValidateFormat(p.Email); err != nil {
			utils.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(p.Name) == 0 {
			utils.RespWrap(w, http.StatusBadRequest, "name field is empty")
			return
		}

		if len(p.Password) == 0 {
			utils.RespWrap(w, http.StatusBadRequest, "password field is empty")
			return
		}

		token, err := playerSvc.SignUp(p)
		if err != nil {
			utils.RespWrap(w, http.StatusConflict, err.Error())
			return
		}

		tkString, err := token.SignedString([]byte(os.Getenv("JWT_PASSWORD")))
		if err != nil {
			utils.RespWrap(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.Wrap(w, map[string]interface{}{"token": tkString})
	}
}

func SignIn(playerSvc *player.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &entities.Player{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			utils.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := checkmail.ValidateFormat(p.Email); err != nil {
			utils.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		_, err := playerSvc.Find(p.Email)
		if err != nil {
			switch err {
			case pkg.ErrNotFound:
				utils.RespWrap(w, http.StatusUnauthorized, "the email does not exist")
			default:
				utils.RespWrap(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		token, err := playerSvc.SignIn(p)
		if err != nil {
			switch err {
			case pkg.ErrNotFound:
				utils.RespWrap(w, http.StatusUnauthorized, "incorrect password")
			default:
				utils.RespWrap(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		tkString, err := token.SignedString([]byte(os.Getenv("JWT_PASSWORD")))
		if err != nil {
			utils.RespWrap(w, http.StatusForbidden, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.Wrap(w, map[string]interface{}{"token": tkString})
	}
}

func ViewLeaderboard(playerSvc *player.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		leaders, err := playerSvc.ViewLeaderboard()
		if err != nil {
			utils.RespWrap(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.Wrap(w, map[string]interface{}{"leaderboard": leaders})
	}
}

func MakePlayerHandlers(router *httprouter.Router, playerSvc *player.Service) {
	router.HandlerFunc("POST", "/api/bandersnatch/signup", SignUp(playerSvc))
	router.HandlerFunc("POST", "/api/bandersnatch/signin", SignIn(playerSvc))
	router.HandlerFunc("POST", "/api/bandersnatch/leaderboard", ViewLeaderboard(playerSvc))
}
