Golang API → Redis Queue → Python AI Container → Redis Queue → Golang Worker → AWS S3 + MQTT

Redis LPUSH/RPOP: ~0.1ms per operation
Redis single-threaded tapi bisa handle 100k+ requests/sec
50k-100k+ ops/sec pada single instance

Scalability
    Bisa scale horizontally:
        Multiple Python AI containers  ---> x Worker
        Multiple Golang workers ---> x Worker
        Redis cluster jika needed