package repo

import (
	"errors"
	"fmt"
	"root/lib/errs"
	"root/lib/psgqueries"
	"root/models"
	"root/store/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PsgRepo struct {
	pgCl *postgres.Client
}

func NewUserRepo(pgCl *postgres.Client) *PsgRepo {
	return &PsgRepo{
		pgCl: pgCl,
	}
}

func (pr *PsgRepo) CreateUser(user *models.User) error {
	q := psgqueries.CreateUser

	if err := pr.pgCl.DbPool.QueryRow(pr.pgCl.Ctx, q, user.Username, user.Password, user.Balance).Scan(&user.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("Ошибка базы данных: код - %v, сообщение - %v, детали - %v", pgErr.Code, pgErr.Message, pgErr.Detail)
			fmt.Println(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (pr *PsgRepo) GetUserByName(u *models.User) error {
	q := psgqueries.GetUserByName

	if err := pr.pgCl.DbPool.QueryRow(pr.pgCl.Ctx, q, u.Username).Scan(&u.ID, &u.Password, &u.Balance); err != nil {
		return err
	}

	return nil
}

func (pr *PsgRepo) GetAllActivity(id int) (*models.UserActivity, error) {
	const itemsInstore = 12
	ua := models.UserActivity{
		Coins:       0,
		Inventory:   make([]*models.Inventory, 0, itemsInstore),
		CoinHistory: &models.CoinHistory{},
	}
	qBal := psgqueries.UserBalanceById

	if err := pr.pgCl.DbPool.QueryRow(pr.pgCl.Ctx, qBal, id).Scan(&ua.Coins); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("Ошибка базы данных: код - %v, сообщение - %v, детали - %v", pgErr.Code, pgErr.Message, pgErr.Detail)
			fmt.Println(newErr)
			return nil, newErr
		}
		return nil, err
	}

	qInvent := psgqueries.MerchNameCntById

	rows, err := pr.pgCl.DbPool.Query(pr.pgCl.Ctx, qInvent, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	invs := make([]*models.Inventory, 0, itemsInstore)
	for rows.Next() {
		var inv models.Inventory

		if err := rows.Scan(&inv.Type, &inv.Quantity); err != nil {
			return nil, err
		}

		invs = append(invs, &inv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	qFrom := psgqueries.GetterNameAmountById

	rows, err = pr.pgCl.DbPool.Query(pr.pgCl.Ctx, qFrom, id)
	if err != nil {
		return nil, err
	}

	getFrom := make([]*models.Received, 0)

	for rows.Next() {
		var gf models.Received

		if err := rows.Scan(&gf.FromUser, &gf.Amount); err != nil {
			return nil, err
		}

		getFrom = append(getFrom, &gf)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	qTo := psgqueries.SenderNameAmountById

	rows, err = pr.pgCl.DbPool.Query(pr.pgCl.Ctx, qTo, id)
	if err != nil {
		return nil, err
	}
	sendTo := make([]*models.Sent, 0)

	for rows.Next() {
		var st models.Sent

		if err := rows.Scan(&st.ToUser, &st.Amount); err != nil {
			return nil, err
		}

		sendTo = append(sendTo, &st)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	ua.Inventory = invs
	ua.CoinHistory.Received = getFrom
	ua.CoinHistory.Sent = sendTo

	return &ua, nil
}

func (pr *PsgRepo) BuyItem(id int, item string) error {

	tx, err := pr.pgCl.DbPool.BeginTx(pr.pgCl.Ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})

	if err != nil {
		return err
	}
	defer tx.Rollback(pr.pgCl.Ctx)

	q := psgqueries.MerchPriceIdByName

	var price int
	var idMerch int

	if err := tx.QueryRow(pr.pgCl.Ctx, q, item).Scan(&price, &idMerch); err != nil {
		return err
	}

	q = psgqueries.MinusUserBalanceById

	if pgComm, err := tx.Exec(pr.pgCl.Ctx, q, price, id); err != nil || pgComm.RowsAffected() == 0 {
		return errs.NotEnoughMoney
	}

	q = psgqueries.UpdateCntItem

	ct, err := tx.Exec(pr.pgCl.Ctx, q, id, idMerch)

	if err != nil || ct.RowsAffected() == 0 {
		return fmt.Errorf("can't update table bucket: %v", err)
	}

	err = tx.Commit(pr.pgCl.Ctx)

	if err != nil {
		return fmt.Errorf("can't commit transaction: %v", err)
	}

	return nil
}

func (pr *PsgRepo) SendCoins(fromId int, toId int, amount int) error {

	tx, err := pr.pgCl.DbPool.BeginTx(pr.pgCl.Ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})

	if err != nil {
		return err
	}
	defer tx.Rollback(pr.pgCl.Ctx)

	q := psgqueries.MinusUserBalanceById

	if cmd, err := tx.Exec(pr.pgCl.Ctx, q, amount, fromId); err != nil || cmd.RowsAffected() == 0 {
		return errs.NotEnoughMoney
	}

	q = psgqueries.PlusUserBalanceById

	if cmd, err := tx.Exec(pr.pgCl.Ctx, q, amount, toId); err != nil || cmd.RowsAffected() == 0 {
		return fmt.Errorf("can't get money %v, user_id -%v", err, toId)
	}

	q = psgqueries.UpdateOperationsHistory

	if cmd, err := tx.Exec(pr.pgCl.Ctx, q, fromId, toId, amount); err != nil || cmd.RowsAffected() == 0 {
		return fmt.Errorf("can't memory operation %v", err)
	}

	if err := tx.Commit(pr.pgCl.Ctx); err != nil {
		return fmt.Errorf("can't commit transaction of sendCoins: %v", err)
	}

	return nil
}
