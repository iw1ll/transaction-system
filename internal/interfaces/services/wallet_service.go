package services

import (
	"context"
	"fmt"
	"transaction-system/internal/domain"
)

type WalletService struct {
	walletRepo domain.WalletRepository
	transRepo  domain.TransactionRepository
}

func NewWalletService(
	walletRepo domain.WalletRepository,
	transRepo domain.TransactionRepository,
) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		transRepo:  transRepo,
	}
}

// GetAllWallets возвращает все кошельки
func (s *WalletService) GetAllWallets(ctx context.Context) ([]domain.Wallet, error) {
	wallets, err := s.walletRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}
	return wallets, nil
}

// GetRecentTransactions возвращает последние транзакции
func (s *WalletService) GetRecentTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	transactions, err := s.transRepo.GetRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	return transactions, nil
}

// GetWalletBalance возвращает баланс кошелька
func (s *WalletService) GetWalletBalance(ctx context.Context, address string) (float64, error) {
	wallet, err := s.walletRepo.GetByAddress(ctx, address)
	if err != nil {
		return 0, fmt.Errorf("failed to get wallet balance: %w", err)
	}
	return wallet.Balance, nil
}

// TransferFunds выполняет перевод средств между кошельками
func (s *WalletService) TransferFunds(
	ctx context.Context,
	req domain.TransferRequest,
) error {
	// Валидация суммы перевода
	if req.Amount <= 0 {
		return domain.ErrInvalidAmount
	}

	// Получение информации об отправителе
	sender, err := s.walletRepo.GetByAddress(ctx, req.From)
	if err != nil {
		return fmt.Errorf("sender wallet lookup failed: %w", err)
	}

	// Проверка достаточности средств
	if sender.Balance < req.Amount {
		return domain.ErrInsufficientFunds
	}

	// Выполнение перевода
	if err := s.walletRepo.Transfer(ctx, req.From, req.To, req.Amount); err != nil {
		return fmt.Errorf("funds transfer failed: %w", err)
	}

	// Логирование транзакции
	if err := s.transRepo.Create(ctx, req.From, req.To, req.Amount); err != nil {
		return fmt.Errorf("transaction logging failed: %w", err)
	}

	return nil
}
