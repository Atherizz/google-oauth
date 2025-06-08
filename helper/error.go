package helper

import "database/sql"

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		if errorRollback != nil {
			panic(errorRollback)
		}
		panic(err)
	} else {
		errComit := tx.Commit()
		if errComit != nil {
			panic(errComit)
		}
	}
}
