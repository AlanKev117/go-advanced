package main

import "fmt"

type Observer interface {
	getId() string
	updateValue(string)
}

type Topic interface {
	register(Observer)
	broadcast()
}

type Item struct {
	observers []Observer
	name      string
	available bool
}

func NewItem(name string) *Item {
	return &Item{
		name: name,
	}
}

func (i *Item) UpdateAvailable() {
	fmt.Printf("Item %v is now available.\n", i.name)
	i.available = true
	i.broadcast()
}

func (i *Item) register(observer Observer) {
	i.observers = append(i.observers, observer)
}

func (i *Item) broadcast() {
	for _, observer := range i.observers {
		observer.updateValue(i.name)
	}
}

type EmailClient struct {
	id string
}

func (ec *EmailClient) getId() string {
	return ec.id
}

func (ec *EmailClient) updateValue(value string) {
	fmt.Printf("Sending email about %v from client %v\n", value, ec.id)
}

func main() {
	nvidiaItem := NewItem("RTX 3080")
	firstObserver := &EmailClient{
		id: "0x30",
	}
	secondObserver := &EmailClient{
		id: "0x45",
	}
	nvidiaItem.register(firstObserver)
	nvidiaItem.register(secondObserver)
	nvidiaItem.UpdateAvailable()
}
