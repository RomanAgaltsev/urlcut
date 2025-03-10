package url

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RomanAgaltsev/urlcut/internal/repository"
	pb "github.com/RomanAgaltsev/urlcut/pkg/shortener/v1"
	"github.com/RomanAgaltsev/urlcut/pkg/transform"
)

// Shorten выполняет обработку запроса на сокращение URL.
func (s *ShortenerServer) Shorten(ctx context.Context, request *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	if err := protovalidate.Validate(request); err != nil {
		slog.Info("shorten request validation", "error", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	url, err := s.shortener.Shorten(ctx, request.LongUrl, uuid.New())
	if err != nil && !errors.Is(err, repository.ErrConflict) {
		slog.Info("failed to short URL", "error", err.Error())
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	if errors.Is(err, repository.ErrConflict) {
		return nil, status.Error(codes.AlreadyExists, "URL already exists")
	}

	response := pb.ShortenResponse{
		ShortUrl: url.Short(),
	}

	return &response, nil
}

// ShortenBatch выполняет обработку запроса на сокращение массива URL (батча).
func (s *ShortenerServer) ShortenBatch(ctx context.Context, request *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	if err := protovalidate.Validate(request); err != nil {
		slog.Info("shorten batch request validation", "error", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	batch := transform.PbToIncomingBatchDTO(request.Items)

	batchShortened, err := s.shortener.ShortenBatch(ctx, batch, uuid.New())
	if err != nil && !errors.Is(err, repository.ErrConflict) {
		slog.Info("failed to short URL", "error", err.Error())
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	if errors.Is(err, repository.ErrConflict) {
		return nil, status.Error(codes.AlreadyExists, "URL already exists")
	}

	response := pb.ShortenBatchResponse{
		Items: transform.OutgoingBatchDTOToPb(batchShortened),
	}

	return &response, nil
}

// Expand выполняет обработку запроса на получение оригинального URL.
func (s *ShortenerServer) Expand(ctx context.Context, request *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	if err := protovalidate.Validate(request); err != nil {
		slog.Info("expand request validation", "error", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	url, err := s.shortener.Expand(ctx, request.ShortUrl)
	if err != nil {
		slog.Info("failed to expand URL", "error", err.Error())
		return nil, status.Error(codes.NotFound, "url not found")
	}

	if url.Deleted || len(url.Long) == 0 {
		return nil, status.Error(codes.NotFound, "url not found")
	}

	response := pb.ExpandResponse{
		LongUrl: url.Long,
	}

	return &response, nil
}

// UserUrls выполняет обработку запроса на получение списка всех сохраненных URL пользователя.
func (s *ShortenerServer) UserUrls(ctx context.Context, request *pb.UserUrlsRequest) (*pb.UserUrlsResponse, error) {
	if err := protovalidate.Validate(request); err != nil {
		slog.Info("user urls request validation", "error", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response := pb.UserUrlsResponse{}

	return &response, nil
}

// DeleteUserUrls выполняет обработку запроса на удаление всех сохраненных URL пользователя.
func (s *ShortenerServer) DeleteUserUrls(ctx context.Context, request *pb.DeleteUserUrlsRequest) (*pb.DeleteUserUrlsResponse, error) {
	if err := protovalidate.Validate(request); err != nil {
		slog.Info("delete user urls request validation", "error", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response := pb.DeleteUserUrlsResponse{}

	return &response, nil
}

// Stats выполняет обработку запроса на получении статистики сервиса.
func (s *ShortenerServer) Stats(ctx context.Context, request *pb.StatsRequest) (*pb.StatsResponse, error) {
	if err := protovalidate.Validate(request); err != nil {
		slog.Info("stats request validation", "error", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	stats, err := s.shortener.Stats(ctx)
	if err != nil {
		slog.Info("failed to get stats", "error", err.Error())
		return nil, status.Error(codes.Internal, "please look at logs")
	}

	response := pb.StatsResponse{
		Urls:  int64(stats.Urls),
		Users: int64(stats.Users),
	}

	return &response, nil
}
