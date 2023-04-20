package server

import (
	"net"

	"github.com/jorgepereirajunior/fclx/chatservice/internal/infra/grpc/pb"
	"github.com/jorgepereirajunior/fclx/chatservice/internal/infra/grpc/service"
	"github.com/jorgepereirajunior/fclx/chatservice/internal/usecase/chatcompletionstream"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream            chatcompletionstream.ChatCompletionConfigInputDTO
	ChatService                 service.ChatService
	Port                        string
	AuthToken                   string
	StreamChannel               chan chatcompletionstream.ChatCompletionOutputDTO
}

func NewGRPCServer(chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase, chatConfigStream chatcompletionstream.ChatCompletionConfigInputDTO, port, authToken string, streamChannel chan chatcompletionstream.ChatCompletionOutputDTO) *GRPCServer {
	chatService := service.NewChatService(chatCompletionStreamUseCase, chatConfigStream, streamChannel)
	return &GRPCServer{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		Port:                        port,
		AuthToken:                   authToken,
		StreamChannel:               streamChannel,
		ChatService:                 *chatService,
	}
}

func (g *GRPCServer) Start() {
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, &g.ChatService)

	lis, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		panic(err.Error())
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err.Error())
	}
}
