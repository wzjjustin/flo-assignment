# Flo Energy Tech Assessment by Justin Wong (wzjjustin@gmail.com)
---
*Disclaimer: Assessment was done on Windows 10*
## 1. Pre-requisites
  - Docker (v28.0.4 Installed) 
  - Golang (v1.25.0+)
## 2. How to run
### Start up postgresql database  
*ensure values are aligned to ~/config/config.yaml
Run in terminal:
> make start-db

### Build executable
Change directory to project folder and run build command in terminal:
> cd ~/flo_assignment

> go build

The flo-assignment.exe will be generated in the root folder

### Running the file parser
This executable has 2 functionalities:
1. Parse file example:
  > ./flo-assignment.exe parse --filepath="test_data/happy_flow_5.csv"
  

2. Clean up database
-  Reason for clean up functionality is due to checksum implemented. To avoid same file from being processed twice, which will result in duplicated data with different UUIDs.
-  This operation will drop existing tables which will be recreated during setup via auto migration in the next parse operation

  > ./flo-assginment.exe clean

### Testing
1. Unit test:
  > make run-unit-test

2. Functional test:
  > make run-happy-flow
  > make run-unhappy-flow

### Clean up docker
Run in terminal:
> make clean-db


## 3. Questions
### Q1. What is the rationale for the technologies you have decided to use?

Answer:
> Golang - familiarity with the language and its in-built concurrency model is easy to use. Goroutines are lightweight (less memory) that comes efficient and safe communication via channels. Little build time and good debugging tools.

> Docker - Lightweight and portable setup for postgresql on local systems.

> github.com/urfave/cli - library to quickly build command line tools, allowing for different functionalities to be performed in a single binary (in this case - parse/clean).

> gorm.io (or originally github.com/jinzhu/gorm) - Full featured golang orm that provides auto migration of models into tables, db transaction capabilities, easy CRUD functionality and supports various sql drivers (even custom drivers).

### Q2. What would you have done differently if you had more time?
Answer:
> Spent more time reading and understanding the NEM12/13 documentation to possibly find more edge cases that I may have missed out during this assignment. There are many details in the documentation that were helpful in the implementation of this assignment. Some examples were the blocking cycle logic, the different fields to the values in each record that led me to validate interval data. Adding logic for the parsing of 400/500 records.

### Q3. What is the rationale for the design choices that you have made?
Answer:
> Supporting config files - to allow for the user to customise their own set up.

> Auto-DB migration - reduce the need for users to perform an additional set just for db set up prior to running the program.

> Using a dynamic worker pool - this splits up the load into multiple goroutines to run segments of operation concurrently while also allowing for efficient resource usage and scalability. 

> Using tx.CreateInBatches() - this function generates a bulk insert [ example: INSERT INTO table(colA,colB) VALUES (a,b),(c,d)(e,f)...(y,z); ] which helps to reduce overhead when creating multiple records in DB.

# Flow diagram
<img width="1300" height="788" alt="Image" src="https://github.com/user-attachments/assets/02c89ab0-448e-4723-98f1-22ea41d69e4c" />