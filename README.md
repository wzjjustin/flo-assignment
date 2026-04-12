# Flo Energy Tech Assessment by Justin Wong (wzjjustin@gmail.com)
---
*Disclaimer: Assessment was done on Windows 10*

# Flow diagram
<img width="1300" height="788" alt="Image" src="https://github.com/user-attachments/assets/02c89ab0-448e-4723-98f1-22ea41d69e4c" />

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
> Spent more time reading and understanding the NEM12/13 documentation to possibly find more edge cases that I may have missed out during this assignment. There are many details in the documentation that were helpful in the implementation of this assignment. Some examples were the blocking cycle logic, the different fields to the values in each record that led me to validate interval data. 

> Implementing more logic to handle specific scenarios, for example (refer to last point in readme - my ideal solution):
>>1. Fail safe - The current implementation terminates all operation when a single occurence of error appears in the parser's flow. To mitigate the inefficiency of having to read all data from the start again in a separate instance, instead, an implementation that accepts data blocks (in chunks of type 200-500) which pass all the validation requirements (store them into db) and buffers the failures into a backup/retry file to attempt retry on a scheduled cronjob or via manual intiation after data rectification.
>>2. Separate extraction from preprocessing stage -  the largest bottle neck of my implementation is the extraction stage as it handles both preprocessing and extraction, slowing down the process. Even with a pool of dynamic workers to pick up the transformation of data for loading, the preprocessing stage is currently read and processed via a single line scanner. In my ideal implementation, I would separate out the preprocessing stage and extraction stage into their individual operations. Preprocessing stage resolve issues such as file format incorrectness, allowing for the extraction stage to flow faster. While also enhancing the extractor pool to process in chunks of data in segments of 200-500 records. (interesting read: https://medium.com/swlh/processing-16gb-file-in-seconds-go-lang-3982c235dfa2)
>>3. With multiple transformers to convert data into data we want to see in the DB, having more workers to handle the db storing/fail safe file creation will also help to speed up the entire parser's process. 

### Q3. What is the rationale for the design choices that you have made?
Answer:
> Supporting config files - to allow for the user to customise their own set up.

> Auto-DB migration - reduce the need for users to perform an additional set just for db set up prior to running the program.

> Using a dynamic worker pool - this splits up the load into multiple goroutines to run segments of operation concurrently while also allowing for efficient resource usage and scalability. 

> Using tx.CreateInBatches() - this function generates a bulk insert [ example: INSERT INTO table(colA,colB) VALUES (a,b),(c,d)(e,f)...(y,z); ] which helps to reduce overhead when creating multiple records in DB.

# My Ideal Solution
<img width="1209" height="679" alt="Image" src="https://github.com/user-attachments/assets/4f86a531-e4c3-444c-8ffb-20fb92882525" />