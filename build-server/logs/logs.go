package logs

import (
	"github.com/chrollo-lucifer-12/shared/db"
	"github.com/google/uuid"
)

type LogDispatcher struct {
	seq     int64
	logChan chan db.LogEvent
}

func NewLogDispatcher(buffer int) *LogDispatcher {
	return &LogDispatcher{
		seq:     0,
		logChan: make(chan db.LogEvent, buffer),
	}
}

func (l *LogDispatcher) Push(deploymentId uuid.UUID, line string) {
	l.seq++

	l.logChan <- db.LogEvent{
		DeploymentID: deploymentId,
		Log:          line,
		Sequence:     l.seq,
	}
}

func (l *LogDispatcher) Channel() <-chan db.LogEvent {
	return l.logChan
}

func (l *LogDispatcher) Close() {
	close(l.logChan)
}
