package handlers

import (
	"io"
	"log"

	pb "gophkeeper/protogen/api/proto/v1"
	pbrpc "gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) BinaryData(stream pb.GophKeeper_BinaryDataServer) error {
	var (
		fileData   []byte
		fileName   string
		totalBytes int64
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Клиент завершил отправку
			return stream.SendAndClose(&pbrpc.BinaryDataResponse{
				FileId:    "generated-file-id",
				SizeBytes: totalBytes,
			})
		}
		if err != nil {
			return err
		}

		switch data := req.Data.(type) {
		case *pbrpc.BinaryDataRequest_Metadata:
			// Первое сообщение содержит метаданные
			fileName = data.Metadata.FileName
			log.Printf("Начало загрузки: %s", fileName)
		case *pbrpc.BinaryDataRequest_Chunk:
			// Последующие сообщения - чанки файла
			fileData = append(fileData, data.Chunk...)
			totalBytes += int64(len(data.Chunk))
		}
	}
}
