### How to run

Run `go build -o server cmd/server/main.go` to build the binary.

Run `./server` to start the server.

By default, it runs on port 8080. Pass `-port` flag to change the port.

#### Load data 

Use `/v1/admin/process-csv` endpoint to process and store csv file data. Pass the csv file path in the request.

Use `/v1/admin/switch-db` endpoint to mark newly process csv data as primary location for read.

Use `/v1/promotions/{id}` to query the data.

[Swagger Doc](api/v1/rest/openapi.yaml)

### Implementation Details

Current implementation is using Primary/Secondary (or Blue/Green) storage 
strategy. This works by using two tables in a database. One is marked primary 
and will serve the read requests. The other one is marked secondary is used to
load the data. When the data the processed when can manually switch the primary
with secondary table and new data is started being served. After that we can 
purge the old data.

This way we can keep read and write part separate and both processes can run 
simultaneously without impacting the performance of each other.

The csv file processing can be moved to separate CLI instead of http API.
This might be needed if csv is very large and take some time to process.



### Questions
- The .csv file could be very big (billions of entries) - how would your application
perform?

I have tried to handle it with by processing the file concurrently and process chunks of 
file in parallel using all available cores on our server.

For more performance this can also done using distributed processing. Divide the file 
in multiple parts and process each on separate machine.

We can also utilze MapReduce workflow using Hadoop or Spark to do the distributed processing of the large file. 


- Every new file is immutable, that is, you should erase and write the whole storage;

I have tried to handle this by using Primary/Secondary strategy. This works by using 
two tables in a database. One is marked primary and will serve the read requests. The other 
one is marked secondary is used to load the data.


- How would your application perform in peak periods (millions of requests per
minute)?

As the read part of application is separate from the write part. The application perform really
well in peak periods. As it is simple read API with simple query. Performance can be increased
by incorporating caching in read side of the api. 


- How would you operate this app in production (e.g. deployment, scaling, monitoring)?

In production the app can have the following separate components

    - Read API
    - Read Cache
    - Storage (database)
    - Data Processing Pipeline

`Read API` can scale independently based on the read request volumes. 

`Read Cache` can also be dynamically scaled during peak periods.

`Storage` can run separately and will have relatively fixed requirement from
disk space and compute point of view.

`Data Processing Pipeline` can be adjust based on the amount of data that
needed to be processed.

We will need to monitor RPS and latency of requests for read API.
We will also need to monitor CPU usage for read API infrastructure.

CPU and Disk need to be monitored for Storage and Cache infrastructure.

