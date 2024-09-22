package settings

import (
	"sort"
	"sync"
)

// Сервис для хранения текущих периодов расчёта
// Данные хранятся в отсортированном виде, т.к. необходимо:
// - получать максимальный элемент (для хранения нужного количества данных в памяти)
// - добавлять новый элемент при подключении нового клиента
// - удалять существующий элемент при отключении клиента.
type Service struct {
	calcPeriods []uint32
	mx          sync.RWMutex
}

func NewService() *Service {
	return &Service{
		calcPeriods: make([]uint32, 0),
	}
}

func (s *Service) Add(calcPeriod uint32) {
	s.mx.Lock()
	defer s.mx.Unlock()

	idx := s.search(calcPeriod)
	s.calcPeriods = append(s.calcPeriods[:idx], append([]uint32{calcPeriod}, s.calcPeriods[idx:]...)...)
}

func (s *Service) Remove(calcPeriod uint32) bool {
	s.mx.Lock()
	defer s.mx.Unlock()

	idx := s.search(calcPeriod)
	if idx < len(s.calcPeriods) && s.calcPeriods[idx] == calcPeriod {
		s.calcPeriods = append(s.calcPeriods[:idx], s.calcPeriods[idx+1:]...)
		return true
	}
	return false
}

func (s *Service) GetMax() (uint32, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if len(s.calcPeriods) == 0 {
		return 0, false
	}
	return s.calcPeriods[len(s.calcPeriods)-1], true
}

func (s *Service) search(calcPeriod uint32) int {
	return sort.Search(len(s.calcPeriods), func(i int) bool { return s.calcPeriods[i] >= calcPeriod })
}
