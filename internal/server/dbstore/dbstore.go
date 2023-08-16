package dbstore

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/greyfox12/Metrics/internal/server/logmy"
	"github.com/greyfox12/Metrics/internal/server/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/lib/pq"
)

// Создаю объекты БД
func CreateDB(db *sql.DB) error {
	var Script string
	//	var errdb error
	var path string

	pwd, _ := os.Getwd()
	fmt.Printf("Currect pass=%v\n", pwd)
	//	ent, _ := os.ReadDir("./")
	//	for _, e := range ent {
	//		fmt.Printf("Currect dirlist=%v\n", e.Name())
	//	}

	if strings.HasPrefix(pwd, "C:\\GoYandex") {
		path = "../../internal/server/dbstore/Script.sql"
	} else {
		path = "./internal/server/dbstore/Script.sql"
	}

	file, err := os.Open(path)
	if err != nil {
		logmy.OutLog(err)
		return error(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Script = Script + scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		logmy.OutLog(err)
		return error(err)
	}
	//	fmt.Printf("%v\n", Script)
	Errdb := ResendDB(db, Script)

	if Errdb != nil {
		logmy.OutLog(Errdb)
		return error(Errdb)
	}

	//	fmt.Printf("%v\n", Script)
	return nil
}

// Считаю задержку
func WaitSec(s int) int {
	switch s {
	case 2:
		return 1
	case 3:
		return 3
	case 4:
		return 5
	default:
		return 0
	}
}

func ResendDB(db *sql.DB, Script string) error {
	var Errdb error
	var pgErr *pgconn.PgError

	for i := 1; i <= 4; i++ {
		if i > 1 {
			fmt.Printf("Pause: %v sec\n", WaitSec(i))
			time.Sleep(time.Duration(WaitSec(i)) * time.Second)
		}

		_, Errdb = db.Exec(Script)
		if Errdb == nil {
			return nil
		}

		// Проверяю тип ошибки
		fmt.Printf("Error DB: %v\n", Errdb)

		if errors.As(Errdb, &pgErr) {
			//			fmt.Println(pgErr.Code)    // => syntax error at end of input
			//			fmt.Println(pgErr.Message) // => 42601
			if !pgerrcode.IsConnectionException(pgErr.Code) {
				return Errdb // Ошибка не коннекта
			}
		}
	}
	return Errdb
}

// Прочитать данные из DB
// Повторяю Чтение
func QueryDBRet(ctx context.Context, db *sql.DB, sql string, id string) (*sql.Rows, error) {
	var err error
	var pgErr *pgconn.PgError

	for i := 1; i <= 4; i++ {
		if i > 1 {
			fmt.Printf("Pause: %v sec\n", WaitSec(i))
			time.Sleep(time.Duration(WaitSec(i)) * time.Second)
		}

		rows, err := db.QueryContext(ctx, sql, id)
		if err == nil {
			return rows, nil
		}

		// Проверяю тип ошибки
		fmt.Printf("Error DB: %v\n", err)

		if errors.As(err, &pgErr) {
			//			fmt.Println(pgErr.Code)    // => syntax error at end of input
			//			fmt.Println(pgErr.Message) // => 42601
			if !pgerrcode.IsConnectionException(pgErr.Code) {
				return nil, err // Ошибка не коннекта
			}
		}
	}
	return nil, err
}

func GetDBGauge(ctx context.Context, db *sql.DB, id string) (float64, error) {
	var mgauge float64

	rows, err := QueryDBRet(ctx, db, "SELECT get_gauge($1) gauge", id)
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}
	defer rows.Close()

	err = rows.Scan(&mgauge)
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}

	err = rows.Err()
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}
	return mgauge, nil
}

func GetDBCounter(ctx context.Context, db *sql.DB, id string) (int64, error) {
	var mcounter int64

	rows, err := QueryDBRet(ctx, db, "SELECT get_counter($1) counter", id)
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}
	defer rows.Close()

	err = rows.Scan(&mcounter)
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}

	err = rows.Err()
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}
	return mcounter, nil
}

// Записать данные в DB
// Повторяю Чтение
func SetQueryDBRet(ctx context.Context, tx *sql.Tx, sql string, id string, par string) (*sql.Rows, error) {
	var err error
	var pgErr *pgconn.PgError

	for i := 1; i <= 4; i++ {
		if i > 1 {
			fmt.Printf("Pause: %v sec\n", WaitSec(i))
			time.Sleep(time.Duration(WaitSec(i)) * time.Second)
		}

		rows, err := tx.QueryContext(ctx, sql, id, par)
		if err == nil {
			return rows, nil
		}

		// Проверяю тип ошибки
		fmt.Printf("Error DB: %v\n", err)

		if errors.As(err, &pgErr) {
			//			fmt.Println(pgErr.Code)    // => syntax error at end of input
			//			fmt.Println(pgErr.Message) // => 42601
			if !pgerrcode.IsConnectionException(pgErr.Code) {
				return nil, err // Ошибка не коннекта
			}
		}
	}
	return nil, err
}

func SetDBGauge(ctx context.Context, tx *sql.Tx, id string, par float64) (float64, error) {
	var mgauge float64
	//	fmt.Printf("id=%v, par=%v\n", id, par)

	rows, err := SetQueryDBRet(ctx, tx, "SELECT set_gauge($1, $2) gauge", id, fmt.Sprint(par))
	//	tx.QueryContext(ctx, "SELECT set_gauge($1, $2) gauge", id, par)

	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&mgauge)
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}

	err = rows.Err()
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}

	return mgauge, nil
}

func SetDBCounter(ctx context.Context, tx *sql.Tx, id string, par int64) (int64, error) {
	var mcounter int64

	rows, err := SetQueryDBRet(ctx, tx, "SELECT set_counter($1, $2) counter", id, fmt.Sprint(par))
	//	tx.QueryContext(ctx, "SELECT set_counter($1, $2) counter", id, par)
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&mcounter)
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}

	err = rows.Err()
	if err != nil {
		logmy.OutLog(err)
		return 0, err
	}

	//	db.Commit()
	return mcounter, nil
}

// Записать метрики в DB
func SaveMetric(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter, db *sql.DB) error {

	fmt.Printf("Save DB\n")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logmy.OutLog(err)
		return err
	}
	// Откладываем откат на случай, если что-то не удастся.
	defer tx.Rollback()

	for _, val := range mgauge.Keylist() {
		v, err := mgauge.Get(val)
		if err != nil {
			continue
		}
		if _, err := SetDBGauge(ctx, tx, val, v); err != nil {
			logmy.OutLog(err)
			continue
		}
	}

	for _, val := range mmetric.Keylist() {
		v, err := mmetric.Get(val)
		if err != nil {
			continue
		}
		if _, err := SetDBCounter(ctx, tx, val, v); err != nil {
			logmy.OutLog(err)
			continue
		}
	}

	//Фиксирую транзакцию
	if err = tx.Commit(); err != nil {
		logmy.OutLog(err)
		return err
	}

	return nil
}

// // Загрузка данных из DB
func LoadMetric(mgauge *storage.GaugeCounter, mmetric *storage.MetricCounter, db *sql.DB) error {
	var g sql.NullFloat64
	var c sql.NullInt64
	var mtype string
	var id string
	fmt.Printf("Load data from DB \n")

	rows, err := db.Query("SELECT id, mtype, gauge, counter FROM metrics")
	if err != nil {
		logmy.OutLog(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &mtype, &g, &c)
		if err != nil {
			logmy.OutLog(err)
			return err
		}
		switch mtype {
		case "gauge":
			mgauge.Set(id, g.Float64)
		case "counter":
			mmetric.Set(id, c.Int64)
		default:
			logmy.OutLog(errors.New("LoadMetric: Неизвестный тип метрики"))
		}
	}

	err = rows.Err()
	if err != nil {
		logmy.OutLog(err)
		return err
	}
	return nil
}
