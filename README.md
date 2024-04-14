# Bitcoin Transactions Console App
[Task #5](https://docs.google.com/document/d/1NgViDj9ELMuFMQDc4aOFCs2jMUkKfd9lHGViSwtFdV0/edit?usp=sharing) in Distributed Lab Challenge
This console application is designed to construct a Bitcoin block from a list of transactions, maximizing the extracted fee while staying within the block size limit. The program takes a CSV file containing transaction data as input and provides various output metrics related to the constructed block.
## How to Run the Code

To run this console app, follow these steps:

1. Clone this repository to your local machine:

```
git clone <repository_url>
```

2. Navigate to the directory containing the Go code:

```
cd bitcoin-transactions
```

3. Ensure you have Go installed on your system. If not, you can download it from [here](https://golang.org/dl/).

4. Execute the Go code by running the following command:

```
go run main.go
```

5. Follow the on-screen instructions to provide the paths to your CSV files containing website analytics data for the first and second days.

6. Once the program finishes executing, it will display the users who visited some pages on both days and users who visited a page on the second day that they hadn't visited on the first day.

## Input CSV File Format
The input CSV file should have the following structure:
```
tx_id,tx_size,tx_fee
```
Where:
  - tx_id is the transaction ID.
  - tx_size is the size of the transaction in bytes.
  - tx_fee is the fee of the transaction in satoshis.

## Output
After processing the transactions, the application will display the following information:

  - Constructed block: List of transactions included in the block.
  - Amount of transactions in the block.
  - The block size in bytes.
  - The total extracted value (total fees collected) in satoshis.
  - Construction time: Time taken to construct the block.
## Algorithm:
  - Read CSV file.
  - Add the transaction to the block if there is enough empty space or if the fee of the element exceeds the total fee of loss of valuable transaction occupying the same memory. Transaction value is the price to size ratio
  - Sort block by transaction value.
  - Return result if the time is up or the file runs out
## Efficiency:
The application utilizes a priority queue (min heap) data structure to efficiently select transactions with the highest fee-to-size ratio first. This approach ensures that transactions with the highest fee relative to their size are included in the block first, maximizing the extracted fee while staying within the block size limit.
## Conclusion
This console application efficiently constructs a Bitcoin block by selecting transactions with the highest fee-to-size ratio, ensuring maximum fee extraction within the 1MB block size limit. It offers a balance between time and space efficiency, making it suitable for processing large transaction datasets. Feel free to reach out if you have any questions or encounter any issues while running the code.
