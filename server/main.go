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
	"context"
	"flag"
	"fmt"
	pb "github.com/AndreySibrinin/grspSendingFiles/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedFileUploadServiceServer
}

func (s *server) UploadFile(ctx context.Context, req *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {

	fileContent := req.GetFileContent()
	fileName := req.GetFileName()
	fmt.Println(fileName + ":" + string(fileContent))

	// Do something with the file content and name, e.g. save to disk

	path := filepath.Join("server", fileName)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println("Error:", err)
		return &pb.FileUploadResponse{Success: false, Message: "File didnt upload."}, err
	}

	defer file.Close()

	_, err = file.Write(fileContent)
	if err != nil {
		fmt.Println("Error:", err)
		return &pb.FileUploadResponse{Success: false, Message: "File didnt upload."}, err
	}

	fmt.Println("Successfully wrote to file")

	return &pb.FileUploadResponse{Success: true, Message: "File uploaded successfully."}, nil
}

func (s *server) GetListFiles(req *pb.ListFilesRequest, stream pb.FileUploadService_GetListFilesServer) error {

	entries, err := os.ReadDir("./server")

	if err != nil {
		return err
	}

	for _, entry := range entries {

		info, err := entry.Info()

		if err != nil {
			return err
		}

		stat, ok := info.Sys().(*syscall.Stat_t)

		if !ok {
			return err
		}

		res := &pb.ListFilesResponse{
			FileName:   entry.Name(),
			DateChange: info.ModTime().String(),
			DateCreate: time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec).String(),
		}

		if err := stream.Send(res); err != nil {
			return err
		}
		// 2 second delay to simulate a long running process
		time.Sleep(2 * time.Second)
	}

	return nil
}

func (s *server) DownloadFile(ctx context.Context, req *pb.FileDownloadRequest) (*pb.FileDownloadResponse, error) {

	fileName := req.GetFileName()

	filePath := filepath.Join("server", fileName)

	fileContent, err := ioutil.ReadFile(filePath)

	if err != nil {
		return &pb.FileDownloadResponse{Success: false}, err
	}

	return &pb.FileDownloadResponse{Success: true, FileContent: fileContent}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileUploadServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
