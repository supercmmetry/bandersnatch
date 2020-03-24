package handlers

import (
	"bandersnatch/api/middleware"
	"bandersnatch/pkg/game"
	"bandersnatch/pkg/player"
	"bandersnatch/utils"
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
			utils.RespWrap(w, http.StatusForbidden, err.Error())
			return
		}

		if tk.Id != p.Id {
			utils.RespWrap(w, http.StatusForbidden, "player id mismatch")
			return
		}

		data, err := gameSvc.StartGame(p)
		if err != nil {
			utils.RespWrap(w, http.StatusForbidden, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.Wrap(w, map[string]interface{}{"data": data})
	}
}


func Play(playerSvc *player.Service, gameSvc *game.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		p, err := playerSvc.Find(tk.Email)
		if err != nil {
			utils.RespWrap(w, http.StatusInternalServerError, err.Error())
			return
		}

		if tk.Id != p.Id {
			utils.RespWrap(w, http.StatusForbidden, "player id mismatch")
			return
		}

		jsonMap := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&jsonMap); err != nil {
			utils.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		option, ok := jsonMap["option"]
		if !ok {
			utils.RespWrap(w, http.StatusBadRequest, "option not found")
			return
		}
		pl := &game.Player{Id:p.Id}
		data, err := gameSvc.Play(pl, utils.OptionTypeCast(option.(float64)))
		if err != nil {
			utils.RespWrap(w, http.StatusBadRequest, err.Error())
			return
		}

		artifacts, err := gameSvc.GetArtifacts(pl)
		if err != nil {
			utils.RespWrap(w, http.StatusNotFound, "player not found")
			return
		}

		p.MaxScore = pl.TotalScore
		if err := playerSvc.SaveMaxScore(p); err != nil {
			utils.RespWrap(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.Wrap(w, map[string]interface{}{"data": data, "artifacts": artifacts, "score": pl.TotalScore})
	}
}

func MakeGameHandlers(router *httprouter.Router, playerSvc *player.Service, gameSvc *game.Service) {
	router.HandlerFunc("POST", "/api/bandersnatch/start", middleware.JwtAuth(StartGame(playerSvc, gameSvc)))
	router.HandlerFunc("POST", "/api/bandersnatch/play", middleware.JwtAuth(Play(playerSvc, gameSvc)))
}

