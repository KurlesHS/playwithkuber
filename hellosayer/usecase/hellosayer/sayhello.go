package hellosayer

import (
	"context"
	"math/rand"
)

func (HelloSayer) SayHello(_ context.Context, name string) string {
	//создай слайс из строк 'привет' на разных языках
	hello := []string{
		"Ahlan wa sahlan",
		"Marhaba",
		"Hola",
		"Привет",
		"Прывитание",
		"Здравейте",
		"Jo napot",
		"Chao",
		"Aloha",
		"Hallo",
		"Geia sou",
		"Shalom",
		"Buenas dias",
		"Bonjour",
		"Gutten tag",
		"Ave",
		"Lab dien, sveiki",
		"Sawatdi",
		"Namaste",
	}
	// выбери случайную строку из массива
	return hello[rand.Intn(len(hello))] + ", " + name + "!"
}
