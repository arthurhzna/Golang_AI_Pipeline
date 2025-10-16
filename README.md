Golang API → Redis Queue → Python AI Container → Redis Queue → Golang Worker → AWS S3 + MQTT

Redis LPUSH/RPOP: ~0.1ms per operation
Redis single-threaded but can handle 100k+ requests/sec
50k-100k+ ops/sec  single instance/thread

Scalability
    scale horizontally:
        Multiple Thread/Worker Python AI ---> x Worker
        Multiple Golang Workers ---> x Worker
        Redis cluster 

todo:
env--> global config 
fix path aws s3 bucket