package basis

import (
	"context"
	"gorono/internal/middlas"
	"gorono/internal/models"
)

func (suite *TstBase) Test00InitDBStorage() {
	tests := []struct {
		name       string
		ctx        context.Context
		dbEndPoint string
		wantErr    bool
	}{
		{
			name:       "InitDB Bad BASE",
			ctx:        context.Background(),
			dbEndPoint: "postgres://postgres:passwordas@localhost:5432/hzwhatbase",
			wantErr:    true,
		},
		{
			name:       "Bad PASSWORD",
			ctx:        context.Background(),
			dbEndPoint: "postgres://postgres:wrongpassword@localhost:5432/forgo",
			wantErr:    true,
		},
		{
			name:       "InitDB Nice manner", // last - RIGHT base params. чтобы база была открыта для дальнейших тестов
			ctx:        context.Background(),
			dbEndPoint: "postgres://postgres:passwordas@localhost:5432/forgo",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			d, err := InitDBStorage(tt.ctx, tt.dbEndPoint)
			suite.dataBase = d
			suite.Require().Equal(err != nil, tt.wantErr) //
			if err == nil {
				suite.dataBase.TablesDrop(tt.ctx)
				err = suite.dataBase.DB.Close(tt.ctx)
				suite.Require().NoError(err)
			}
			d, err = InitDBStorage(tt.ctx, tt.dbEndPoint) // reopen after drop
			suite.dataBase = d
			suite.Require().Equal(err != nil, tt.wantErr) //
		})
	}
}

func (suite *TstBase) TestDBstruct_PutMetric() {

	tests := []struct {
		name    string
		metr    Metrics
		wantErr bool
	}{
		{
			name:    "Nice gauge PUT metric",
			metr:    models.Metrics{MType: "gauge", ID: "Alloc", Value: middlas.Ptr(777.77)},
			wantErr: false,
		},
		// {
		// 	name:    "Nice counter PUT metric",
		// 	metr:    models.Metrics{MType: "counter", ID: "coooo", Delta: middlas.Ptr[int64](777)},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.dataBase.PutMetric(suite.ctx, &tt.metr, nil)
			suite.Require().Equal(err != nil, tt.wantErr)
		})
	}
}
