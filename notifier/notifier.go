package notifier

import (
	"fmt"
)

type NotifierService struct {
}

type Notifier interface {
	Notify(string, int)
}

func (n NotifierService) Notify(todoText string, userId int) {
	fmt.Println("notified")
}
