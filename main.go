package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	limitTimeMillisecond = 1000 * time.Millisecond
	limitSizeBlockBytes  = 10
)

type transaction struct {
	id   string
	size int
	fee  int
}

func (data *transaction) getFeePerSize() float32 {
	return float32(data.fee) / float32(data.size)
}

type block struct {
	transactions []transaction
	totalSize    int
	totalFee     int
}

func (data *block) init() {
	data.totalSize = 0
	data.totalFee = 0
	data.transactions = make([]transaction, 0, 125)
}

func (data *block) handleRecord(record []string) {

	var tempTransaction transaction
	tempTransaction.id = record[0]
	tempTransaction.size, _ = strconv.Atoi(record[1])
	tempTransaction.fee, _ = strconv.Atoi(record[2])

	if tempTransaction.size > limitSizeBlockBytes {
		return
	}
	newSize := data.totalSize + tempTransaction.size
	if newSize <= limitSizeBlockBytes {
		data.addTransaction(tempTransaction, 0)
		return
	}
	data.weighTransaction(tempTransaction)
}

func (data *block) weighTransaction(transactionData transaction) {

	if data.transactions[0].getFeePerSize() > transactionData.getFeePerSize() {
		return
	}

	freeSpace := limitSizeBlockBytes - data.totalSize
	feeLost := 0
	sizeLost := 0
	index := 0

	for ; index < len(data.transactions); index++ {
		sizeLost += data.transactions[index].size
		feeLost += data.transactions[index].fee

		if feeLost >= transactionData.fee {
			return
		}
		if sizeLost+freeSpace >= transactionData.size {
			break
		}
	}
	index++

	data.totalSize -= sizeLost
	data.totalFee -= feeLost

	data.addTransaction(transactionData, index)
}

func (data *block) addTransaction(transactionData transaction, index int) {
	data.totalSize += transactionData.size
	data.totalFee += transactionData.fee
	data.transactions = append(data.transactions, transactionData)

	tempMas := data.transactions[index:]
	sort.Slice(tempMas, func(i, j int) bool {
		return tempMas[i].getFeePerSize() < tempMas[j].getFeePerSize()
	})

	data.transactions = tempMas
}

func getReader(filePath string) (*csv.Reader, *os.File) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return csv.NewReader(file), file
}

func getRecord(reader *csv.Reader) ([]string, bool) {
	record, err := reader.Read()
	if err == io.EOF {
		return record, false
	}
	if err != nil {
		panic(err)
	}
	if len(record) != 3 {
		panic(errors.New("incorrect record"))
	}
	return record, true
}

func readTransactionsFromCSV(path string, startTime time.Time) block {

	reader, file := getReader(path)
	defer file.Close()

	record, _ := getRecord(reader)

	var result block
	result.init()

	isBreakFlag := true
	for isBreakFlag {

		record, isBreakFlag = getRecord(reader)
		if isBreakFlag {
			result.handleRecord(record)
		}

		endTime := time.Now()
		if endTime.Sub(startTime) > limitTimeMillisecond {
			isBreakFlag = false
		}
	}
	return result
}

func main() {
	fmt.Println(readTransactionsFromCSV("transactions.csv", time.Now()))
}
