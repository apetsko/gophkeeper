package handlers

import (
	"context"

	pbmodels "github.com/apetsko/gophkeeper/protogen/api/proto/v1/models"
	pbrpc "github.com/apetsko/gophkeeper/protogen/api/proto/v1/rpc"
)

func (s *ServerAdmin) Records(ctx context.Context, in *pbrpc.RecordsRequest) (*pbrpc.RecordsResponse, error) {
	// TODO: Здесь должна быть логика получения записей из БД с учетом пагинации

	records := []*pbmodels.Record{
		{
			Id:   1,
			Type: "login",
			Meta: &pbmodels.Meta{
				Content: "Мой логин от VK",
			},
		},
		{
			Id:   2,
			Type: "card",
			Meta: &pbmodels.Meta{
				Content: "Банковская карта Таджикистана",
			},
		},
		{
			Id:   3,
			Type: "file",
			Meta: &pbmodels.Meta{
				Content: "Файл с выгрузкой транзакций",
			},
		},
	}

	return &pbrpc.RecordsResponse{
		Records: records,
		Count:   int32(len(records)), // TODO: Общее количество записей (может быть больше чем len(records) из-за пагинации)
	}, nil
}
