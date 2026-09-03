package server

import (
	"context"
	"fmt"
	"io"
	"main/proto"
	"net"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedThumbnailServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func Start(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterThumbnailServiceServer(grpcServer, NewServer())
	return grpcServer.Serve(lis)
}

func (s *Server) GetThumbnail(ctx context.Context, req *proto.ThumbnailRequest) (*proto.ThumbnailResponse, error) {
	videoURL := req.GetVideoUrl()
	videoID, err := getVideoID(videoURL)
	if err != nil {
		return nil, err
	}
	//todo проверка кэша
	data, err := downloadThumbnail(ctx, videoID)
	if err != nil {
		return nil, err
	}
	return &proto.ThumbnailResponse{Data: data, VideoId: videoID}, nil
}

func downloadThumbnail(ctx context.Context, videoID string) ([]byte, error) {
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", videoID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, thumbnailURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	thumbnail, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return thumbnail, nil
}

func getVideoID(videoURL string) (string, error) {
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", err
	}

	var videoID string

	switch parsedURL.Hostname() {
	case "youtu.be":
		videoID = strings.Trim(parsedURL.Path, "/")

	case "youtube.com", "www.youtube.com":
		switch {
		case parsedURL.Path == "/watch":
			videoID = parsedURL.Query().Get("v")

		case strings.HasPrefix(parsedURL.Path, "/shorts/"):
			videoID = strings.TrimPrefix(parsedURL.Path, "/shorts/")
			videoID = strings.Trim(videoID, "/")

		case strings.HasPrefix(parsedURL.Path, "/embed/"):
			videoID = strings.TrimPrefix(parsedURL.Path, "/embed/")
			videoID = strings.Trim(videoID, "/")

		default:
			return "", fmt.Errorf("unsupported youtube url: %s", videoURL)
		}

	default:
		return "", fmt.Errorf("not a youtube video url: %s", videoURL)
	}

	if videoID == "" {
		return "", fmt.Errorf("video id not found in url: %s", videoURL)
	}

	return videoID, nil
}
