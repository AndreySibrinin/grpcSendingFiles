package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/AndreySibrinin/grspSendingFiles/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	addr = flag.String("addr", "localhost:50050", "the address to connect to")
)

func uploadFile(client pb.FileUploadServiceClient, filePath string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fileName := filepath.Base(filePath)

	_, err = client.UploadFile(ctxTimeout, &pb.FileUploadRequest{FileContent: fileContent, FileName: fileName})
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(client pb.FileUploadServiceClient, fileName string) error {

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	response, err := client.DownloadFile(ctxTimeout, &pb.FileDownloadRequest{FileName: fileName})

	if err != nil {
		return err
	}

	fileContent := response.GetFileContent()

	fmt.Println("File content: " + string(fileContent))

	path := filepath.Join("client", fileName)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println("Error:", err)
	}

	defer file.Close()

	_, err = file.Write(fileContent)

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Successfully wrote to file")

	if err != nil {
		return err
	}

	return nil
}

func getListFiles(client pb.FileUploadServiceClient) {

	log.Printf("Streaming started")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	stream, err := client.GetListFiles(ctxTimeout, &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming %v", err)
		}

		log.Println(message)
	}

	log.Printf("Streaming finished")
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileUploadServiceClient(conn)
	/*
		err = uploadFile(client, "client/Screenshot.png")

		if err != nil {
			fmt.Printf("error upload file: %v", err)
		}

		err = downloadFile(client, "test.txt")

		if err != nil {
			fmt.Printf("error download file: %v", err)
		}



		err = downloadFile(client, "test.txt")

		if err != nil {
			fmt.Printf("error download file: %v", err)
		}
	*/
	err = uploadFile(client, "client/text.txt")

	if err != nil {
		fmt.Printf("error upload file: %v", err)
	}
	getListFiles(client)

}
