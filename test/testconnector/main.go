package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/qlik-oss/corectl/test/testconnector/qlik_connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := &server{}
	qlik_connect.RegisterConnectorServer(s, srv)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	fmt.Println("Server started", port)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
}

func buildChunk(values [][]string) *qlik_connect.DataChunk {
	chunk := &qlik_connect.DataChunk{
		StringBucket: []string{},
		StringCodes:  []int32{},
		NumberCodes:  []int64{},
	}

	var i int32
	for _, row := range values {
		for _, cell := range row {
			chunk.StringBucket = append(chunk.StringBucket, cell)
			chunk.StringCodes = append(chunk.StringCodes, i)
			chunk.NumberCodes = append(chunk.NumberCodes, -1)
			i++
		}
	}
	return chunk
}
func (s *server) GetData(dataRequest *qlik_connect.DataRequest, stream qlik_connect.Connector_GetDataServer) error {

	if strings.ToLower(dataRequest.Parameters.Statement) == "select a" {
		fields := []*qlik_connect.FieldInfo{{Name: "abcs"}, {Name: "numbers"}}
		s.sendFieldList(fields, stream)

		chunk := buildChunk([][]string{
			{"a", "1"},
			{"b", "2"},
			{"b", "3"},
			{"b", "4"},
			{"c", "5"},
		})
		stream.Send(chunk)
	} else if strings.ToLower(dataRequest.Parameters.Statement) == "select b" {
		fields := []*qlik_connect.FieldInfo{{Name: "xyz"}, {Name: "numbers"}}
		s.sendFieldList(fields, stream)

		chunk := buildChunk([][]string{
			{"x", "5"},
			{"y", "4"},
			{"y", "3"},
			{"y", "2"},
			{"z", "1"},
		})
		stream.Send(chunk)

	}
	return nil
}

func (s *server) sendFieldList(fields []*qlik_connect.FieldInfo, stream qlik_connect.Connector_GetDataServer) {
	// Set header with postgresRowData format
	meta := &qlik_connect.GetDataResponse{FieldInfo: fields, TableName: ""}
	headerMap := make(map[string]string)
	getDataResponseBytes, _ := proto.Marshal(meta)
	headerMap["x-qlik-getdata-bin"] = string(getDataResponseBytes)
	stream.SendHeader(metadata.New(headerMap))
}

func (s *server) GetMetaInfo(ctx context.Context, metaInfoRequest *qlik_connect.MetaInfoRequest) (*qlik_connect.MetaInfo, error) {
	var metaInfo = qlik_connect.MetaInfo{Name: "Test connector", Version: "0.0.0", Developer: "Qlik"}
	return &metaInfo, nil
}
