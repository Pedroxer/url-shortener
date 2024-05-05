package service

import (
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"ozon-fintech/internal/models"
	"ozon-fintech/internal/storage"
	mock_storage "ozon-fintech/internal/storage/mocks"
	"testing"
)

func TestGetFullURL(t *testing.T) {
	type mockBehavior func(r *mock_storage.MockDbType, input string)

	testTable := []struct {
		name         string
		input        string
		want         string
		mockBehavior mockBehavior
	}{
		{
			name:  "OK",
			input: "abc_012_yz",
			want:  "https://yandex.ru",
			mockBehavior: func(r *mock_storage.MockDbType, input string) {
				r.EXPECT().GetFullURL(input).Return("https://yandex.ru", nil)
			},
		},
		{
			name:  "ERROR",
			input: "abc_012_yz",
			mockBehavior: func(r *mock_storage.MockDbType, input string) {
				r.EXPECT().GetFullURL(input).Return("", fmt.Errorf("some error"))
			},
		},
		{
			name:  "ERROR_NOT_FOUND",
			input: "abc_012_yz",
			mockBehavior: func(r *mock_storage.MockDbType, input string) {
				r.EXPECT().GetFullURL(input).Return("", sql.ErrNoRows)
			},
		},
	}
	c := gomock.NewController(t)
	defer c.Finish()

	store := mock_storage.NewMockDbType(c)
	service := NewService(store)

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(store, tc.input)

			got, err := service.GetFullURL(tc.input)
			if err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestLoadShortURL(t *testing.T) {
	type mockBehavior func(r *mock_storage.MockDbType, input models.Link)

	testTable := []struct {
		name         string
		input        models.Link
		want         string
		mockBehavior mockBehavior
	}{
		{
			name:  "OK",
			input: models.Link{FullUrl: "https://ya.ru"},
			want:  "TOKEN_3482",
			mockBehavior: func(r *mock_storage.MockDbType, input models.Link) {
				r.EXPECT().LoadShortURL(input).Return("TOKEN_3482", nil)
			},
		},
		{
			name:  "ERROR",
			input: models.Link{FullUrl: "https://ya.ru"},
			mockBehavior: func(r *mock_storage.MockDbType, input models.Link) {
				r.EXPECT().LoadShortURL(input).Return("", fmt.Errorf("some error"))
			},
		},
		{
			name:  "DUPLICATED_ERROR",
			input: models.Link{FullUrl: "https://ya.ru"},
			mockBehavior: func(r *mock_storage.MockDbType, input models.Link) {
				r.EXPECT().LoadShortURL(input).Return("", storage.DuplicateErr)
			},
		},
	}
	c := gomock.NewController(t)
	defer c.Finish()

	store := mock_storage.NewMockDbType(c)
	service := NewService(store)

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(store, tc.input)

			got, err := service.LoadShortURL(tc.input)
			if err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
