package repository

import (
	"database/sql"
	. "rest_module/model"
)

type OperationRepository struct {
	Db *DBManager // база данных
}

func InitOperationRepository(db *DBManager) *OperationRepository {
	repo := OperationRepository{}
	repo.Db = db
	return &repo
}

func (repo *OperationRepository) Database() *sql.DB {
	if repo.Db == nil {
		panic("База данных не подключена!")
	}

	return repo.Db.database
}

// Сохранение новой операции в БД
func (repo *OperationRepository) InsertOperation(operation *Operation) (int64, error) {
	insertStmt := `insert into "operations" ("sum_value", "operation_type", "user_id", "account_id") values($1, $2, $3, $4) returning "id"`

	var id int64 = 0
	err := repo.Database().QueryRow(insertStmt, operation.SumValue, operation.OperationType, operation.UserId, operation.AccountId).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Список операций пользователя
func (repo *OperationRepository) GetOperationsByUserId(user_id int64) (*[]Operation, error) {
	selectStmt := `select "id", "sum_value", "operation_type", "account_id" from "operations" where "user_id" = $1`
	rows, err := repo.Database().Query(selectStmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []Operation
	for rows.Next() {
		var id int64
		var sumValue float64
		var operationType string
		var account_id int64

		err = rows.Scan(&id, &sumValue, &operationType, &account_id)
		if err != nil {
			return nil, err
		}

		operation := Operation{
			ID:            id,
			SumValue:      sumValue,
			OperationType: operationType,
			AccountId:     account_id,
			UserId:        user_id,
		}
		operations = append(operations, operation)
	}

	return &operations, nil
}

// Список операций пользователя
func (repo *OperationRepository) GetOperationsByAccountId(user_id, account_id int64) (*[]Operation, error) {
	selectStmt := `select "id", "sum_value", "operation_type" from "operations" where "user_id" = $1 and "account_id" = $2`
	rows, err := repo.Database().Query(selectStmt, user_id, account_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []Operation
	for rows.Next() {
		var id int64
		var sumValue float64
		var operationType string

		err = rows.Scan(&id, &sumValue, &operationType)
		if err != nil {
			return nil, err
		}

		operation := Operation{
			ID:            id,
			SumValue:      sumValue,
			OperationType: operationType,
			AccountId:     account_id,
			UserId:        user_id,
		}
		operations = append(operations, operation)
	}

	return &operations, nil
}
