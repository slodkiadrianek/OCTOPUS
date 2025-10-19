package controllers

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type RouteController struct {
	RouteService *services.RouteService
	Logger       *utils.Logger
}

func NewRouteController(routeService *services.RouteService, logger *utils.Logger) *RouteController {
	return &RouteController{
		RouteService: routeService,
		Logger:       logger,
	}
}

func (rc *RouteController) AddWorkingRoutes(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ReadBody[[]DTO.CreateRoute](r)
	if err != nil {
		utils.SetError(w, r, err)
	}
	for i := 0; i < len(*body); i++ {
		if i < len(*body)-1 {
			if len((*body)[i].NextRouteBody) > 0 {
				resBody := utils.CheckIsNextRouteBodyInTheBodyAndInTheBodyOfTheNextRoute((*body)[i], (*body)[i+1])
				if !resBody {
					err := models.NewError(400, "Validation", "provided next route body data is malformed, make sure next route body data are in response and in the next route")
					rc.Logger.Info(err.Error(), *body)
					utils.SetError(w, r, err)
					return
				}
			}
			if len((*body)[i].NextRouteQuery) > 0 {
				resQuery := utils.CheckIsNextRouteQueryInTheBodyAndInTheQueryOfTheNextRoute((*body)[i], (*body)[i+1])
				if !resQuery {
					err := models.NewError(400, "Validation", "provided next route query data is malformed, make sure next route query data are in response and in the next route")
					rc.Logger.Info(err.Error(), *body)
					utils.SetError(w, r, err)
					return
				}
			}
			if len((*body)[i].NextRouteParams) > 0 {
				resParams := utils.CheckIsNextRouteParamsInTheBodyAndInTheParamsOfTheNextRoute((*body)[i], (*body)[i+1])
				if !resParams {
					err := models.NewError(400, "Validation", "provided next route params data is malformed, make sure next route params data are in response and in the next route")
					rc.Logger.Info(err.Error(), *body)
					utils.SetError(w, r, err)
					return
				}
			}
		}
		resParams := utils.CheckRouteParams((*body)[i])
		if !resParams {
			err := models.NewError(400, "Validation", "provided next route params data is malformed, make sure next route params data are in response and in the next route")
			rc.Logger.Info(err.Error(), *body)
			utils.SetError(w, r, err)
			return
		}
		(*body)[i].ParentId = i
	}
	_ = rc.RouteService.AddWorkingRoutes(body, 2)
	utils.SendResponse(w, 201, map[string]string{})
}
