## Steps to run the project

1. Bring up the database using the following command
```
docker-compose up -d
```

2. Bring up rabbitmq using the following command
```
docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
```

3. Build the project using the following command
```
go build
```

4. Run the project using the following command
```
./producer_consumer
```

5. Send a POST request to the following endpoint to create a new task
```
http://localhost:8000/api/v1/product
```

6. To run the consumer, run the following command
```
cd consumer
go build
./consumer
```

