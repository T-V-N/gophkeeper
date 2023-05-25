package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/T-V-N/gophkeeper/internal/app"
	"github.com/T-V-N/gophkeeper/internal/config"
	internalGRPC "github.com/T-V-N/gophkeeper/internal/grpc"
	cpb "github.com/T-V-N/gophkeeper/internal/grpc/card"
	fpb "github.com/T-V-N/gophkeeper/internal/grpc/file"
	lpb "github.com/T-V-N/gophkeeper/internal/grpc/logPassword"
	tnpb "github.com/T-V-N/gophkeeper/internal/grpc/textNote"
	upb "github.com/T-V-N/gophkeeper/internal/grpc/user"

	cardService "github.com/T-V-N/gophkeeper/internal/service/cardService"
	fileService "github.com/T-V-N/gophkeeper/internal/service/fileService"
	logPassService "github.com/T-V-N/gophkeeper/internal/service/logPasswordService"
	textNoteService "github.com/T-V-N/gophkeeper/internal/service/textNoteService"

	userService "github.com/T-V-N/gophkeeper/internal/service/userService"

	"github.com/T-V-N/gophkeeper/internal/helpers"

	confirmService "github.com/T-V-N/gophkeeper/internal/service/confirmService"
	"github.com/go-chi/chi/v5"
	grpc "google.golang.org/grpc"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	cfg, err := config.Init()
	if err != nil {
		sugar.Fatalw("Unable to load config",
			"Error", err,
		)
	}

	emailSender := helpers.InitEmailSender(cfg)

	userApp, err := app.InitUserApp(cfg, sugar, emailSender)
	if err != nil {
		sugar.Fatalw("Unable to init application",
			"Error", err,
		)
	}
	defer userApp.User.Close()

	cardApp, err := app.InitCardApp(cfg, sugar)
	if err != nil {
		sugar.Fatalw("Unable to init application",
			"Error", err,
		)
	}
	defer cardApp.Card.Close()

	textNoteApp, err := app.InitTextNoteApp(cfg, sugar)
	if err != nil {
		sugar.Fatalw("Unable to init application",
			"Error", err,
		)
	}
	defer textNoteApp.TextNote.Close()

	logPassApp, err := app.InitLogPasswordApp(cfg, sugar)
	if err != nil {
		sugar.Fatalw("Unable to init application",
			"Error", err,
		)
	}
	defer textNoteApp.TextNote.Close()

	fileApp, err := app.InitFileApp(cfg, sugar)
	if err != nil {
		sugar.Fatalw("Unable to init application",
			"Error", err,
		)
	}
	defer fileApp.File.Close()

	confHandler := confirmService.InitConfirmationService(userApp)
	router := chi.NewRouter()
	router.Get("/confirm", confHandler.HandleConfirmUser)

	listen, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(internalGRPC.InitAuthInterceptor(cfg)))
	userService := userService.InitUserService(cfg, userApp)
	cardService := cardService.InitCardService(cfg, cardApp)
	logPassService := logPassService.InitLogPasswordService(cfg, logPassApp)
	textNoteService := textNoteService.InittextNoteService(cfg, textNoteApp)
	fileService := fileService.InitFileService(cfg, fileApp)

	cpb.RegisterCardServer(s, cardService)
	lpb.RegisterLogPasswordServer(s, logPassService)
	tnpb.RegisterTextNoteServer(s, textNoteService)
	fpb.RegisterFileServer(s, fileService)
	upb.RegisterUserServer(s, userService)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	gr := &sync.WaitGroup{}

	server := http.Server{
		Handler: router,
		Addr:    cfg.RunAddress,
	}

	gr.Add(1)
	go func() {
		sugar.Infow("Starting confirmation server",
			"Port", cfg.RunAddress,
		)

		server.ListenAndServe()
		gr.Done()
	}()

	gr.Add(1)
	go func() {
		sugar.Infow("Starting rpc",
			"Port", cfg.RPCPort,
		)
		err := s.Serve(listen)

		if err != nil {
			fmt.Println(err)
		}
		gr.Done()
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, stopShutdownCtx := context.WithTimeout(context.Background(), time.Duration(cfg.ContextCancelTimeout)*time.Second)
	defer stopShutdownCtx()

	err = server.Shutdown(shutdownCtx)

	if err != nil {
		sugar.Errorw("Unable to shutdown http server",
			"Error", err,
		)
	}

	s.Stop()

	gr.Wait()
	sugar.Info("Server stopped")
}
