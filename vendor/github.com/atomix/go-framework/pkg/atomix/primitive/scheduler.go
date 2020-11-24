// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package primitive

import (
	"container/list"
	"time"
)

// Scheduler provides deterministic scheduling for a state machine
type Scheduler interface {
	// Execute executes a function asynchronously
	Execute(f func())

	// ScheduleOnce schedules a function to be run once after the given delay
	ScheduleOnce(delay time.Duration, f func()) Timer

	// ScheduleRepeat schedules a function to run repeatedly every interval starting after the given delay
	ScheduleRepeat(delay time.Duration, interval time.Duration, f func()) Timer

	// ScheduleIndex schedules a function to run at a specific index
	ScheduleIndex(index Index, f func())
}

// Timer is a cancellable timer
type Timer interface {
	// Cancel cancels the timer, preventing it from running in the future
	Cancel()
}

func newScheduler() *scheduler {
	return &scheduler{
		tasks:          list.New(),
		scheduledTasks: list.New(),
		indexTasks:     make(map[Index]*list.List),
		time:           time.Now(),
	}
}

type scheduler struct {
	Scheduler
	tasks          *list.List
	scheduledTasks *list.List
	indexTasks     map[Index]*list.List
	time           time.Time
}

func (s *scheduler) Execute(f func()) {
	s.tasks.PushBack(f)
}

func (s *scheduler) ScheduleOnce(delay time.Duration, f func()) Timer {
	task := &task{
		scheduler: s,
		time:      time.Now().Add(delay),
		interval:  0,
		callback:  f,
	}
	s.schedule(task)
	return task
}

func (s *scheduler) ScheduleRepeat(delay time.Duration, interval time.Duration, f func()) Timer {
	task := &task{
		scheduler: s,
		time:      time.Now().Add(delay),
		interval:  interval,
		callback:  f,
	}
	s.schedule(task)
	return task
}

func (s *scheduler) ScheduleIndex(index Index, f func()) {
	tasks, ok := s.indexTasks[index]
	if !ok {
		tasks = list.New()
		s.indexTasks[index] = tasks
	}
	tasks.PushBack(f)
}

// runImmediateTasks runs the immediate tasks in the scheduler queue
func (s *scheduler) runImmediateTasks() {
	task := s.tasks.Front()
	for task != nil {
		task.Value.(func())()
		task = task.Next()
	}
	s.tasks = list.New()
}

// runScheduleTasks runs the scheduled tasks in the scheduler queue
func (s *scheduler) runScheduledTasks(time time.Time) {
	s.time = time
	element := s.scheduledTasks.Front()
	if element != nil {
		complete := list.New()
		for element != nil {
			task := element.Value.(*task)
			if task.isRunnable(time) {
				next := element.Next()
				s.scheduledTasks.Remove(element)
				s.time = task.time
				task.run()
				complete.PushBack(task)
				element = next
			} else {
				break
			}
		}

		element = complete.Front()
		for element != nil {
			task := element.Value.(*task)
			if task.interval > 0 {
				task.time = s.time.Add(task.interval)
				s.schedule(task)
			}
			element = element.Next()
		}
	}
}

// runIndex runs functions pending at the given index
func (s *scheduler) runIndex(index Index) {
	tasks, ok := s.indexTasks[index]
	if ok {
		task := tasks.Front()
		for task != nil {
			task.Value.(func())()
			task = task.Next()
		}
		delete(s.indexTasks, index)
	}
}

// schedule schedules a task
func (s *scheduler) schedule(t *task) {
	if s.scheduledTasks.Len() == 0 {
		t.element = s.scheduledTasks.PushBack(t)
	} else {
		element := s.scheduledTasks.Back()
		for element != nil {
			time := element.Value.(*task).time
			if element.Value.(*task).time.UnixNano() < time.UnixNano() {
				t.element = s.scheduledTasks.InsertAfter(t, element)
				return
			}
			element = element.Prev()
		}
		t.element = s.scheduledTasks.PushFront(t)
	}
}

// Scheduler task
type task struct {
	Timer
	scheduler *scheduler
	interval  time.Duration
	callback  func()
	time      time.Time
	element   *list.Element
}

func (t *task) isRunnable(time time.Time) bool {
	return time.UnixNano() > t.time.UnixNano()
}

func (t *task) run() {
	t.callback()
}

func (t *task) Cancel() {
	if t.element != nil {
		t.scheduler.scheduledTasks.Remove(t.element)
	}
}
