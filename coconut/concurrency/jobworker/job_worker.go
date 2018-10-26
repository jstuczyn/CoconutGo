// job_worker.go - Worker for job queue tasks.
// Copyright (C) 2018  Jedrzej Stuczynski.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package jobworker implements worker for performing tasks defined by jobpacket that are in the queue.
package jobworker

import (
	"fmt"

	"github.com/jstuczyn/CoconutGo/coconut/concurrency/jobpacket"
	"github.com/jstuczyn/CoconutGo/coconut/concurrency/worker"
	"github.com/jstuczyn/CoconutGo/constants"
)

// Worker is an instance of jobWorker.
type Worker struct {
	worker.Worker

	id       uint64
	jobQueue <-chan interface{}
}

func (w *Worker) worker() {
	for {
		var jobpkt *jobpacket.JobPacket
		select {
		case <-w.HaltCh():
			if constants.DEBUG {
				fmt.Println("Halting worker", w.id)
			}
			return
		case e := <-w.jobQueue:
			jobpkt = e.(*jobpacket.JobPacket)
		}

		res, err := jobpkt.Op()
		// job provider will be able to distinguish those cases thanks to type assertions
		if err != nil {
			jobpkt.OutCh <- err
		} else {
			jobpkt.OutCh <- res
		}
	}
}

// New creates new instance of a jobWorker.
func New(jobQueue <-chan interface{}, id uint64) *Worker {
	w := &Worker{
		jobQueue: jobQueue,
		id:       id,
	}

	w.Go(w.worker)
	return w
}
