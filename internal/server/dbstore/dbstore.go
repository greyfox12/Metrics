package dbstore

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/greyfox12/Metrics/internal/server/logmy"
	"github.com/greyfox12/Metrics/internal/server/storage"
)

// Создаю объекты БД
func CreateDB(db *sql.DB) error {
	var Script string
	var errdb error

	pwd, _ := os.Getwd()
	fmt.Printf("Currect pass=%v\n", pwd)
	ent, _ := os.ReadDir("./")
	for _, e := range ent {
		fmt.Printf("Currect dirlist=%v\n", e.Name())
	}

	file, err := os.Open("./internal/server/dbstore/Script.sql")
	//	file, err := os.Open("../../internal/server/dbstore/Script.sql")
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
	_, errdb = db.Exec(Script)

	if errdb != nil {
		fmt.Printf("%v\n", errdb)
		logmy.OutLog(errdb)
		return error(errdb)
	}

	//	fmt.Printf("%v\n", Script)
	return nil
}

// Прочитать данные из DB
func GetDBGauge(db *sql.DB, id string) (float64, error) {
	var mgauge float64

	rows, err := db.Query("SELECT get_gauge($1) gauge", id)
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

func GetDBCounter(db *sql.DB, id string) (int64, error) {
	var mcounter int64

	rows, err := db.Query("SELECT get_counter($1) counter", id)
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
func SetDBGauge(db *sql.DB, id string, par float64) (float64, error) {
	var mgauge float64
	//	fmt.Printf("id=%v, par=%v\n", id, par)
	rows, err := db.Query("SELECT set_gauge($1, $2) gauge", id, par)
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

	//	db.Commit()
	return mgauge, nil
}

func SetDBCounter(db *sql.DB, id string, par int64) (int64, error) {
	var mcounter int64

	rows, err := db.Query("SELECT set_counter($1, $2) counter", id, par)
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

	for _, val := range mgauge.Keylist() {
		v, err := mgauge.Get(val)
		if err != nil {
			continue
		}
		if _, err := SetDBGauge(db, val, v); err != nil {
			logmy.OutLog(err)
			continue
		}
	}

	for _, val := range mmetric.Keylist() {
		v, err := mmetric.Get(val)
		if err != nil {
			continue
		}
		if _, err := SetDBCounter(db, val, v); err != nil {
			logmy.OutLog(err)
			continue
		}
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
