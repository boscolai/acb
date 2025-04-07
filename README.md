# ACB - Average Cost Base
acb is a command line tool to compute average cost base of security sells.  The current implementation can parse downloaded csv files from National Bank Direct Brokerage (NBDB).

For each given csv file containing security transactions in reversed chronicle order,
this tool computes ACB values for each sell and outputs 3 corresponding csv files:

- {infile}_gains.csv: a list of SELL transactions with Gains (or loss)
- {infile}_gains_summary: a list of Gains (Loss) by ticker symbols
- {infile}_holdings.csv: a list of current holdings

```
Usage example:
acb -d \; account1.csv account2.csv

Usage:
acb csv_file1 [csv_file2 [csv_file2 ...]] [flags]

Flags:
-d, --delimiter string   field delimiter in csv file (default ",")
-h, --help               help for acb

```


## Build

To build the executable, use the make target below.  The executable will be created as bin/acb.

```
make build
```









