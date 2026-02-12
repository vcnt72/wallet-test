package service_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	"github.com/vcnt72/go-boilerplate/internal/domain"
	"github.com/vcnt72/go-boilerplate/internal/repository"
	"github.com/vcnt72/go-boilerplate/internal/service"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}

	dsn := os.Getenv("DB_TEST_URL")
	if dsn == "" {
		os.Exit(m.Run())
	}

	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatal(err)
	}
	dbName := strings.TrimPrefix(u.Path, "/")

	ensureDatabase(dsn, dbName)

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	testDB = sqlx.NewDb(sqlDB, "pgx")

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(sqlDB, "../../migrations"); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func ensureDatabase(dsn string, dbName string) {
	adminDSN := strings.Replace(dsn, dbName, "postgres", 1)

	adminDB, err := sql.Open("pgx", adminDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer adminDB.Close()

	var exists bool
	err = adminDB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM pg_database WHERE datname = $1
		)
	`, dbName).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err = adminDB.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func cleanDB(t *testing.T) {
	_, err := testDB.Exec(`
		TRUNCATE TABLE ledgers RESTART IDENTITY CASCADE;
		TRUNCATE TABLE wallets RESTART IDENTITY CASCADE;
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func seedUser(t *testing.T, userID int64) {
	_, err := testDB.Exec(`
		INSERT INTO users (id, name, created_at, updated_at)
		VALUES ($1, 'test', now(), now())
	`, userID)
	require.NoError(t, err)
}

func seedWallet(t *testing.T, userID int64, balance int64) {
	_, err := testDB.Exec(`
		INSERT INTO wallets (user_id, balance, created_at, updated_at)
		VALUES ($1, $2, now(), now())
	`, userID, balance)
	require.NoError(t, err)
}

func getBalance(t *testing.T, userID int64) int64 {
	var b int64
	err := testDB.QueryRowx(`
		SELECT balance FROM wallets WHERE user_id = $1
	`, userID).Scan(&b)
	require.NoError(t, err)
	return b
}

func countLedgers(t *testing.T, key string) int {
	var c int
	err := testDB.QueryRowx(`
		SELECT COUNT(1) FROM ledgers WHERE idempotency_key = $1
	`, key).Scan(&c)
	require.NoError(t, err)
	return c
}

func newWalletService() *service.WalletService {
	walletRepo := repository.NewWalletRepository(testDB)
	ledgerRepo := repository.NewLedgerRepository(testDB)
	txProvider := repository.NewTxProvider(testDB)
	return service.NewWalletService(walletRepo, ledgerRepo, txProvider)
}

func TestIntegration_Withdraw_Success(t *testing.T) {
	cleanDB(t)

	svc := newWalletService()
	userID := int64(1)

	seedUser(t, userID)
	seedWallet(t, userID, 100_000)

	res, err := svc.Withdraw(context.Background(), service.WithdrawWalletSpec{
		UserID:         userID,
		Amount:         30_000,
		IdempotencyKey: "k-success",
	})
	require.NoError(t, err)
	require.Equal(t, int64(70_000), res.Balance)
	require.Equal(t, int64(70_000), getBalance(t, userID))
}

func TestIntegration_Withdraw_InsufficientFunds(t *testing.T) {
	cleanDB(t)

	svc := newWalletService()
	userID := int64(1)

	seedUser(t, userID)
	seedWallet(t, userID, 50_000)

	res, err := svc.Withdraw(context.Background(), service.WithdrawWalletSpec{
		UserID:         userID,
		Amount:         60_000,
		IdempotencyKey: "k-insufficient",
	})
	require.Nil(t, res)
	require.True(t, errors.Is(err, domain.ErrInsufficientFund))
	require.Equal(t, int64(50_000), getBalance(t, userID))
}

func TestIntegration_Withdraw_Concurrent(t *testing.T) {
	cleanDB(t)

	svc := newWalletService()
	userID := int64(1)

	seedUser(t, userID)
	seedWallet(t, userID, 100_000)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	withdraw := func(key string) {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := svc.Withdraw(ctx, service.WithdrawWalletSpec{
			UserID:         userID,
			Amount:         80_000,
			IdempotencyKey: key,
		})
		errCh <- err
	}

	wg.Add(2)
	go withdraw("k-concurrent-1")
	go withdraw("k-concurrent-2")
	wg.Wait()
	close(errCh)

	success := 0
	insufficient := 0

	for err := range errCh {
		if err == nil {
			success++
		} else if errors.Is(err, domain.ErrInsufficientFund) {
			insufficient++
		} else {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	require.Equal(t, 1, success)
	require.Equal(t, 1, insufficient)
	require.Equal(t, int64(20_000), getBalance(t, userID))
}

func TestIntegration_Withdraw_Idempotency_DeductOnce(t *testing.T) {
	cleanDB(t)

	svc := newWalletService()
	userID := int64(1)

	seedUser(t, userID)
	seedWallet(t, userID, 100_000)

	key := "k-idempotent"

	res1, err := svc.Withdraw(context.Background(), service.WithdrawWalletSpec{
		UserID:         userID,
		Amount:         30_000,
		IdempotencyKey: key,
	})
	require.NoError(t, err)
	require.Equal(t, int64(70_000), res1.Balance)
	require.Equal(t, 1, countLedgers(t, key))

	res2, err := svc.Withdraw(context.Background(), service.WithdrawWalletSpec{
		UserID:         userID,
		Amount:         30_000,
		IdempotencyKey: key,
	})
	require.NoError(t, err)
	require.Equal(t, int64(70_000), res2.Balance)
	require.Equal(t, 1, countLedgers(t, key))
	require.Equal(t, int64(70_000), getBalance(t, userID))
}
