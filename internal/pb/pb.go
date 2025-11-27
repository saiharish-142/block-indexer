// Code generated placeholder. TODO: replace with protoc-gen-go output.
package pb

import (
	"context"

	"google.golang.org/grpc"
)

type BlockSummary struct {
	Number     uint64 `json:"number"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parent_hash"`
	Timestamp  int64  `json:"timestamp"`
}

type TxSummary struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	BlockNumber uint64 `json:"block_number"`
	Status      string `json:"status"`
}

type AddressActivity struct {
	Address string     `json:"address"`
	Tx      *TxSummary `json:"tx"`
}

type BlockRequest struct {
	Hash   string `json:"hash"`
	Number uint64 `json:"number"`
}

type TxRequest struct {
	Hash string `json:"hash"`
}

type AddressRequest struct {
	Address string `json:"address"`
}

// QueryServiceServer mirrors the generated gRPC server interface.
type QueryServiceServer interface {
	GetBlock(context.Context, *BlockRequest) (*BlockSummary, error)
	GetTransaction(context.Context, *TxRequest) (*TxSummary, error)
}

// StreamServiceServer mirrors the generated gRPC streaming interface.
type StreamServiceServer interface {
	StreamHeads(*Empty, StreamService_StreamHeadsServer) error
	StreamAddress(*AddressRequest, StreamService_StreamAddressServer) error
}

// Empty placeholder message.
type Empty struct{}

type StreamService_StreamHeadsServer interface {
	Send(*BlockSummary) error
	grpc.ServerStream
}

type StreamService_StreamAddressServer interface {
	Send(*AddressActivity) error
	grpc.ServerStream
}

// Register functions keep gRPC wiring consistent even as placeholders.
func RegisterQueryServiceServer(s grpc.ServiceRegistrar, srv QueryServiceServer) {}
func RegisterStreamServiceServer(s grpc.ServiceRegistrar, srv StreamServiceServer) {}
