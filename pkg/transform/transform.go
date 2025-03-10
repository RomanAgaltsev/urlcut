package transform

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"
	pb "github.com/RomanAgaltsev/urlcut/pkg/shortener/v1"
)

func PbToIncomingBatchDTO(src []*pb.ShortenBatchRequest_ShortenBatchRequestItem) []model.IncomingBatchDTO {
	return nil
}

func OutgoingBatchDTOToPb(src []model.OutgoingBatchDTO) []*pb.ShortenBatchResponse_ShortenBatchResponseItem {
	return nil
}
