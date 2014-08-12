go test -v -run=TestSerialization
redis-cli FLUSHALL
go test -v -run=TestMysqlWithCache
