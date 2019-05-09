rm chaindata/*
go run test/full_store.go > chaindata/store.txt
go run test/full_load.go > chaindata/load.txt
echo "done, check store.txt and load.txt in chaindata directory"