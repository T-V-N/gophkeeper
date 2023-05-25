package service_test

import (
	"context"
	"testing"

	"github.com/caarlos0/env/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	cardPB "github.com/T-V-N/gophkeeper/internal/grpc/card"
	mock "github.com/T-V-N/gophkeeper/internal/mocks"
	service "github.com/T-V-N/gophkeeper/internal/service/cardService"
	"github.com/T-V-N/gophkeeper/internal/storage"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

func TestCardService_CreateCard(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response cardPB.CreateCardResponse
	}

	tests := []struct {
		name           string
		request        cardPB.CreateCardRequest
		want           want
		ctx            context.Context
		omitCreateCard bool
	}{
		{
			name: "create log pass normally",
			request: cardPB.CreateCardRequest{
				CardNumberHash: "123",
				ValidUntilHash: "123",
				CVVHash:        "123",
				LastFourDigits: "1337",
			},
			want: want{
				errCode:  codes.OK,
				response: cardPB.CreateCardResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "create log pass unauthorized",
			request: cardPB.CreateCardRequest{
				CardNumberHash: "123",
				ValidUntilHash: "123",
				CVVHash:        "123",
				LastFourDigits: "1337",
			},
			want: want{
				errCode:  codes.Unauthenticated,
				response: cardPB.CreateCardResponse{},
			},
			ctx:            context.Background(),
			omitCreateCard: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	c := mock.NewMockCard(ctrl)

	cApp := app.CardApp{Card: c, Cfg: cfg}
	cService := service.CardService{cardPB.UnimplementedCardServer{}, &cApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.omitCreateCard {
				c.EXPECT().CreateCard(tt.ctx, "1", tt.request.CardNumberHash, tt.request.ValidUntilHash, tt.request.CVVHash, tt.request.LastFourDigits, gomock.Any())
			}
			_, err := cService.CreateCard(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestCardService_UpdateCard(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response cardPB.UpdateCardResponse
	}

	tests := []struct {
		name              string
		request           cardPB.UpdateCardRequest
		existingCard      storage.Card
		existingCardError error
		want              want
		ctx               context.Context
		omitUpdateTN      bool
	}{
		{
			name: "update lp log password normally",
			request: cardPB.UpdateCardRequest{
				Id:             "1",
				CardNumberHash: "asvdsf",
				ValidUntilHash: "asvdsf",
				CVVHash:        "asvdsf",
				LastFourDigits: "asvdsf",
				ForceUpdate:    false,
				PreviousHash:   "sameHash",
				IsDeleted:      false,
			},
			existingCard: storage.Card{
				ID:             "1",
				UID:            "1",
				CardNumberHash: "asvdsf3123",
				ValidUntilHash: "asvdsf",
				CVVHash:        "asvdsf",
				LastFourDigits: "asvdsf",
				EntryHash:      "sameHash",
			},
			want: want{
				errCode:  codes.OK,
				response: cardPB.UpdateCardResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "update lp which doesn't exist",
			request: cardPB.UpdateCardRequest{
				Id:             "10",
				CardNumberHash: "asvdsf",
				ValidUntilHash: "asvdsf",
				CVVHash:        "asvdsf",
				LastFourDigits: "asvdsf",
				ForceUpdate:    false,
				PreviousHash:   "sameHash",
				IsDeleted:      false,
			},
			existingCard: storage.Card{
				ID:             "1",
				UID:            "1",
				CardNumberHash: "asvdsf",
				ValidUntilHash: "asvdsf",
				CVVHash:        "asvdsf",
				LastFourDigits: "asvdsf",
				EntryHash:      "sameHash",
			},
			existingCardError: utils.ErrNotFound,
			want: want{
				errCode:  codes.NotFound,
				response: cardPB.UpdateCardResponse{},
			},
			omitUpdateTN: true,
			ctx:          metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "try update lp which already updated",
			request: cardPB.UpdateCardRequest{
				Id:             "1",
				CardNumberHash: "asvdsf",
				ValidUntilHash: "asvdsf",
				CVVHash:        "asvdsf",
				LastFourDigits: "asvdsf",
				ForceUpdate:    false,
				PreviousHash:   "differentHash",
				IsDeleted:      false,
			},
			existingCard: storage.Card{
				ID:             "1",
				UID:            "1",
				CardNumberHash: "asvdsf",
				ValidUntilHash: "asvdsf",
				CVVHash:        "asvdsf",
				LastFourDigits: "asvdsf",
				EntryHash:      "sameHash",
			},
			existingCardError: nil,
			want: want{
				errCode:  codes.Unavailable,
				response: cardPB.UpdateCardResponse{},
			},
			omitUpdateTN: true,
			ctx:          metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
		{
			name: "try update lp which already updated but force",
			request: cardPB.UpdateCardRequest{
				Id:             "1",
				CardNumberHash: "asvdsf",
				ValidUntilHash: "asvdsf",
				CVVHash:        "asvdsf",
				LastFourDigits: "asvdsf",
				ForceUpdate:    true,
				PreviousHash:   "differentHash",
				IsDeleted:      false,
			},
			existingCard: storage.Card{
				ID:             "1",
				UID:            "1",
				CardNumberHash: "asvdsf",
				ValidUntilHash: "asvdsf",
				CVVHash:        "asvdsf",
				LastFourDigits: "asvdsf",
				EntryHash:      "sameHash",
			},
			existingCardError: nil,
			want: want{
				errCode:  codes.OK,
				response: cardPB.UpdateCardResponse{},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	lp := mock.NewMockCard(ctrl)

	lpApp := app.CardApp{Card: lp, Cfg: cfg}
	tnService := service.CardService{cardPB.UnimplementedCardServer{}, &lpApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp.EXPECT().GetCardByID(tt.ctx, tt.request.Id).Return(&tt.existingCard, tt.existingCardError)
			if !tt.omitUpdateTN {
				lp.EXPECT().UpdateCard(tt.ctx, tt.request.Id, tt.request.CardNumberHash, tt.request.ValidUntilHash, tt.request.CVVHash, tt.request.LastFourDigits, gomock.Any(), tt.request.IsDeleted)
			}
			_, err := tnService.UpdateCard(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}

func TestCardService_ListCard(t *testing.T) {
	type want struct {
		errCode  codes.Code
		response cardPB.ListCardResponse
	}

	tests := []struct {
		name            string
		request         cardPB.ListCardRequest
		existingCards   []storage.Card
		existingCardErr error
		want            want
		ctx             context.Context
		uid             string
	}{
		{
			name:          "list log password normally",
			request:       cardPB.ListCardRequest{},
			existingCards: []storage.Card{{EntryHash: "1", UID: "1"}, {EntryHash: "2", UID: "1"}, {EntryHash: "4", UID: "1"}},
			want: want{
				errCode: codes.OK,
				response: cardPB.ListCardResponse{
					Cards: []*cardPB.CardEntry{{EntryHash: "1"}, {EntryHash: "2"}, {EntryHash: "3"}},
				},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
			uid: "1",
		},
		{
			name: "list log password but remove some existing log passwords",
			request: cardPB.ListCardRequest{
				ExistingHashes: []*cardPB.ExistingCardHash{{Id: "1", EntryHash: "1"}, {Id: "3", EntryHash: "3"}},
			},
			existingCards: []storage.Card{{ID: "1", EntryHash: "1", UID: "1"}, {ID: "2", EntryHash: "2"}, {ID: "3", EntryHash: "3", UID: "1"}, {ID: "1", EntryHash: "4"}},
			want: want{
				errCode: codes.OK,
				response: cardPB.ListCardResponse{
					Cards: []*cardPB.CardEntry{{EntryHash: "1"}, {EntryHash: "3"}},
				},
			},
			ctx: metadata.NewIncomingContext(context.Background(), metadata.Pairs("uid", "1")),
			uid: "1",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	env.Parse(cfg)

	tn := mock.NewMockCard(ctrl)

	tnApp := app.CardApp{Card: tn, Cfg: cfg}
	tnService := service.CardService{cardPB.UnimplementedCardServer{}, &tnApp}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tn.EXPECT().ListCardByUID(tt.ctx, tt.uid).Return(tt.existingCards, nil)
			resp, err := tnService.ListCard(tt.ctx, &tt.request)

			if unwrapError, ok := status.FromError(err); ok {
				assert.Equal(t, unwrapError.Code(), tt.want.errCode)
				assert.Equal(t, len(resp.Cards), len(tt.want.response.Cards))

			} else {
				t.Error("can't parse error, unexpected test result")
			}
		})
	}
}
