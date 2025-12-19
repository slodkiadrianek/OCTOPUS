package controllers

import (
	"context"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
	"github.com/slodkiadrianek/octopus/internal/utils/validation"
)

type routeService interface {
	AddWorkingRoutes(ctx context.Context, routes *[]DTO.CreateRoute, appId string, name string) error
}
type RouteController struct {
	routeService  routeService
	loggerService utils.LoggerService
}

func NewRouteController(routeService routeService, loggerService utils.LoggerService) *RouteController {
	return &RouteController{
		routeService:  routeService,
		loggerService: loggerService,
	}
}

func (rc *RouteController) AddWorkingRoutes(w http.ResponseWriter, r *http.Request) {
	appId, err := request.ReadParam(r, "appId")
	if err != nil {
		rc.loggerService.Error(failedToReadParamFromRequest, r.URL.Path)
		response.SetError(w, r, err)
		return
	}

	body, err := request.ReadBody[DTO.CreateRouteData](r)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	routes := body.Routes
	for i := 0; i < len(routes); i++ {

		if i < len(routes)-1 {

			if len(routes[i].NextRouteBody) > 0 {
				resBody := validation.CheckIsNextRouteBodyInTheBodyAndInTheBodyOfTheNextRoute(routes[i], routes[i+1])
				if !resBody {
					err := models.NewError(400, "Validation", "provided next route body data is malformed, make sure next route body data are in response and in the next route")
					rc.loggerService.Info(err.Error(), routes)
					response.SetError(w, r, err)
					return
				}

			}

			if len(routes[i].NextRouteQuery) > 0 {

				resQuery := validation.CheckIsNextRouteQueryInTheBodyAndInTheQueryOfTheNextRoute(routes[i], routes[i+1])
				if !resQuery {
					err := models.NewError(400, "Validation", "provided next route query data is malformed, make sure next route query data are in response and in the next route")
					rc.loggerService.Info(err.Error(), routes)
					response.SetError(w, r, err)
					return
				}

			}

			if len(routes[i].NextRouteParams) > 0 {

				resParams := validation.CheckIsNextRouteParamsInTheBodyAndInTheParamsOfTheNextRoute(routes[i], routes[i+1])
				if !resParams {
					err := models.NewError(400, "Validation", "provided next route params data is malformed, make sure next route params data are in response and in the next route")
					rc.loggerService.Info(err.Error(), routes)
					response.SetError(w, r, err)
					return
				}

			}

		} else {

			if len(routes[i].NextRouteBody) > 0 {
				err := models.NewError(400, "Validation", "provided next route body data is malformed, make sure next route body data are in response and in the next route")
				rc.loggerService.Info(err.Error(), routes)
				response.SetError(w, r, err)
				return
			}

			if len(routes[i].NextRouteQuery) > 0 {
				err := models.NewError(400, "Validation", "provided next route query data is malformed, make sure next route query data are in response and in the next route")
				rc.loggerService.Info(err.Error(), routes)
				response.SetError(w, r, err)
				return
			}

			if len(routes[i].NextRouteParams) > 0 {
				err := models.NewError(400, "Validation", "provided next route params data is malformed, make sure next route params data are in response and in the next route")
				rc.loggerService.Info(err.Error(), routes)
				response.SetError(w, r, err)
				return
			}
		}
		resParams := request.CheckRouteParams(routes[i])
		if !resParams {
			err := models.NewError(400, "Validation", "provided next route params data is malformed, make sure next route params data are in response and in the next route")
			rc.loggerService.Info(err.Error(), routes)
			response.SetError(w, r, err)
			return
		}

		routes[i].ParentId = i
	}

	err = rc.routeService.AddWorkingRoutes(r.Context(), &routes, appId, body.Name)
	if err != nil {
		response.SetError(w, r, err)
		return
	}

	response.Send(w, 201, map[string]string{})
}
