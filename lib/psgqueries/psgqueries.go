package psgqueries

const CreateUser string = `
		INSERT INTO shop.users (username,password,balance)
		VALUES ($1,$2,$3)
		RETURNING id;
	`

const GetUserByName string = `
		SELECT 
			id, password, balance 
		FROM 
			shop.users 
		WHERE username=$1;
	`
const UserBalanceById string = `
		SELECT 
			shop.users.balance 
		FROM 
			shop.users 
		WHERE shop.users.ID = $1;
	`

const MerchNameCntById string = `
		SELECT 
			shop.merch.name_type, shop.bucket.cnt 
		FROM 
			shop.users 
		JOIN 
			shop.bucket ON shop.users.ID = shop.bucket.fk_user_id 
		JOIN 
			shop.merch ON shop.bucket.fk_merch_id = shop.merch.id
		WHERE shop.users.ID = $1;
	`
const GetterNameAmountById string = `
		SELECT 
			shop.users.username, shop.operations.amount 
		FROM 
			shop.operations
		JOIN 
			shop.users on shop.users.id = shop.operations.fk_id_sender
		WHERE shop.operations.fk_id_getter = $1;
	`

const SenderNameAmountById string = `
		SELECT 
			shop.users.username, shop.operations.amount 
		FROM 
			shop.operations
		JOIN 
			shop.users on shop.users.id = shop.operations.fk_id_getter
		WHERE shop.operations.fk_id_sender = $1;
	`
const MerchPriceIdByName string = `
		SELECT 
			shop.merch.price, shop.merch.id 
		FROM 
			shop.merch
		WHERE shop.merch.name_type=$1;
	`
const MinusUserBalanceById string = `
		UPDATE 
			shop.users 
		SET 
			balance = balance - $1 
		WHERE id = $2 AND balance >= $1;
	`

const UpdateCntItem string = `
		INSERT INTO shop.bucket (fk_user_id, fk_merch_id, cnt)
		VALUES ($1, $2, 1)
		ON CONFLICT (fk_user_id, fk_merch_id) 
		DO UPDATE SET cnt = shop.bucket.cnt + 1;
	`

const PlusUserBalanceById string = `
		UPDATE 
			shop.users 
		SET
			balance = balance + $1 
		WHERE id = $2;
	`
const UpdateOperationsHistory string = `
		INSERT INTO 
		shop.operations (fk_id_sender,fk_id_getter,amount)
		VALUES ($1,$2,$3)
		ON CONFLICT (fk_id_sender, fk_id_getter) 
		DO UPDATE SET amount = shop.operations.amount + $3;
	`
