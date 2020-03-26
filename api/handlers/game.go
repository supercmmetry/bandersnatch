package handlers

import (
	"bandersnatch/api/middleware"
	"bandersnatch/pkg"
	"bandersnatch/pkg/game"
	"bandersnatch/pkg/player"
	"bandersnatch/utils"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Play(playerSvc *player.Service, gameSvc *game.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tk := ctx.Value(middleware.JwtContextKey("token")).(*middleware.Token)

		p, err := playerSvc.Find(tk.Email)
		if err != nil {
			switch err {
			case pkg.ErrNotFound:
				utils.RespWrap(w, http.StatusNotFound, "player not found")
			default:
				utils.RespWrap(w, http.StatusInternalServerError, err.Error())
			}
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
		pl := &game.Player{Id:p.Id}

		if !gameSvc.CheckIfPlayerExists(pl) {
			data, err := gameSvc.StartGame(p)
			if err != nil {
				utils.RespWrap(w, http.StatusForbidden, err.Error())
				return
			}

			w.WriteHeader(http.StatusOK)
			utils.Wrap(w, map[string]interface{}{"data": data})
			return
		}
		if result, ok := jsonMap["start"]; ok {
			if shouldStart, ok := result.(bool); shouldStart && ok {
				data, err := gameSvc.StartGame(p)
				if err != nil {
					utils.RespWrap(w, http.StatusForbidden, err.Error())
					return
				}

				w.WriteHeader(http.StatusOK)
				utils.Wrap(w, map[string]interface{}{"data": data})
				return
			}
		}


		if result, ok := jsonMap["resume"]; ok {
			if shouldResume, ok := result.(bool); shouldResume && ok {
				data, err := gameSvc.GetNodeData(pl)
				if err != nil {
					utils.RespWrap(w, http.StatusInternalServerError, err.Error())
				} else {
					artifacts, err := gameSvc.GetArtifacts(pl)
					if err != nil {
						utils.RespWrap(w, http.StatusNotFound, "player not found")
						return
					}
					w.WriteHeader(http.StatusOK)
					utils.Wrap(w, map[string]interface{}{"data": data, "artifacts": artifacts, "score": pl.TotalScore})
					return
				}

			}
		}

		option, ok := jsonMap["option"]
		if !ok {
			utils.RespWrap(w, http.StatusBadRequest, "option not found")
			return
		}


		optionFloat, ok := option.(float64)
		if !ok {
			utils.RespWrap(w, http.StatusBadRequest, "option must be of type float")
			return
		}
		data, err := gameSvc.Play(pl, utils.OptionTypeCast(optionFloat))
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
			switch err {
			case pkg.ErrNotFound:
				utils.RespWrap(w, http.StatusNotFound, "player not found")
			default:
				utils.RespWrap(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.Wrap(w, map[string]interface{}{"data": data, "artifacts": artifacts, "score": pl.TotalScore})
	}
}

func MakeGameHandlers(router *httprouter.Router, playerSvc *player.Service, gameSvc *game.Service) {
	router.HandlerFunc("POST", "/api/bandersnatch/play", middleware.JwtAuth(Play(playerSvc, gameSvc)))
}

