package grpc

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/middleware"
)

type Clients struct {
	ArtistClient   *grpc.ClientConn
	AlbumClient    *grpc.ClientConn
	TrackClient    *grpc.ClientConn
	AuthClient     *grpc.ClientConn
	UserClient     *grpc.ClientConn
	PlaylistClient *grpc.ClientConn
}

func requestIdUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	requestId := ctx.Value(middleware.RequestIDKey{}).(string)
	md := metadata.New(map[string]string{
		"request_id": requestId,
	})
	return invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
}

func InitGrpc(cfg *config.Services, logger *zap.SugaredLogger) (*Clients, error) {
	artistClient, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.ArtistService.Host, cfg.ArtistService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(requestIdUnaryClientInterceptor))
	if err != nil {
		logger.Fatal("Error creating artist client:", zap.Error(err))
	}

	albumClient, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.AlbumService.Host, cfg.AlbumService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(requestIdUnaryClientInterceptor))
	if err != nil {
		logger.Fatal("Error creating album client:", zap.Error(err))
	}

	trackClient, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.TrackService.Host, cfg.TrackService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(requestIdUnaryClientInterceptor))
	if err != nil {
		logger.Fatal("Error creating track client:", zap.Error(err))
	}

	playlistClient, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.PlaylistService.Host, cfg.PlaylistService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(requestIdUnaryClientInterceptor))
	if err != nil {
		logger.Fatal("Error creating playlist client:", zap.Error(err))
	}

	authClient, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.AuthService.Host, cfg.AuthService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(requestIdUnaryClientInterceptor))
	if err != nil {
		logger.Fatal("Error creating auth client:", zap.Error(err))
	}

	userClient, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.UserService.Host, cfg.UserService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(requestIdUnaryClientInterceptor))
	if err != nil {
		logger.Fatal("Error creating user client:", zap.Error(err))
	}

	return &Clients{
		ArtistClient:   artistClient,
		AlbumClient:    albumClient,
		TrackClient:    trackClient,
		PlaylistClient: playlistClient,
		AuthClient:     authClient,
		UserClient:     userClient,
	}, nil
}
