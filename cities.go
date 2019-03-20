package airstorage

import (
	"bytes"
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

//Город
type City struct {
	Name string //Наименование города
	Next *City  //Указатель на следующий город в односвязном списке
}

//Одновязанный список городов
type CitiesList struct {
	Head  *City //Указатель на начало списка
	Count int   //Кол-во элементов в списке
}

func (citiesList *CitiesList) Get(number int) *City {
	var currentNumber int
	var city *City = new(City)
	city = citiesList.Head
	for city != nil {
		if currentNumber == number {
			return city
		}
		city = city.Next
		currentNumber++
	}
	return nil
}

func (city *City) Set(data map[string]interface{}) {
	err := mapstructure.Decode(data, city)
	if err != nil {
		panic(err)
	}
}

func (citiesList *CitiesList) Insert(data map[string]interface{}, pos int) *City {
	if pos > citiesList.Count {
		panic("pos > than list count")
	} else {
		var city *City = new(City)
		city.Set(data)
		if pos == 0 {
			citiesList.Head, city.Next = city, citiesList.Head
		} else if pos == citiesList.Count {
			citiesList.Get(pos - 1).Next = city
		} else {
			citiesList.Get(pos-1).Next, city.Next = city, citiesList.Get(pos)
		}
		citiesList.Count++
		return city
	}
}

func (citiesList *CitiesList) Push(data map[string]interface{}) *City {
	return citiesList.Insert(data, citiesList.Count)
}

func (citiesList *CitiesList) Unhift(data map[string]interface{}) *City {
	return citiesList.Insert(data, 0)
}

func (citiesList *CitiesList) Remove(pos int) *City {
	if pos >= citiesList.Count {
		panic("pos >= than list count")
	} else {
		var city *City = new(City)
		city = citiesList.Get(pos)
		if pos == 0 {
			head := citiesList.Head.Next
			citiesList.Head = head
		} else if pos == citiesList.Count-1 {
			citiesList.Get(pos - 1).Next = nil
		} else {
			prev := citiesList.Get(pos - 1)
			prev.Next = prev.Next.Next
		}
		citiesList.Count--
		return city
	}
}

func (citiesList *CitiesList) Pop() *City {
	return citiesList.Remove(citiesList.Count - 1)
}

func (citiesList *CitiesList) Shift() *City {
	return citiesList.Remove(0)
}

func (city *City) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name string `json:"name"`
	}{
		Name: city.Name,
	})
}

func (citiesList *CitiesList) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	ptr := citiesList.Head
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
