package handlers

import (
	"bandersnatch/api"
	"bandersnatch/pkg"
	"bandersnatch/pkg/entities"
	"bandersnatch/pkg/player"
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
			api.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := checkmail.ValidateFormat(p.Email); err != nil {
			api.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		token, err := playerSvc.SignUp(p)
		if err != nil {
			api.RespWrap(w, http.StatusConflict, err.Error())
			return
		}

		tkString, err := token.SignedString([]byte(os.Getenv("JWT_PASSWORD")))
		if err != nil {
			api.RespWrap(w, http.StatusInternalServerError, err.Error())
			return
		}

		api.Wrap(w, map[string]interface{}{"token": tkString, "id": p.Id})
	}
}

func SignIn(playerSvc *player.Service) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		p := &entities.Player{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			api.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := checkmail.ValidateFormat(p.Email); err != nil {
			api.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		p, err := playerSvc.Find(p.Email)
		if err != nil {
			switch err {
			case pkg.ErrNotFound:
				api.RespWrap(w, http.StatusUnauthorized, "the email does not exist")
			default:
				api.RespWrap(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		token, err := playerSvc.SignIn(p)
		if err != nil {
			switch err {
			case pkg.ErrNotFound:
				api.RespWrap(w, http.StatusUnauthorized, "incorrect password")
			default:
				api.RespWrap(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		tkString, err := token.SignedString([]byte(os.Getenv("JWT_PASSWORD")))
		if err != nil {
			api.RespWrap(w, http.StatusInternalServerError, err.Error())
			return
		}

		api.Wrap(w, map[string]interface{}{"token": tkString, "id": p.Id})
	}
}

func MakeHandler(router *httprouter.Router, playerSvc *player.Service) {
	router.HandlerFunc("POST", "/api/bandersnatch/signup", SignUp(playerSvc))
	router.HandlerFunc("POST", "/api/bandersnatch/signin", SignIn(playerSvc))
}

