package fileservices

import (
	"context"
	"fmt"
	pb "github.com/AndreySibrinin/grspSendingFiles/proto"
)

const (
	maxCallsUpload    = 10
	maxCallsDownload  = 10
	maxCallsListFiles = 100
)

var (
	semUpload    = make(chan struct{}, maxCallsUpload)
	semDownload  = make(chan struct{}, maxCallsDownload)
	semListFiles = make(chan struct{}, maxCallsListFiles)
)

type Service struct {
	pb.UnimplementedFileUploadServiceServer
	repo FileRepo
}

func New(fr FileRepo) *Service {
	return &Service{
		repo: fr,
	}
}

func (s *Service) UploadFile(ctx context.Context, req *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {

	select {
	case <-ctx.Done():
		return &pb.FileUploadResponse{Success: false, Message: "File didnt upload. high load"}, ctx.Err()

	case semUpload <- struct{}{}:
		defer func() { <-semUpload }()

		data := req.GetFileContent()
		fileName := req.GetFileName()
		err := s.repo.UploadFile(fileName, data)

		if err != nil {
			fmt.Println("Error:", err)
			return &pb.FileUploadResponse{Success: false, Message: "File didnt upload."}, err
		}

		return &pb.FileUploadResponse{Success: true, Message: "File uploaded successfully."}, nil
	}
}

func (s *Service) GetListFiles(req *pb.ListFilesRequest, stream pb.FileUploadService_GetListFilesServer) error {

	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case semListFiles <- struct{}{}:
		defer func() { <-semListFiles }()

		entries, err := s.repo.GetListFiles()

		if err != nil {
			return err
		}

		for _, entry := range entries {

			if err != nil {
				return err
			}

			info, err := entry.Info()

			if err != nil {
				return err
			}

			res := &pb.ListFilesResponse{
				FileName:   entry.Name(),
				DateChange: info.ModTime().String(),
				//DateCreate: time.Unix(stat.Atim.Sec, stat.Atim.Nsec).String(),
			}

			if err := stream.Send(res); err != nil {
				return err
			}
		}

		return nil
	}

}

func (s *Service) DownloadFile(ctx context.Context, req *pb.FileDownloadRequest) (*pb.FileDownloadResponse, error) {

	select {
	case <-ctx.Done():
		return &pb.FileDownloadResponse{Success: false}, ctx.Err()
	case semDownload <- struct{}{}:
		defer func() { <-semDownload }()

		fileName := req.GetFileName()

		fileContent, err := s.repo.DownloadFile(fileName)

		if err != nil {
			return &pb.FileDownloadResponse{Success: false}, err
		}

		return &pb.FileDownloadResponse{Success: true, FileContent: fileContent}, nil
	}

}
