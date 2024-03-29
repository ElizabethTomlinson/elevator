package elevator

import (
	"elevator/elevatorDirections"
	"errors"
	"fmt"
	"math"
	"sort"
)

const Speed = 0.5

type Elevator struct {
	floorQueue    []int
	floorTarget   int
	currentHeight float64
	people        []Agent
	maxFloor      int
	doorsOpen     bool
	direction     elevatorDirections.Direction
}

type Agent interface {
	GetCurrentFloor() int
	GetDesiredFloor() int
	SetCurrentFloor(int)
	GetId() string
}

func New(maxFloor int) (*Elevator, error) {
	var e Elevator
	if maxFloor <= 0 {
		return nil, errors.New("maxFloor must be > 0")
	}
	e.floorTarget = 0
	e.maxFloor = maxFloor
	e.direction = elevatorDirections.UP
	return &e, nil
}

func (e *Elevator) CallElevator(a Agent) {
	e.SetTargetFloor(a.GetCurrentFloor())
}

func (e *Elevator) Board(a Agent) {
	if e.doorsOpen && e.isAtAgentFloor(a) {
		e.people = append(e.people, a)
		e.SetTargetFloor(a.GetDesiredFloor())
		a.SetCurrentFloor(-1)
		fmt.Printf("  ** %v has boarded the elevator ** \n", a.GetId())
	}
}

func (e *Elevator) Exit(a Agent) {
	if e.IsInElevator(a) {
		if e.doorsOpen {
			if e.isAtAgentsTargetFloor(a) {
				a.SetCurrentFloor(a.GetDesiredFloor())
				fmt.Printf("  ** %v has left the elevator ** \n", a.GetId())
			}
		}
	}
}

func (e *Elevator) Move() {
	if e.isAtTargetFloor() {
		e.doorsOpen = true
		e.nextTargetFloor()
		return
	}
	e.doorsOpen = false
	if e.currentHeight < float64(e.floorTarget) {
		e.goUp()
		return
	}

	if e.currentHeight > float64(e.floorTarget) {
		e.goDown()
		return
	}
}

func (e *Elevator) SetTargetFloor(floor int) {
	if floor == e.floorTarget {
		return
	}
	if floor <= e.maxFloor && floor >= 0 {
		for _, f := range e.floorQueue {
			if f == floor {
				return
			}
		}
		e.floorQueue = append(e.floorQueue, floor)
	}
	e.sortFloors()
}

func (e Elevator) GetMaxFloor() int {
	return e.maxFloor
}

func (e Elevator) IsInElevator(a Agent) bool {
	for _, person := range e.people {
		if a.GetId() == person.GetId() {
			return true
		}
	}
	return false
}

func (e Elevator) GetHeight() float64 {
	return e.currentHeight
}

func (e Elevator) PrintStatus() {
	var msg string
	msg = fmt.Sprintf("   Elevator currently at height %v", e.GetHeight())
	fmt.Println(msg)
	msg = fmt.Sprintf("   Elevator moving to %v with queue of %v", e.floorTarget, e.floorQueue)
	fmt.Println(msg)
	msg = fmt.Sprintf("   Elevator direction %v", e.direction.ToString())
	fmt.Println(msg)
}

///////////////////////////
///////// Private /////////
///////////////////////////

func (e *Elevator) nextTargetFloor() {
	if len(e.floorQueue) == 1 {
		e.floorTarget = e.floorQueue[0]
		e.floorQueue = []int{}
		return
	}
	if len(e.floorQueue) > 1 {
		e.getNextFloorWithDirection()
		return
	}
	e.floorTarget = 0
}

func (e *Elevator) isAtAgentFloor(a Agent) bool {
	return e.isAtFloor(a.GetCurrentFloor())
}

func (e *Elevator) goUp() {
	if e.currentHeight <= float64(e.maxFloor) {
		e.currentHeight += Speed
	}
}

func (e *Elevator) goDown() {
	if e.currentHeight >= 0 {
		e.currentHeight -= Speed
	}
}

func (e *Elevator) isAtTargetFloor() bool {
	return e.isAtFloor(e.floorTarget)
}

func (e *Elevator) isAtFloor(floor int) bool {
	if math.Abs(e.currentHeight-float64(floor)) < Speed {
		e.currentHeight = float64(floor)
		return true
	}
	return false
}

func (e *Elevator) isAtAgentsTargetFloor(a Agent) bool {
	return e.isAtFloor(a.GetDesiredFloor())
}

func (e *Elevator) sortFloors() {
	sort.Ints(e.floorQueue)
}

func (e *Elevator) switchDirections() {
	if e.direction == elevatorDirections.UP {
		e.direction = elevatorDirections.DOWN
	} else {
		e.direction = elevatorDirections.UP
	}
}

func (e *Elevator) getNextFloorWithDirection() {
	needToSwitchDir := true
	if e.direction == elevatorDirections.DOWN {
		for i := len(e.floorQueue) - 1; i >= 0; i-- {
			if float64(e.floorQueue[i]) < e.currentHeight {
				e.floorQueue[i], e.floorQueue[0] = e.floorQueue[0], e.floorQueue[i]
				needToSwitchDir = false
				break
			}
		}
	} else {
		for i, floor := range e.floorQueue {
			if float64(floor) > e.currentHeight {
				e.floorQueue[i], e.floorQueue[0] = e.floorQueue[0], e.floorQueue[i]
				needToSwitchDir = false
				break
			}
		}
	}

	if needToSwitchDir {
		e.switchDirections()
	}

	e.floorTarget = e.floorQueue[0]
	e.floorQueue = e.floorQueue[1:]
	e.sortFloors()
}
