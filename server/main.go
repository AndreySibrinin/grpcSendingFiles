/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"fmt"
	pb "github.com/AndreySibrinin/grspSendingFiles/proto"
	"github.com/AndreySibrinin/grspSendingFiles/server/repo"
	"github.com/AndreySibrinin/grspSendingFiles/server/services/fileservices"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50050, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fileRepo := repo.NewFileRepo("files")
	fileService := fileservices.New(fileRepo)

	s := grpc.NewServer()
	pb.RegisterFileUploadServiceServer(s, fileService)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
