package usecase_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/usecase"
	domainmock "github.com/topi0247/hitori/domain/mock"
)

func newProfileUsecase(ctrl *gomock.Controller) (*usecase.ProfileUsecase, *domainmock.MockProfileRepository) {
	profileRepository := domainmock.NewMockProfileRepository(ctrl)
	uc := usecase.NewProfileUsecase(profileRepository)
	return uc, profileRepository
}

func TestCreateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, profileRepository := newProfileUsecase(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-uuid").
		Return(nil, domain.ErrNotFound)
	profileRepository.EXPECT().Create(ctx, "user-uuid", "たろう").
		Return(&domain.Profile{AuthUserID: "user-uuid", UserName: "たろう"}, nil)

	out, err := uc.Create(ctx, usecase.CreateProfileInput{
		AuthUserID: "user-uuid",
		UserName:   "たろう",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.UserName != "たろう" {
		t.Errorf("UserName=%s, want たろう", out.UserName)
	}
}

func TestCreateProfile_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, profileRepository := newProfileUsecase(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-uuid").
		Return(&domain.Profile{UserName: "既存"}, nil)

	_, err := uc.Create(ctx, usecase.CreateProfileInput{
		AuthUserID: "user-uuid",
		UserName:   "たろう",
	})
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("wantErr=%v, got=%v", domain.ErrAlreadyExists, err)
	}
}

func TestCreateProfile_InvalidUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, profileRepository := newProfileUsecase(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-uuid").
		Return(nil, domain.ErrNotFound)

	_, err := uc.Create(ctx, usecase.CreateProfileInput{
		AuthUserID: "user-uuid",
		UserName:   "",
	})
	if !errors.Is(err, domain.ErrInvalidUserName) {
		t.Errorf("wantErr=%v, got=%v", domain.ErrInvalidUserName, err)
	}
}

func TestGetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, profileRepository := newProfileUsecase(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-uuid").
		Return(&domain.Profile{UserName: "たろう"}, nil)

	out, err := uc.Get(ctx, "user-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.UserName != "たろう" {
		t.Errorf("UserName=%s, want たろう", out.UserName)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, profileRepository := newProfileUsecase(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-uuid").
		Return(nil, domain.ErrNotFound)

	_, err := uc.Get(ctx, "user-uuid")
	if err == nil {
		t.Error("プロフィールが存在しない場合はエラーを返すべき")
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, profileRepository := newProfileUsecase(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-uuid").
		Return(&domain.Profile{UserName: "たろう"}, nil)
	profileRepository.EXPECT().UpdateUserName(ctx, "user-uuid", "じろう").Return(nil)

	out, err := uc.Update(ctx, usecase.UpdateProfileInput{
		AuthUserID: "user-uuid",
		UserName:   "じろう",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.UserName != "じろう" {
		t.Errorf("UserName=%s, want じろう", out.UserName)
	}
}

func TestUpdateProfile_InvalidUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, profileRepository := newProfileUsecase(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-uuid").
		Return(&domain.Profile{UserName: "たろう"}, nil)

	_, err := uc.Update(ctx, usecase.UpdateProfileInput{
		AuthUserID: "user-uuid",
		UserName:   "あいうえおかきくけこさ", // 11文字
	})
	if !errors.Is(err, domain.ErrInvalidUserName) {
		t.Errorf("wantErr=%v, got=%v", domain.ErrInvalidUserName, err)
	}
}

func TestDeleteProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc, profileRepository := newProfileUsecase(ctrl)
	ctx := context.Background()

	profileRepository.EXPECT().FetchByAuthUserID(ctx, "user-uuid").
		Return(&domain.Profile{UserName: "たろう"}, nil)
	profileRepository.EXPECT().Delete(ctx, "user-uuid").Return(nil)

	if err := uc.Delete(ctx, "user-uuid"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
