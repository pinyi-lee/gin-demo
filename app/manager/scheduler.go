package manager

var (
	schedulerManager *SchedulerManager
)

func GetScheduler() *SchedulerManager {
	return schedulerManager
}

type SchedulerManager struct {
}

func (manager *SchedulerManager) Setup() (err error) {

	return nil
}

func (manager *SchedulerManager) Close() {

}

func (manager *SchedulerManager) Run() (err error) {

	return nil
}
