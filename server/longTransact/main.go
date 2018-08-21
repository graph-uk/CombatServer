package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func clearDB(db *sql.DB) error {
	db.Exec(`
	DELETE FROM Table1;
	DELETE FROM Table2;
	INSERT INTO Table2(sum) VALUES("0");`)
	return nil
}

func writerTransact(db *sql.DB, id string, wg *sync.WaitGroup) {
	for i := 1; i <= 300; i++ {
		tx, err := db.Begin()

		req, err := tx.Prepare("INSERT INTO Table1(someString) VALUES(?)")
		check(err)
		_, err = req.Exec(id + `-` + strconv.Itoa(i))
		if err != nil {
			fmt.Println(id + `-` + strconv.Itoa(i) + `-` + err.Error())
			time.Sleep(10 * time.Millisecond)
			continue
		}

		req, err = tx.Prepare(`SELECT sum FROM Table2`)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		rows, err := req.Query()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var sum int
		rows.Next()
		rows.Scan(&sum)
		rows.Close()

		sum += i

		req, err = tx.Prepare(`UPDATE Table2 SET sum=?`)
		check(err)
		_, err = req.Exec(sum)
		check(err)
		req, err = tx.Prepare("INSERT INTO Table1(someString) VALUES(?)")
		check(err)
		_, err = req.Exec(id + `-` + strconv.Itoa(i))

		tx.Commit()
		time.Sleep(900 * time.Microsecond) // sleep to allow other routines work with DB
		if i%100 == 0 {
			fmt.Println(id + `-` + strconv.Itoa(i))
		}
	}
	fmt.Println(id + `-finished`)
	wg.Done()
}

func writerNoTransact(db *sql.DB, id string, wg *sync.WaitGroup) {
	for i := 1; i <= 300; i++ {
		req, err := db.Prepare("INSERT INTO Table1(someString) VALUES(?)")
		check(err)
		_, err = req.Exec(id + `-` + strconv.Itoa(i))
		if err != nil {
			fmt.Println(id + `-` + strconv.Itoa(i) + `-` + err.Error())
			time.Sleep(10 * time.Millisecond)
			continue
		}

		req, err = db.Prepare(`SELECT sum FROM Table2`)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		rows, err := req.Query()
		if err != nil {
			fmt.Println(id + `-` + strconv.Itoa(i) + `-` + err.Error())
			time.Sleep(10 * time.Millisecond)
			continue
		}

		var sum int
		rows.Next()
		rows.Scan(&sum)
		rows.Close()

		sum += i

		req, err = db.Prepare(`UPDATE Table2 SET sum=?`)
		check(err)
		_, err = req.Exec(sum)
		if err != nil {
			fmt.Println(id + `-` + strconv.Itoa(i) + `-` + err.Error())
			time.Sleep(10 * time.Millisecond)
			continue
		}
		req, err = db.Prepare("INSERT INTO Table1(someString) VALUES(?)")
		check(err)
		_, err = req.Exec(id + `-` + strconv.Itoa(i))

		time.Sleep(900 * time.Microsecond) // sleep to allow other routines work with DB
		if i%100 == 0 {
			fmt.Println(id + `-` + strconv.Itoa(i))
		}
	}
	fmt.Println(id + `-finished`)
	wg.Done()
}

func checkResult(db *sql.DB) error {

	req, err := db.Prepare(`SELECT sum FROM Table2`)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	rows, err := req.Query()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var sum int
	rows.Next()
	rows.Scan(&sum)
	rows.Close()

	fmt.Println(`sum=` + strconv.Itoa(sum))
	return nil
}

func main() {

	db, err := sql.Open("sqlite3", "db.sl3")
	check(err)
	defer db.Close()

	clearDB(db)

	workersCount := 4
	var wg sync.WaitGroup
	wg.Add(workersCount)

	for i := 1; i <= workersCount; i++ {
		go writerNoTransact(db, strconv.Itoa(i), &wg)
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	checkResult(db)

	clearDB(db)
	wg.Add(workersCount)
	for i := 1; i <= workersCount; i++ {
		go writerTransact(db, strconv.Itoa(i), &wg)
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	checkResult(db)

	fmt.Println(`All routines are finished. Press enter to exit.`)
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
