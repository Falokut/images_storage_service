package handler_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Falokut/images_storage_service/internal/handler"
	"github.com/Falokut/images_storage_service/internal/models"
	mock_repository "github.com/Falokut/images_storage_service/internal/repository/mocks"
	"github.com/Falokut/images_storage_service/internal/service"
	service_mock "github.com/Falokut/images_storage_service/internal/service/mocks"
	"github.com/Falokut/images_storage_service/pkg/images_storage_service/v1/protos"
	"github.com/Falokut/images_storage_service/pkg/logging"
	"github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const testResoursesDir = "test/resources/"
const testPNGImageName = "test.png"
const testJPGImageName = "test.jpg"

type replaceMockBehavior func(s *mock_repository.MockImageStorage, img []byte, filename string, relativePath string)
type isExistMockBehavior func(s *mock_repository.MockImageStorage, filename string, relativePath string)
type saveMockBehavior func(s *mock_repository.MockImageStorage, image []byte, relativePath string)
type getMockBehavior func(s *mock_repository.MockImageStorage, imageID string, relativePath string)
type deleteMockBehavior func(s *mock_repository.MockImageStorage, imageID string, relativePath string)

func newServer(t *testing.T, register func(srv *grpc.Server)) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	register(srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	return conn
}

func newClient(t *testing.T, s *handler.ImagesStorageServiceHandler) *grpc.ClientConn {
	return newServer(t, func(srv *grpc.Server) { protos.RegisterImagesStorageServiceV1Server(srv, s) })
}

func TestGetImage(t *testing.T) {
	testCases := []struct {
		ImageID        string
		imageBody      []byte
		Category       string
		mockBehavior   getMockBehavior
		expectedStatus codes.Code
		caseMessage    string
	}{
		{
			ImageID:   uuid.NewString(),
			imageBody: []byte("203 123 212 121"),
			Category:  "asweqeqw",
			mockBehavior: func(s *mock_repository.MockImageStorage,
				imageID string, relativePath string) {
				s.EXPECT().
					GetImage(gomock.Any(), imageID, relativePath).
					Return([]byte("203 123 212 121"), nil).
					Times(1)
			},
			expectedStatus: codes.OK,
			caseMessage:    "Case num %d, check repository return valid data",
		},
		{
			ImageID:   uuid.NewString(),
			imageBody: []byte("123 121 21 21 99 12"),
			mockBehavior: func(s *mock_repository.MockImageStorage,
				imageID string, relativePath string) {
				s.EXPECT().
					GetImage(gomock.Any(), imageID, relativePath).
					Return([]byte{}, models.Error(models.NotFound, "")).
					Times(1)
			},
			expectedStatus: codes.NotFound,
			caseMessage:    "Case num %d, check repository return not found error",
		},
	}

	logger := logging.GetNullLogger()
	logger.Logger.SetLevel(logrus.ErrorLevel)
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		metr := service_mock.NewMockMetrics(mockController)
		metr.EXPECT().IncBytesUploaded(gomock.Any()).AnyTimes()

		repo := mock_repository.NewMockImageStorage(mockController)
		service := service.NewImagesStorageService(logger.Logger, metr, repo, service.Config{})
		conn := newClient(t, handler.NewImagesStorageServiceHandler(logger.Logger,
			handler.Config{}, service))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.mockBehavior(repo, testCase.ImageID, testCase.Category)

		req := &protos.ImageRequest{Category: testCase.Category, ImageId: testCase.ImageID}
		res, err := client.GetImage(ctx, req)

		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)

		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res, caseMessage, "Response mustn't be null")
			assert.Equal(
				t,
				testCase.imageBody,
				res.Data,
				caseMessage,
				"handler mustn't change image data from repository",
			)
		}
	}
}

func TestUploadImage(t *testing.T) {
	imagePNG, err := os.ReadFile(filepath.Clean(testResoursesDir + testPNGImageName))
	assert.NoError(t, err)
	imageJPG, err := os.ReadFile(filepath.Clean(testResoursesDir + testJPGImageName))
	assert.NoError(t, err)

	testCases := []struct {
		imageBody      []byte
		Category       string
		mockBehavior   saveMockBehavior
		expectedStatus codes.Code
		maxImageSize   int
		caseMessage    string
	}{
		{
			imageBody: []byte("31231 12312"),
			Category:  "Test",
			mockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			maxImageSize:   100,
			caseMessage:    "Case num %d, check non image []byte",
		},
		{
			imageBody: []byte{},
			Category:  "Test",
			mockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			maxImageSize:   1,
			caseMessage:    "Case num %d, check receiving image with zero size",
		},
		{
			imageBody: imagePNG,
			mockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			expectedStatus: codes.OK,
			maxImageSize:   len(imagePNG),
			caseMessage:    "Case num %d, check receiving valid PNG image",
		},
		{
			imageBody: imageJPG,
			mockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			expectedStatus: codes.OK,
			maxImageSize:   len(imagePNG),
			caseMessage:    "Case num %d, check receiving valid JPG image",
		},
		{
			imageBody: imagePNG,
			Category:  "Te132st",
			mockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			maxImageSize:   len(imagePNG) / 2,
			caseMessage:    "Case num %d, check receiving image with size bigger than maxSize",
		},
		{
			imageBody: imagePNG,
			Category:  "Te231st",
			mockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().
					SaveImage(gomock.Any(), image, gomock.Any(), relativePath).
					Return(models.Error(models.Internal, "")).
					Times(1)
			},
			expectedStatus: codes.Internal,
			maxImageSize:   len(imagePNG),
			caseMessage:    "Case num %d, check receiving valid image, but repository return error",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)
		metr := service_mock.NewMockMetrics(mockController)
		metr.EXPECT().IncBytesUploaded(gomock.Any()).AnyTimes()

		service := service.NewImagesStorageService(logger.Logger, metr,
			repo, service.Config{MaxImageSize: testCase.maxImageSize})
		conn := newClient(t, handler.NewImagesStorageServiceHandler(logger.Logger,
			handler.Config{MaxImageSize: testCase.maxImageSize}, service))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.mockBehavior(repo, testCase.imageBody, testCase.Category)

		req := &protos.UploadImageRequest{Category: testCase.Category, Image: testCase.imageBody}
		res, err := client.UploadImage(ctx, req)
		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)

		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res.ImageId, caseMessage, "must return valid image id")
			continue
		}

	}
}

func TestDeleteImage(t *testing.T) {
	testCases := []struct {
		Category       string
		imageID        string
		mockBehavior   deleteMockBehavior
		expectedStatus codes.Code
		caseMessage    string
	}{
		{
			Category: "Test/Any/Category",
			imageID:  "AnyId.png",
			mockBehavior: func(s *mock_repository.MockImageStorage, imageID, relativePath string) {
				s.EXPECT().DeleteImage(gomock.Any(), imageID, relativePath).Return(nil).Times(1)
			},
			expectedStatus: codes.OK,
			caseMessage:    "Case num %d, checks if method work when no errors",
		},
		{
			Category: "Test/Category",
			imageID:  "23910312.png",
			mockBehavior: func(s *mock_repository.MockImageStorage, imageID, relativePath string) {
				s.EXPECT().DeleteImage(gomock.Any(), imageID, relativePath).Return(models.Error(models.Internal, "")).Times(1)
			},
			expectedStatus: codes.Internal,
			caseMessage:    "Case num %d, check receiving valid image, but repository return error",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)
		metr := service_mock.NewMockMetrics(mockController)
		metr.EXPECT().IncBytesUploaded(gomock.Any()).Times(0)

		service := service.NewImagesStorageService(logger.Logger, metr, repo, service.Config{})
		conn := newClient(t, handler.NewImagesStorageServiceHandler(logger.Logger,
			handler.Config{}, service))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.mockBehavior(repo, testCase.imageID, testCase.Category)

		req := &protos.ImageRequest{Category: testCase.Category, ImageId: testCase.imageID}
		_, err := client.DeleteImage(ctx, req)

		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)
	}
}

func TestReplaceImage(t *testing.T) {
	imagePNG, err := os.ReadFile(filepath.Clean(testResoursesDir + testPNGImageName))
	assert.NoError(t, err)
	imageJPG, err := os.ReadFile(filepath.Clean(testResoursesDir + testJPGImageName))
	assert.NoError(t, err)

	testCases := []struct {
		imageBody           []byte
		Category            string
		replaceMockBehavior replaceMockBehavior
		isExistMockBehavior isExistMockBehavior
		saveMockBehavior    saveMockBehavior
		ImageID             string
		CreateIfNotExist    bool
		expectedStatus      codes.Code
		maxImageSize        int
		caseMessage         string
	}{

		{
			imageBody: []byte("31231 12312"),
			Category:  "Test",
			replaceMockBehavior: func(s *mock_repository.MockImageStorage,
				img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), img, gomock.Any(), relativePath).Times(0)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage,
				filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Times(0)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			maxImageSize:   100,
			ImageID:        "3212dada",
			caseMessage:    "Case num %d, check non image []byte",
		},
		{
			imageBody: imagePNG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage,
				img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), img, filename, relativePath).Return(nil).Times(1)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage,
				filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(true, nil).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus: codes.OK,
			ImageID:        "321312312sadqweq",
			maxImageSize:   len(imagePNG),
			caseMessage:    "Case num %d, check receiving valid image with existing image file in repo",
		},
		{
			imageBody: imageJPG,
			Category:  "Test",
			replaceMockBehavior: func(s *mock_repository.MockImageStorage,
				img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), img, filename, relativePath).Times(0)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage,
				filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(false, nil).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus:   codes.NotFound,
			CreateIfNotExist: false,
			maxImageSize:     len(imagePNG),
			ImageID:          "3212dada",
			caseMessage:      "Case num %d, check receiving valid image without existing image file in repo and CreateIfNotExist=false",
		},
		{
			imageBody: imageJPG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage,
				img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), img, filename, relativePath).Times(0)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage,
				filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(false, nil).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image,
					gomock.Any(), relativePath).Return(models.Error(models.Internal, "")).Times(1)
			},
			expectedStatus:   codes.Internal,
			ImageID:          "321312312sadqweq",
			maxImageSize:     len(imagePNG),
			CreateIfNotExist: true,
			caseMessage:      "Case num %d, check receiving valid image without existing image file in repo and CreateIfNotExist=true",
		},
		{
			imageBody: imagePNG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage,
				img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), img, filename, relativePath).Times(0)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage,
				filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(false, nil).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			expectedStatus:   codes.OK,
			ImageID:          "321312312sadqweq",
			maxImageSize:     len(imagePNG),
			CreateIfNotExist: true,
			caseMessage:      "Case num %d, check receiving valid image without existing image file in repo and CreateIfNotExist=true",
		},
		{
			imageBody: imagePNG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage,
				img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), img, filename, relativePath).Return(nil).Times(1)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage,
				filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(true, nil).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus:   codes.OK,
			ImageID:          "321312312sadqweq",
			maxImageSize:     len(imagePNG),
			CreateIfNotExist: true,
			caseMessage:      "Case num %d, check receiving valid image with existing image file in repo",
		},
		{
			imageBody: imagePNG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage,
				img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), img, filename, relativePath).Return(models.Error(models.Internal, "")).Times(1)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage,
				filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(true, nil).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage,
				image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus:   codes.Internal,
			ImageID:          "321312312sadqweq",
			maxImageSize:     len(imagePNG),
			CreateIfNotExist: true,
			caseMessage:      "Case num %d, check receiving valid image with existing image file in repo, but with error while rewriting",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)
		metr := service_mock.NewMockMetrics(mockController)
		metr.EXPECT().IncBytesUploaded(gomock.Any()).AnyTimes()

		service := service.NewImagesStorageService(logger.Logger,
			metr, repo, service.Config{MaxImageSize: testCase.maxImageSize})
		conn := newClient(t, handler.NewImagesStorageServiceHandler(logger.Logger,
			handler.Config{MaxImageSize: testCase.maxImageSize}, service))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.replaceMockBehavior(repo, testCase.imageBody, testCase.ImageID, testCase.Category)
		testCase.isExistMockBehavior(repo, testCase.ImageID, testCase.Category)
		testCase.saveMockBehavior(repo, testCase.imageBody, testCase.Category)

		req := &protos.ReplaceImageRequest{
			Category:         testCase.Category,
			ImageId:          testCase.ImageID,
			ImageData:        testCase.imageBody,
			CreateIfNotExist: testCase.CreateIfNotExist,
		}
		res, err := client.ReplaceImage(ctx, req)
		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)
		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res, caseMessage)
			continue
		}
	}
}

func TestIsImageExist(t *testing.T) {
	testCases := []struct {
		Category         string
		imageID          string
		mockBehavior     isExistMockBehavior
		caseMessage      string
		expectedResponse bool
	}{
		{
			imageID:  "AnyID",
			Category: "AnyPAth/Patrh/ase1",
			mockBehavior: func(s *mock_repository.MockImageStorage, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(true, nil).Times(1)
			},
			expectedResponse: true,
			caseMessage:      "Case num %d, checks response, if image exist",
		},
		{
			imageID:  "AnyID",
			Category: "AnyPAth/1231asdweqq/ase1",
			mockBehavior: func(s *mock_repository.MockImageStorage, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(false, nil).Times(1)
			},
			expectedResponse: false,
			caseMessage:      "Case num %d, checks response, if image not exist",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)
		metr := service_mock.NewMockMetrics(mockController)
		metr.EXPECT().IncBytesUploaded(gomock.Any()).AnyTimes()

		service := service.NewImagesStorageService(logger.Logger, metr, repo, service.Config{})
		conn := newClient(t, handler.NewImagesStorageServiceHandler(logger.Logger,
			handler.Config{}, service))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.mockBehavior(repo, testCase.imageID, testCase.Category)

		req := &protos.ImageRequest{Category: testCase.Category, ImageId: testCase.imageID}
		res, err := client.IsImageExist(ctx, req)

		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			codes.OK,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)

		assert.NotNil(t, res, caseMessage, "Must return valid response")
		assert.Equal(t, testCase.expectedResponse, res.ImageExist, caseMessage, "Must return expected response")
		assert.NoError(t, err, caseMessage, "Mustn't return error")
	}
}

func TestStreamingUploadImage(t *testing.T) {
	imagePNG, err := os.ReadFile(filepath.Clean(testResoursesDir + testPNGImageName))
	assert.NoError(t, err)
	imageJPG, err := os.ReadFile(filepath.Clean(testResoursesDir + testJPGImageName))
	assert.NoError(t, err)

	testCases := []struct {
		imageBody      []byte
		Category       string
		mockBehavior   saveMockBehavior
		expectedStatus codes.Code
		maxImageSize   int
		chunkSize      int
		caseMessage    string
		cancelContext  bool
	}{
		{
			imageBody: []byte("31231 12312"),
			Category:  "Test",
			mockBehavior: func(s *mock_repository.MockImageStorage, image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			maxImageSize:   100,
			caseMessage:    "Case num %d, check non image []byte",
		},
		{
			imageBody: []byte{},
			Category:  "Test",
			mockBehavior: func(s *mock_repository.MockImageStorage, image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			maxImageSize:   1,
			caseMessage:    "Case num %d, check receiving image with zero size",
		},
		{
			imageBody: imagePNG,
			mockBehavior: func(s *mock_repository.MockImageStorage, image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			expectedStatus: codes.OK,
			maxImageSize:   len(imagePNG) + 100,
			chunkSize:      len(imagePNG) / 16,
			caseMessage:    "Case num %d, check receiving valid image",
		},
		{
			imageBody: imageJPG,
			Category:  "Te132st",
			mockBehavior: func(s *mock_repository.MockImageStorage, image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			maxImageSize:   len(imageJPG) / 2,
			chunkSize:      len(imageJPG) / 10,
			caseMessage:    "Case num %d, check receiving image with size bigger than maxSize",
		},
		{
			imageBody: imagePNG,
			Category:  "Te231st",
			mockBehavior: func(s *mock_repository.MockImageStorage, image []byte, relativePath string) {
				s.EXPECT().
					SaveImage(gomock.Any(), image, gomock.Any(), relativePath).
					Return(models.Error(models.Internal, "")).
					Times(1)
			},
			expectedStatus: codes.Internal,
			maxImageSize:   len(imagePNG),
			caseMessage:    "Case num %d, check receiving valid image, but repository return error",
		},
		{
			imageBody: imageJPG,
			mockBehavior: func(s *mock_repository.MockImageStorage, image []byte, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.Canceled,
			maxImageSize:   len(imagePNG),
			chunkSize:      60,
			cancelContext:  true,
			caseMessage:    "Case num %d, check receiving valid image with cancel",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)
		metr := service_mock.NewMockMetrics(mockController)
		metr.EXPECT().IncBytesUploaded(gomock.Any()).AnyTimes()

		service := service.NewImagesStorageService(logger.Logger, metr, repo, service.Config{MaxImageSize: testCase.maxImageSize})
		conn := newClient(t, handler.NewImagesStorageServiceHandler(logger.Logger,
			handler.Config{MaxImageSize: testCase.maxImageSize}, service))

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		testCase.mockBehavior(repo, testCase.imageBody, testCase.Category)

		streamingReq, err := client.StreamingUploadImage(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, streamingReq)
		if testCase.chunkSize == 0 {
			testCase.chunkSize = 20
		}

		for j := 0; j < len(testCase.imageBody); j += testCase.chunkSize {
			last := j + testCase.chunkSize
			if last > len(testCase.imageBody) {
				last = len(testCase.imageBody)
			}
			var chunk []byte
			chunk = append(chunk, testCase.imageBody[j:last]...)
			req := &protos.StreamingUploadImageRequest{Category: testCase.Category, Data: chunk}
			err := streamingReq.Send(req)
			if !errors.Is(err, io.EOF) {
				assert.NoError(t, err)
			}
			if testCase.cancelContext {
				cancel()
			}
		}

		res, err := streamingReq.CloseAndRecv()
		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)
		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res, caseMessage, "Expected Response not to be nil")
			assert.NoError(t, err, "Mustn't return error when returning non nil Response")
			continue
		}
	}
}
