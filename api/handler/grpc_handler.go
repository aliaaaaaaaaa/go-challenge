package handler

import (
	"context"
	estag "es/api/proto/src"
	serviceError "es/internal/errors"
	"es/internal/service"
	"github.com/pkg/errors"
)

type GrpcHandler struct {
	segmentTagService *service.EstimationService
}

func NewGrpcHandler(segmentTagService *service.EstimationService) *GrpcHandler {
	return &GrpcHandler{segmentTagService: segmentTagService}
}

func (g *GrpcHandler) EsSaveSegmentTag(ctx context.Context, request *estag.EsRequest) (*estag.EsResponse, error) {
	var response estag.EsResponse
	err := g.segmentTagService.SaveSegmentTagForUser(request.UserId, request.Segment)
	if err != nil {
		err := errors.Wrap(serviceError.InternalError, err.Error())
		return nil, err
	}
	response.Response = "successfully insert new SegmentTag for user"
	return &response, nil
}

func (g *GrpcHandler) Estimation(ctx context.Context, request *estag.EstimationRequest) (*estag.EstimationResponse, error) {
	days, err := g.segmentTagService.GetSegmentTagFor14dLastDays(request.Segment)
	var response estag.EstimationResponse
	response.Estimation = days
	return &response, err
}
