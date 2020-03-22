package handlers

import (
	"bandersnatch/api"
	"bandersnatch/api/middleware"
	"bandersnatch/pkg/game"
	"bandersnatch/pkg/player"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func StartGame(playerSvc *player.Service, gameSvc *game.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		p, err := playerSvc.Find(tk.Email)
		if err != nil {
			api.RespWrap(w, http.StatusForbidden, err.Error())
			return
		}

		if tk.Id != p.Id {
			api.RespWrap(w, http.StatusForbidden, "player id mismatch")
			return
		}

		data := gameSvc.StartGame(p)

		w.WriteHeader(http.StatusOK)
		api.Wrap(w, map[string]interface{}{"data": data})
	}
}


func Play(playerSvc *player.Service, gameSvc *game.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		p, err := playerSvc.Find(tk.Email)
		if err != nil {
			api.RespWrap(w, http.StatusInternalServerError, err.Error())
			return
		}

		if tk.Id != p.Id {
			api.RespWrap(w, http.StatusForbidden, "player id mismatch")
			return
		}

		jsonMap := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&jsonMap); err != nil {
			api.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		option, ok := jsonMap["option"]
		if !ok {
			api.RespWrap(w, http.StatusBadRequest, "option not found")
			return
		}
		pl := &game.Player{Id:p.Id}
		data, err := gameSvc.Play(pl, api.OptionTypeCast(option.(float64)))
		if err != nil {
			api.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		artifacts, err := gameSvc.GetArtifacts(pl)
		if err != nil {
			api.RespWrap(w, http.StatusNotFound, "player not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		api.Wrap(w, map[string]interface{}{"data": data, "artifacts": artifacts})
	}
}

func MakeGameHandlers(router *httprouter.Router, playerSvc *player.Service, gameSvc *game.Service) {
	router.HandlerFunc("POST", "/api/bandersnatch/start", middleware.JwtAuth(StartGame(playerSvc, gameSvc)))
	router.HandlerFunc("POST", "/api/bandersnatch/play", middleware.JwtAuth(Play(playerSvc, gameSvc)))
}

