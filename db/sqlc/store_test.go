package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T){
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurent transfer
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func ()  {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxInput{
				FromAccountID:	account1.ID,
				ToAccountID:    account2.ID,
				Amount: 		amount,	
			})

			errs <- err
			results <- result

		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		errTemp := <- errs
		require.NoError(t, errTemp)
	
		resultTemp := <- results
		require.NotEmpty(t, resultTemp)
		
		//check transfer 
		transfer := resultTemp.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)		
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, errTransfer:= store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, errTransfer)

		//check from entries
		fromEntry := resultTemp.FromEntry
		require.NotEmpty(t, fromEntry)

		require.Equal(t, -amount, fromEntry.Amount)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, errEntry := store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, errEntry)

		//check to entries
		toEntry := resultTemp.ToEntry
		require.NotEmpty(t, toEntry)

		require.Equal(t, amount, toEntry.Amount)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, errEntry = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, errEntry)

		//fromAccount
		fromAccount := resultTemp.FromAccount
		require.NotEmpty(t, fromAccount.ID)
		require.Equal(t, fromAccount.ID, account1.ID)

		//toAccount
		toAccount := resultTemp.ToAccount
		require.NotEmpty(t, toAccount.ID)
		require.Equal(t, toAccount.ID, account2.ID)

		fmt.Println(">> transaction", fromAccount.Balance, toAccount.Balance)

		//TODO check account balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance - int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T){
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurent transfer
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		// i jika huruf ganjil
		if i % 2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func ()  {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxInput{
				FromAccountID:	fromAccountID,
				ToAccountID:    toAccountID,
				Amount: 		amount,	
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		errTemp := <- errs
		require.NoError(t, errTemp)


	}

	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}