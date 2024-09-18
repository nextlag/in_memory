package in_memory

type InMemoryOption func(*Engine)

func WithPartitions(partitionsNumber uint) InMemoryOption {
	return func(db *Engine) {
		db.partitions = make([]*HashTable, partitionsNumber)
		for i := 0; i < int(partitionsNumber); i++ {
			db.partitions[i] = NewHashTable()
		}
	}
}
