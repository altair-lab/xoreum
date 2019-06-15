#### go-ethereum의 LevelDB 스키마를 대부분 따름
#### 일부 스키마 및 인코딩 방식 등을 커스터마이징


1. database.go
    - DB 생성 인터페이스 함수 구현
    - in-memory (NewMemoryDatabase) 및 persistent (=disk) (NewLevelDBDatabase) DB 구현

2. io_chain.go
    - 블록 데이터 I/O 함수 (read, write, delete) 구현
    - Hash, Header, Body, Td (=total difficulty), Block 등의 I/O 구현
    - GenesisHeaderHash, LastHeaderHash 접근 가능

3. io_state.go
    - State 데이터 I/O 함수 (read, write, delete) 구현
    - 추가적으로 테스팅을 위한 ReadStates, CountStates 구현

4. io_txs.go
    - Tx (=transaction) 데이터 I/O 함수 (read, write, delete) 구현
    - 블록 안에 저장된 tx와 밖에 따로 저장된 tx가 있으므로 두 케이스 모두에 대해 구현
    
    (1) 블록 안에 저장된 tx: TxLookupEntry, tx가 저장된 블록의 해시 return
    
    (2) 블록 밖에 저장된 tx: RawTxData, tx 자체를 저장

5. schema.go
    - var (...)에 로우레벨 DB prefixing 스키마 설명 (실제 LevelDB에 저장되는 key 값)
    
      e.g. lastHeaderKey = []byte("LastHeader")
    - 각각의 key 값에 대한 인터페이스 함수 
