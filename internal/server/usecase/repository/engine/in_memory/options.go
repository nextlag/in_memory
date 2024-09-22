package in_memory

type EngineOption func(*Engine)

func WithPartitions(partitionsNum int) EngineOption {
	return func(e *Engine) {
		e.partitions = make([]*HashTable, partitionsNum)
		for i := 0; i < partitionsNum; i++ {
			e.partitions[i] = NewHashTable()
		}
	}
}
