package dbaser

import (
	"context"
	"errors"
	"fmt"

	//	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Struct4db struct {
	Ctx        context.Context
	IsBase     bool
	MetricBase *pgx.Conn
}

// func TableGetAllCounters(ctx context.Context, db *pgx.Conn) (map[string]int64, error) {
func TableGetAllCounters(MetricBaseStruct *Struct4db) (map[string]int64, error) {
	var inta int64
	var str string
	mappa := map[string]int64{}
	zapros := "SELECT * FROM counter;"
	rows, err := MetricBaseStruct.MetricBase.Query(MetricBaseStruct.Ctx, zapros)
	if err != nil {
		return nil, fmt.Errorf("error Query %[2]s:%[3]d database  %[1]w", err,
			MetricBaseStruct.MetricBase.Config().Host, MetricBaseStruct.MetricBase.Config().Port)
	}
	for rows.Next() {
		err = rows.Scan(&str, &inta)
		if err != nil {
			return nil, fmt.Errorf("error counter table Scan %[2]s:%[3]d database\n%[1]w", err,
				MetricBaseStruct.MetricBase.Config().Host, MetricBaseStruct.MetricBase.Config().Port)
		}
		mappa[str] = inta
	}
	return mappa, nil
}
func TableGetAllGauges(MetricBaseStruct *Struct4db) (map[string]float64, error) {
	var flo float64
	var str string
	mappa := map[string]float64{}
	zapros := "SELECT * FROM gauge;"
	rows, err := MetricBaseStruct.MetricBase.Query(MetricBaseStruct.Ctx, zapros)
	if err != nil {
		return nil, fmt.Errorf("error Query %[2]s:%[3]d database  %[1]w", err,
			MetricBaseStruct.MetricBase.Config().Host, MetricBaseStruct.MetricBase.Config().Port)
	}
	for rows.Next() {
		err = rows.Scan(&str, &flo)
		if err != nil {
			return nil, fmt.Errorf("error gauge table Scan %[2]s:%[3]d database\n%[1]w", err,
				MetricBaseStruct.MetricBase.Config().Host, MetricBaseStruct.MetricBase.Config().Port)
		}
		mappa[str] = flo
	}
	return mappa, nil
}

func TableCreation(MetricBaseStruct *Struct4db) error {
	crea := "CREATE TABLE IF NOT EXISTS Gauge(metricname VARCHAR(30) PRIMARY KEY, value FLOAT8);"
	tag, err := MetricBaseStruct.MetricBase.Exec(MetricBaseStruct.Ctx, crea)
	if err != nil {
		return fmt.Errorf("error create Gauge table. Tag is \"%s\" error is %w", tag.String(), err)
	}
	crea = "CREATE TABLE IF NOT EXISTS Counter(metricname VARCHAR(30) PRIMARY KEY, value BIGINT);"
	tag, err = MetricBaseStruct.MetricBase.Exec(MetricBaseStruct.Ctx, crea)
	if err != nil {
		return fmt.Errorf("error create Counter table. Tag is \"%s\" error is %w", tag.String(), err)
	}
	return nil
}

func TablePutGauge(MetricBaseStruct *Struct4db, mname string, value float64) error {

	order := fmt.Sprintf("INSERT INTO Gauge(metricname, value) VALUES ('%[1]s',%[2]g);", mname, value)
	tag1, err := MetricBaseStruct.MetricBase.Exec(MetricBaseStruct.Ctx, order)

	//	log.Printf("TableInsertGauge err %v\n db %v\n\n", err, db)
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
		return fmt.Errorf("error Insert Gauge %s with %g value. TagInsert is \"%s\" *pgconn.PgError %+v error is %w",
			mname, value, tag1.String(), pgErr, err)
	}

	order = fmt.Sprintf("UPDATE Gauge SET value=%[2]g WHERE metricname='%[1]s'", mname, value)
	tag2, err := MetricBaseStruct.MetricBase.Exec(MetricBaseStruct.Ctx, order)
	//	log.Printf("TableUpdateGauge err %v\n db %v\n\n", err, db)
	if err == nil {
		return nil
	}
	return fmt.Errorf("error UPDATE Gauge %s with %f value. TagInsert is \"%s\" TagUpdate is \"%s\" error is %w",
		mname, value, tag1.String(), tag2.String(), err)
}

func TableGetMetric() error {

	return nil
}

func TableGetGauge(MetricBaseStruct *Struct4db, mname string, flo *float64) error {
	//var flo float64
	str := "SELECT value FROM gauge WHERE metricname = $1;"
	row := MetricBaseStruct.MetricBase.QueryRow(MetricBaseStruct.Ctx, str, mname)
	err := row.Scan(flo)
	if err != nil {
		return fmt.Errorf("error get %s gauge metric.  %w", mname, err)
	}
	return nil
}
func TableGetCounter(MetricBaseStruct *Struct4db, mname string, inta *int64) error {
	//var inta int64
	str := "SELECT value FROM counter WHERE metricname = $1;"
	row := MetricBaseStruct.MetricBase.QueryRow(MetricBaseStruct.Ctx, str, mname)
	err := row.Scan(inta)
	if err != nil {
		return fmt.Errorf("error get %s counter metric.  %w", mname, err)
	}
	return nil
}

func TablePutCounter(MetricBaseStruct *Struct4db, mname string, value int64) error {
	order := fmt.Sprintf("INSERT INTO Counter(metricname, value) VALUES ('%[1]s',%[2]d);", mname, value)
	tag1, err := MetricBaseStruct.MetricBase.Exec(MetricBaseStruct.Ctx, order)
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
		return fmt.Errorf("error Insert Counter %s with %d value. TagInsert is \"%s\" *pgconn.PgError %+v error is %w",
			mname, value, tag1.String(), pgErr, err)
	}
	order = fmt.Sprintf("UPDATE Counter SET value=value+%[2]d WHERE metricname='%[1]s'", mname, value)
	tag2, err := MetricBaseStruct.MetricBase.Exec(MetricBaseStruct.Ctx, order)
	if err == nil {
		return nil
	}
	return fmt.Errorf("error UPDATE Counter %s with %d value. TagInsert is \"%s\" TagUpdate is \"%s\" error is %w",
		mname, value, tag1.String(), tag2.String(), err)
}

type Gauge float64
type Counter int64
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func TableBuncher(MetricBaseStruct *Struct4db, metrArray []Metrics) error {
	tx, err := MetricBaseStruct.MetricBase.Begin(MetricBaseStruct.Ctx)
	if err != nil {
		return fmt.Errorf("error db.Begin  %[1]w", err)
	}
	var order string
	for _, metrica := range metrArray {
		if metrica.MType == "gauge" {
			order = fmt.Sprintf("UPDATE gauge SET value=%[2]g WHERE metricname='%[1]s'", metrica.ID, *metrica.Value)
		} else {
			order = fmt.Sprintf("UPDATE counter SET value=value+%[2]d WHERE metricname='%[1]s'", metrica.ID, *metrica.Delta)
		}
		tagUpdate, _ := tx.Exec(MetricBaseStruct.Ctx, order)
		tu := tagUpdate.RowsAffected()
		if tu != 0 { // если удалось записать - уже существует и INSERT не нужен
			continue
		}
		if metrica.MType == "gauge" {
			order = fmt.Sprintf("INSERT INTO Gauge(metricname, value) VALUES ('%[1]s',%[2]g);", metrica.ID, *metrica.Value)
		} else {
			order = fmt.Sprintf("INSERT INTO counter(metricname, value) VALUES ('%[1]s',%[2]d);", metrica.ID, *metrica.Delta)
		}
		tagInsert, err := tx.Exec(MetricBaseStruct.Ctx, order)
		if err != nil {
			return fmt.Errorf("TableBuncher error UPDATE Metric %-v TagInsert is \"%s\" TagUpdate is \"%s\" error is %w",
				metrica, tagInsert.String(), tagUpdate.String(), err)
		}
	}
	return tx.Commit(MetricBaseStruct.Ctx)
}
