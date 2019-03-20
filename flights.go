package airstorage

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Flight struct {
	ArrivalTime   time.Time   //Время прибытия
	DepartuteTime time.Time   //Время отбытия
	Type          string      //Тип воздушного судна
	Number        int         //Номер воздушного судна
	Cities        *CitiesList //Указатель на список городов рейса
	Prev, Next    *Flight     //Указатели на следующий/предыдущий рейс
}

//Двусвязанный список рейсов гражданской авиации
type FlightsList struct {
	Head, Tail *Flight //Указатели на начало и конец списка
	Count      int     //Кол-во элементов в списке
}

func (flightsList *FlightsList) Get(number int) *Flight {
	var currentNumber int
	var flight *Flight = new(Flight)
	if number > flightsList.Count/2 {
		currentNumber = flightsList.Count - 1
		flight = flightsList.Tail
		for flight != nil {
			if currentNumber == number {
				return flight
			}
			flight = flight.Prev
			currentNumber--
		}
	} else {
		flight = flightsList.Head
		for flight != nil {
			if currentNumber == number {
				return flight
			}
			flight = flight.Next
			currentNumber++
		}
	}
	return nil
}

func (flight *Flight) Set(data map[string]interface{}) {
	err := mapstructure.Decode(data, flight)
	if err != nil {
		panic(err)
	}
}

func (flightsList *FlightsList) Insert(data map[string]interface{}, pos int) *Flight {
	if pos > flightsList.Count {
		panic("pos > than list Count")
	} else {
		var flight *Flight = new(Flight)
		*flight = Flight{}
		flight.Set(data)
		flight.Cities = &CitiesList{}
		if pos == 0 {
			next := flightsList.Head
			flightsList.Head = flight
			flight.Next = next
			if flight.Next == nil {
				flightsList.Tail = flight
			} else {
				flight.Next.Prev = flight
			}
		} else if pos == flightsList.Count {
			prev := flightsList.Tail
			flightsList.Tail = flight
			flight.Prev = prev
			flight.Prev.Next = flight
		} else {
			next := flightsList.Get(pos)
			prev := next.Prev
			flight.Next = next
			flight.Prev = prev
			flight.Prev.Next = flight
			flight.Next.Prev = flight

		}
		flightsList.Count++
		return flight
	}
}

func (flightsList *FlightsList) Push(data map[string]interface{}) *Flight {
	return flightsList.Insert(data, flightsList.Count)
}

func (flightsList *FlightsList) Unhift(data map[string]interface{}) *Flight {
	return flightsList.Insert(data, 0)
}

func (flightsList *FlightsList) Remove(pos int) *Flight {
	if pos >= flightsList.Count {
		panic("pos >= than list Count")
	} else {
		var flight *Flight = new(Flight)
		flight = flightsList.Get(pos)
		if pos == 0 {
			head := flightsList.Head.Next
			if head != nil {
				head.Prev = nil
			} else {
				flightsList.Tail = head
			}
			flightsList.Head = head
		} else if pos == flightsList.Count-1 {
			tail := flightsList.Tail.Prev
			tail.Next = nil
			flightsList.Tail = tail
		} else {
			prev := flight.Prev
			next := flight.Next
			next.Prev = prev
			prev.Next = next
		}
		flightsList.Count--
		return flight
	}
}

func (flightsList *FlightsList) Pop() *Flight {
	return flightsList.Remove(flightsList.Count - 1)
}

func (flightsList *FlightsList) Shift() *Flight {
	return flightsList.Remove(0)
}

func (flight *Flight) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ArrivalTime   string      `json:"arrival_time"`   //Время прибытия
		DepartuteTime string      `json:"departute_time"` //Время отбытия
		Type          string      `json:"type"`           //Тип воздушного судна
		Number        int         `json:"number"`         //Номер воздушного судна
		Cities        *CitiesList `json:"cities"`         //Указатель на список городов рейса
	}{
		ArrivalTime:   flight.ArrivalTime.String(),
		DepartuteTime: flight.DepartuteTime.String(),
		Type:          flight.Type,
		Number:        flight.Number,
		Cities:        flight.Cities,
	})
}

func (flightList *FlightsList) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	ptr := flightList.Head
	for ptr != nil {
		jsonValue, err := json.Marshal(ptr)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(string(jsonValue))
		if ptr.Next != nil {
			buffer.WriteString(",")
		}
		ptr = ptr.Next
	}
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}
