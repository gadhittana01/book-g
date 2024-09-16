package mockrepo

import (
	gomock "github.com/golang/mock/gomock"
)

func SetupMockTxPool(ctrl *gomock.Controller, repo *MockRepository, isRollback ...bool) {
	pgxMock := NewMockPGXPool(ctrl)
	mockPool := NewMockPgxIface(ctrl)

	repo.EXPECT().GetDB().Return(pgxMock).Times(1)
	repo.EXPECT().WithTx(gomock.Any()).Return(repo).Times(1)
	pgxMock.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockPool, nil).Times(1)

	if len(isRollback) > 0 && isRollback[0] {
		mockPool.EXPECT().Rollback(gomock.Any()).Return(nil).Times(1)
	} else {
		mockPool.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)
	}
}
