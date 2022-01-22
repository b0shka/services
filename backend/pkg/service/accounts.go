package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/korpgoodness/service.git/internal/domain"
	"github.com/korpgoodness/service.git/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountsService struct {
	repo repository.Accounts
}

func NewAccountsService(repo repository.Accounts) *AccountsService {
	return &AccountsService{repo: repo}
}

func RandomInterval() uint8 {
	min := 15
	max := 40
	rand.Seed(time.Now().UnixNano())
	interval := min + rand.Intn(max-min+1)

	return uint8(interval)
}

func (s *AccountsService) Create(ctx context.Context, accountCreate domain.Account) error {
	accountCreate.Interval = RandomInterval()
	accountCreate.Verify = false
	accountCreate.Launch = false
	accountCreate.Status_block = "clean"

	err := s.repo.Create(ctx, accountCreate)
	return err
}

func (s *AccountsService) GetSettings(ctx context.Context, folderID, accountID primitive.ObjectID) (domain.AccountSettings, error) {
	var accountSettings domain.AccountSettings

	account, err := s.repo.GetData(ctx, accountID)
	if err != nil {
		return domain.AccountSettings{}, err
	}
	accountSettings.ID = account.ID
	accountSettings.Name = account.Name
	accountSettings.Phone = account.Phone
	accountSettings.Launch = account.Launch
	accountSettings.Interval = account.Interval
	accountSettings.Status_block = account.Status_block

	folder, err := s.repo.GetFolderByID(ctx, folderID)
	if err != nil {
		return domain.AccountSettings{}, err
	}
	accountSettings.FolderName = folder.Name
	accountSettings.FolderID = folderID
	accountSettings.Chat = folder.Chat

	foldersMove := map[string]primitive.ObjectID{}
	foldersMove_, err := s.repo.GetFolders(ctx)
	if err != nil {
		return domain.AccountSettings{}, err
	}

	for Name, ObjectID := range foldersMove_ {
		if ObjectID != folderID {
			foldersMove[Name] = ObjectID
		}
	}

	accountSettings.FoldersMove = foldersMove
	return accountSettings, nil
}

func (s *AccountsService) UpdateAccount(ctx context.Context, account domain.AccountUpdate) error {
	err := s.repo.UpdateAccount(ctx, account)
	return err
}

func (s *AccountsService) Delete(ctx context.Context, accountID primitive.ObjectID) error {
	err := s.repo.Delete(ctx, accountID)
	return err
}

func (s *AccountsService) GenerateInterval(ctx context.Context, folderID primitive.ObjectID) error {
	err := s.repo.GenerateInterval(ctx, folderID)
	return err
}
