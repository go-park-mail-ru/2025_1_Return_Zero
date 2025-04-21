package grpc

import (
	"fmt"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	ArtistClient *grpc.ClientConn
}

func InitGrpc(cfg *config.Services, logger *zap.SugaredLogger) (*Clients, error) {
	artistClient, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.ArtistService.Host, cfg.ArtistService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Error creating artist client:", zap.Error(err))
	}

	return &Clients{
		ArtistClient: artistClient,
	}, nil
}
