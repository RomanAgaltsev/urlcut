package transform

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"
	pb "github.com/RomanAgaltsev/urlcut/pkg/shortener/v1"
)

// PbToIncomingBatchDTO трансформирует батч елементов запроса во входящий батч DTO.
func PbToIncomingBatchDTO(src []*pb.ShortenBatchRequest_ShortenBatchRequestItem) []model.IncomingBatchDTO {
	return nil
}

// OutgoingBatchDTOToPb трансформирует исходящий батч DTO в батч элементов ответа.
func OutgoingBatchDTOToPb(src []model.OutgoingBatchDTO) []*pb.ShortenBatchResponse_ShortenBatchResponseItem {
	return nil
}
