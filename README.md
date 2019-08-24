# comp5617_self_assessment
code for unit COMP5617 self_assessment

# Requirement
* go1.12.9

# Running 
## A.HTTP scanner
To get the output of the part A, run following commands:

`go run scanner_a.go top100.csv`

`go run scanner_a.go rand100.csv`

The output files would be generated as `top-100_output_a.csv` or `rand-100_output_a.csv` under the current directory

## B.Finetuning
To get the output of the part B, run following commands:

`go run scanner_b.go top100.csv`

`go run scanner_b.go rand100.csv`

The output files would be generated as `top-100_output_c.csv` or `rand-100_output_c.csv` under the current directory

## C.Speeding it up
To get the output of the part C, run following commands:

`go run scanner_c.go top100.csv`

`go run scanner_c.go rand100.csv`

The output files would be generated as `top-100_output_c.csv` or `rand-100_output_c.csv` under the current directory

## D.Working with the output
To get the output of the part C, run following commands:

`go run count_afterwhile_redirect.go top100.csv`

`go run count_afterwhile_redirect.go rand100.csv`

`go run count_ciphers.go top100.csv`

`go run count_ciphers.go rand100.csv`

The output files would be generated as `top-100_count_afterwhile_redirect.csv`, `rand-100_count_afterwhile_redirect.csv` ,`top-100_count_cipher.csv`,`rand-100_count_cipher.csv`under the current directory. These two files are used to count in related parts in jupyter notebook. 
