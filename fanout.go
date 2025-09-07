package thebus

import "time"

func (b *bus) runFanOut(topic string, state *topicState) {
	defer state.wg.Done()

	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}

	for mr := range state.inQueue {
		// snapshot sous RLock
		b.mutex.RLock()
		subs := snapshotSubsLocked(state.subs)
		b.mutex.RUnlock()

		for _, sub := range subs {
			msg := makeMessage(topic, mr, sub)
			if tryDeliver(sub, msg, timer) {
				state.counters.Delivered.Add(1)
				b.totals.Delivered.Add(1)
			} else {
				state.counters.Dropped.Add(1)
				b.totals.Dropped.Add(1)
			}
		}
	}
}

// snapshotSubsLocked caller must hold RLock
func snapshotSubsLocked(m map[string]*subscription) []*subscription {
	subs := make([]*subscription, 0, len(m))
	for _, s := range m {
		subs = append(subs, s)
	}
	return subs
}

func tryDeliver(sub *subscription, msg Message, timer *time.Timer) bool {
	cfg := sub.cfg
	if cfg.DropIfFull {
		select {
		case sub.messageChan <- msg:
			return true
		default:
			return false
		}
	}
	// timeout
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(cfg.SendTimeout)

	select {
	case sub.messageChan <- msg:
		// flosh timer if needed
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		return true
	case <-timer.C:
		return false
	}
}
