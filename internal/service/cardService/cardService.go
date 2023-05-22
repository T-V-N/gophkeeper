package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	cardPB "github.com/T-V-N/gophkeeper/internal/grpc/card"
	"github.com/T-V-N/gophkeeper/internal/service"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type CardService struct {
	cardPB.UnimplementedCardServer
	CardApp *app.CardApp
}

func (cs *CardService) CreateCard(ctx context.Context, in *cardPB.CreateCardRequest) (*cardPB.CreateCardResponse, error) {
	response := cardPB.CreateCardResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	id, err := cs.CardApp.CreateCard(ctx, uid, in.CardNumberHash, in.ValidUntilHash, in.CVVHash, in.LastFourDigits)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrDBLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		}
	}

	response.Id = id

	return &response, nil
}

func (cs *CardService) UpdateCard(ctx context.Context, in *cardPB.UpdateCardRequest) (*cardPB.UpdateCardResponse, error) {
	response := cardPB.UpdateCardResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unathorized")
	}

	err = cs.CardApp.UpdateCard(ctx, uid, in.Id, in.CardNumberHash, in.ValidUntilHash, in.CVVHash, in.LastFourDigits, in.PreviousHash, in.IsDeleted, in.ForceUpdate)

	if err != nil {
		switch errors.Unwrap(err) {
		case utils.ErrNoData:
			return nil, status.Error(codes.NotFound, err.(utils.WrappedAPIError).Message())
		case utils.ErrConflict:
			return nil, status.Error(codes.Unavailable, "conflict")
		case utils.ErrNotAuthorized:
			return nil, status.Error(codes.Unauthenticated, "unathorized")
		case utils.ErrDBLayer:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		default:
			return nil, status.Error(codes.Internal, err.(utils.WrappedAPIError).Message())
		}
	}

	return &response, nil
}

func (cs *CardService) ListCard(ctx context.Context, in *cardPB.ListCardRequest) (*cardPB.ListCardResponse, error) {
	response := cardPB.ListCardResponse{}

	uid, err := service.ExtractUIDFromCtx(ctx)

	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	existingHashes := []app.ExistingHash{}
	for _, hash := range in.ExistingHashes {
		existingHashes = append(existingHashes, app.ExistingHash{ID: hash.Id, EntryHash: hash.EntryHash})
	}

	cards, err := cs.CardApp.ListCard(ctx, uid, existingHashes)

	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}

	cardPBResponse := []*cardPB.CardEntry{}
	for _, card := range cards {
		cardPBResponse = append(cardPBResponse, &cardPB.CardEntry{
			Id:             card.ID,
			CardNumberHash: card.CardNumberHash,
			ValidUntilHash: card.ValidUntilHash,
			CvvHash:        card.CVVHash,
			LastFourDigits: card.LastFourDigits,
			EntryHash:      card.EntryHash,
		})
	}

	response.Cards = cardPBResponse

	return &response, nil
}

func InitCardService(cfg *config.Config, a *app.CardApp) *CardService {
	return &CardService{cardPB.UnimplementedCardServer{}, a}
}
