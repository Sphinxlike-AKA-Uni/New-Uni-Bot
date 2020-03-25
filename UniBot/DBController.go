// Script resposible for controlling the database
package Uni
import (
	"fmt"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/lib/pq"
	"database/sql"
)

// Open The Database
func (Uni *UniBot) OpenDB() (err error) {
	if !IsStringInArray(Uni.Config.DBDriver, availabledrivers) {
		return errors.New(fmt.Sprintf("Driver \"%s\" is not present inside available drivers", Uni.Config.DBDriver))
	}
	Uni.DB, err = sql.Open(Uni.Config.DBDriver, Uni.Config.DBContent)
	if err != nil {
		goto End
	}
	
	Uni.DBExec("Startup")
	End: return
}

// Execute instruction for database (sometimes different SQL servers have different syntax)
func (Uni *UniBot) DBExec(index string, a ...interface{}) (sql.Result, error) {
	fmt.Printf(unisql[Uni.Config.DBDriver][index], a...)
	fmt.Println()
	return Uni.DB.Exec(fmt.Sprintf(unisql[Uni.Config.DBDriver][index], a...))
}

// Get one index of the first returned column from database (variable remains unaffected if nothing is retrieved)
func (Uni *UniBot) DBGetFirstVar(dest interface{}, index string, a ...interface{}) error {
	rows, err := Uni.DB.Query(fmt.Sprintf(unisql[Uni.Config.DBDriver][index], a...))
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	return rows.Scan(dest)
}

// Same thing as "DBGetFirstVar" except create a new variable
func (Uni *UniBot) DBGetFirst(index string, a ...interface{}) (interface{}, error) {
	rows, err := Uni.DB.Query(fmt.Sprintf(unisql[Uni.Config.DBDriver][index], a...))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	var d interface{}
	err = rows.Scan(d)
	return d, err
}


// UniBot's universal Query function
func (Uni *UniBot) DBQuery(index string, a ...interface{}) (*sql.Rows, error) {
	return Uni.DB.Query(fmt.Sprintf(unisql[Uni.Config.DBDriver][index], a...))	
}