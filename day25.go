package adventofcode2017

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TuringMachineDirection int

const (
	TmRight = TuringMachineDirection(0)
	TmLeft  = TuringMachineDirection(1)
)

type TuringMachineStateName string

type TuringMachineInstruction struct {
	Write     int
	Move      TuringMachineDirection
	NextState TuringMachineStateName
}

type TuringMachineState struct {
	Branch [2]TuringMachineInstruction
}

type TuringMachine struct {
	nextState      TuringMachineStateName
	stepsRemaining int
	position       int
	states         map[TuringMachineStateName]TuringMachineState
	tape           map[int]int // position → written value
}

func NewTuringMachine(blueprint_raw string) *TuringMachine {
	blueprint := strings.Split(blueprint_raw, "\n")
	tm := TuringMachine{states: make(map[TuringMachineStateName]TuringMachineState), tape: make(map[int]int)}

	var re *regexp.Regexp
	var line string
	var match []string

	// preamble
	re = regexp.MustCompile(`Begin in state (\w+)\.`)
	line = blueprint[0]
	match = re.FindStringSubmatch(line)
	if match == nil {
		panic(fmt.Sprintf("could not parse %q", line))
	}
	tm.nextState = TuringMachineStateName(match[1])

	re = regexp.MustCompile(`Perform a diagnostic checksum after (\d+) steps\.`)
	line = blueprint[1]
	match = re.FindStringSubmatch(line)
	if match == nil {
		panic(fmt.Sprintf("could not parse %q", line))
	}
	tm.stepsRemaining, _ = strconv.Atoi(match[1])

	// repeating state sections
	stateRe := regexp.MustCompile(`In state (\w+):`)
	writeRe := regexp.MustCompile(`Write the value (\d+)`)
	moveRe := regexp.MustCompile(`Move one slot to the (\w+)`)
	continueRe := regexp.MustCompile(`Continue with state (\w+)`)

	jline := 3
	for jline < len(blueprint) && stateRe.MatchString(blueprint[jline]) {
		tms := TuringMachineState{}

		match = stateRe.FindStringSubmatch(blueprint[jline])
		if match == nil {
			panic(fmt.Sprintf("could not parse %q on line", blueprint[jline], jline))
		}
		state := TuringMachineStateName(match[1])

		jline++
		for jcurr := 0; jcurr <= 1; jcurr++ {
			jline++
			match = writeRe.FindStringSubmatch(blueprint[jline])
			if match == nil {
				panic(fmt.Sprintf("could not parse %q on line", blueprint[jline], jline))
			}
			tms.Branch[jcurr].Write, _ = strconv.Atoi(match[1])

			jline += 1
			match = moveRe.FindStringSubmatch(blueprint[jline])
			if match == nil {
				panic(fmt.Sprintf("could not parse %q on line", blueprint[jline], jline))
			}
			switch match[1] {
			case "left":
				tms.Branch[jcurr].Move = TmLeft
			case "right":
				tms.Branch[jcurr].Move = TmRight
			default:
				panic(fmt.Sprintf("could not figure out direction %q", match[1]))
			}

			jline += 1
			match = continueRe.FindStringSubmatch(blueprint[jline])
			if match == nil {
				panic(fmt.Sprintf("could not parse %q on line", blueprint[jline], jline))
			}
			tms.Branch[jcurr].NextState = TuringMachineStateName(match[1])

			jline++
		}

		tm.states[state] = tms
		jline++
	}

	return &tm
}

func (tm *TuringMachine) NextState() string {
	return string(tm.nextState)
}

func (tm *TuringMachine) Position() int {
	return tm.position
}

func (tm *TuringMachine) TapeAt(position int) int {
	rval, ok := tm.tape[position]
	if !ok {
		return 0
	}
	return rval
}

func (tm *TuringMachine) StepsRemaining() int {
	return tm.stepsRemaining
}

func (tm *TuringMachine) Checksum() int {
	count := 0
	for _, value := range tm.tape {
		if value == 1 {
			count++
		}
	}
	return count
}

func (tm *TuringMachine) Step() {
	state := tm.states[tm.nextState]
	instructions := state.Branch[tm.TapeAt(tm.Position())]
	tm.tape[tm.Position()] = instructions.Write
	tm.nextState = instructions.NextState
	if instructions.Move == TmRight {
		tm.position += 1
	} else {
		tm.position -= 1
	}
	tm.stepsRemaining -= 1
}

func (tm *TuringMachine) Run() {
	for tm.stepsRemaining > 0 {
		tm.Step()
	}
}

func (tm *TuringMachine) State(name string) TuringMachineState {
	return tm.states[TuringMachineStateName(name)]
}
