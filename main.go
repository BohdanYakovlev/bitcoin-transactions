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

const limitSizeBlockBytes = 1000

var limitTimeMillisecond time.Duration

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

func (data *block) validateTransaction(transactionData transaction) error {
	if transactionData.size > limitSizeBlockBytes {
		return errors.New(fmt.Sprintf("id: %s transaction too large", transactionData.id))
	}

	if transactionData.fee <= 0 {
		return errors.New(fmt.Sprintf("id: %s transaction fee must be greater than zero", transactionData.id))
	}

	if transactionData.size <= 0 {
		return errors.New(fmt.Sprintf("id: %s transaction size must be greater than zero", transactionData.id))
	}

	return nil
}

func (data *block) recordToTransaction(record []string) (transaction, error) {

	var result transaction

	result.id = record[0]
	result.size, _ = strconv.Atoi(record[1])
	result.fee, _ = strconv.Atoi(record[2])

	return result, data.validateTransaction(result)
}

func (data *block) handleRecord(record []string) {

	transactionData, err := data.recordToTransaction(record)
	if err != nil {
		//fmt.Println(err)
		return
	}

	if data.totalSize+transactionData.size <= limitSizeBlockBytes {
		data.addTransaction(transactionData, 0)
		return
	}
	data.weighTransaction(transactionData)
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
			index++
			break
		}
	}
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

func (data *block) countTotalSize() int64 {
	var totalSize int64 = 0
	for _, transactionData := range data.transactions {
		totalSize += int64(transactionData.size)
	}
	return totalSize
}

func (data *block) printTransactions() {
	fmt.Println("Selected transactions:")
	for _, iterator := range data.transactions {
		fmt.Printf("Id: %s, Size: %d, Fee: %d\n", iterator.id, iterator.size, iterator.fee)
	}
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
		if endTime.Sub(startTime) >= limitTimeMillisecond {
			isBreakFlag = false
		}
	}
	return result
}

func getDataFilePath() string {
	var path string
	var tempTime int

	fmt.Println("Enter data path:")

	_, err := fmt.Scan(&path)
	if err != nil {
		fmt.Println("Can not read data")
	}

	fmt.Println("Enter preferred program running time in milliseconds:")
	_, err = fmt.Scan(&tempTime)
	if err != nil {
		fmt.Println("Can not read time")
	}

	limitTimeMillisecond = time.Duration(tempTime) * time.Millisecond

	return path
}

func printResult(result block, timeOfWork time.Duration) {

	result.printTransactions()

	fmt.Printf("Amount of transactions: %d\n", len(result.transactions))

	fmt.Printf("Block size: %d bytes\n", result.totalSize)

	fmt.Printf("The total extracted value: %d\n", result.totalFee)

	fmt.Printf("Construction time: %s\n", timeOfWork.String())
}

func main() {

	//dataFilePath := "transactions.csv"
	//limitTimeMillisecond = 30 * time.Millisecond

	dataFilePath := getDataFilePath()

	startTime := time.Now()
	result := readTransactionsFromCSV(dataFilePath, startTime)
	endTime := time.Now()

	printResult(result, endTime.Sub(startTime))

}
