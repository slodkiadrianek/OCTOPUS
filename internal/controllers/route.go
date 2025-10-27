package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type routeService interface {
	AddWorkingRoutes(ctx context.Context, routes *[]DTO.CreateRoute, appId string, name string) error
}
type RouteController struct {
	RouteService routeService
	Logger       *utils.Logger
}

func NewRouteController(routeService routeService, logger *utils.Logger) *RouteController {
	return &RouteController{
		RouteService: routeService,
		Logger:       logger,
	}
}

func (rc *RouteController) AddWorkingRoutes(w http.ResponseWriter, r *http.Request) {
	appId, err := utils.ReadParam(r, "appId")
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	body, err := utils.ReadBody[DTO.CreateRouteData](r)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	routes := body.Routes
	for i := 0; i < len(routes); i++ {
		if i < len(routes)-1 {
			if len(routes[i].NextRouteBody) > 0 {
				resBody := utils.CheckIsNextRouteBodyInTheBodyAndInTheBodyOfTheNextRoute(routes[i], routes[i+1])
				if !resBody {
					err := models.NewError(400, "Validation", "provided next route body data is malformed, make sure next route body data are in response and in the next route")
					rc.Logger.Info(err.Error(), routes)
					utils.SetError(w, r, err)
					return
				}
			}
			if len(routes[i].NextRouteQuery) > 0 {
				resQuery := utils.CheckIsNextRouteQueryInTheBodyAndInTheQueryOfTheNextRoute(routes[i], routes[i+1])
				if !resQuery {
					err := models.NewError(400, "Validation", "provided next route query data is malformed, make sure next route query data are in response and in the next route")
					rc.Logger.Info(err.Error(), routes)
					utils.SetError(w, r, err)
					return
				}
			}
			if len(routes[i].NextRouteParams) > 0 {
				resParams := utils.CheckIsNextRouteParamsInTheBodyAndInTheParamsOfTheNextRoute(routes[i], routes[i+1])
				if !resParams {
					err := models.NewError(400, "Validation", "provided next route params data is malformed, make sure next route params data are in response and in the next route")
					rc.Logger.Info(err.Error(), routes)
					utils.SetError(w, r, err)
					return
				}
			}
		} else {
			if len(routes[i].NextRouteBody) > 0 {
				err := models.NewError(400, "Validation", "provided next route body data is malformed, make sure next route body data are in response and in the next route")
				rc.Logger.Info(err.Error(), routes)
				utils.SetError(w, r, err)
				return
			}
			if len(routes[i].NextRouteQuery) > 0 {
				err := models.NewError(400, "Validation", "provided next route query data is malformed, make sure next route query data are in response and in the next route")
				rc.Logger.Info(err.Error(), routes)
				utils.SetError(w, r, err)
				return
			}
			if len(routes[i].NextRouteParams) > 0 {
				err := models.NewError(400, "Validation", "provided next route params data is malformed, make sure next route params data are in response and in the next route")
				rc.Logger.Info(err.Error(), routes)
				utils.SetError(w, r, err)
				return
			}
		}
		resParams := utils.CheckRouteParams(routes[i])
		if !resParams {
			err := models.NewError(400, "Validation", "provided next route params data is malformed, make sure next route params data are in response and in the next route")
			rc.Logger.Info(err.Error(), routes)
			utils.SetError(w, r, err)
			return
		}
		routes[i].ParentId = i
	}
	err = rc.RouteService.AddWorkingRoutes(r.Context(), &routes, appId, body.Name)
	if err != nil {
		utils.SetError(w, r, err)
		return
	}
	utils.SendResponse(w, 201, map[string]string{})
}
