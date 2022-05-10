package goroutines

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func AddDataIntoChannel(ch chan string, name string) {
	ch <- name
}


func TestCreateChannel(t *testing.T) {
	// channels are a medium that the goroutines use in order to communicate effectively
	ch := make(chan string) // only can receive string value

	//ch <- go "kaguya chan" // send data into channel
	//value := <- ch // receive data from channel

	go AddDataIntoChannel(ch, "kaguya chan") // must be goroutines, and have sender and receiver

	value := <- ch

	defer close(ch) // always close channel
	fmt.Println(value)
}

// use type struct for type data channel

type User interface {
	CreateAccount(ch chan Account, account Account)
	GetAccount(ch chan Account)
}
type Account struct {
	Id int
	Username string
	Name string
	Password string
}

func CreateAccount(ch chan Account, account Account) {
	ch <- account
}

func (account *Account)GetAccount(ch chan <- Account, id int) { // <- only receive data
	if id == account.Id {
		ch <- *account
	} else {
		ch <- Account{}
	}

}
func TestCreateChannelStructType(t *testing.T) {
	chUser := make(chan Account)

	account := Account{
		Id:       10,
		Username: "hikaru",
		Name:     "hikaruch",
		Password: "wakaranai",
	}

	go CreateAccount(chUser, account)

	go account.GetAccount(chUser, 10)


	// close channel
	value, ok := <-chUser // check chUser if ok send data into value <-
		if !ok {
		defer close(chUser)
	}

	fmt.Println(value)
}

func ReadDataOnly(chData <- chan Account)  {
	acoount := <- chData // receive only
	time.Sleep(time.Second)
	fmt.Println(acoount)
}
func SendDataOnly(ch chan <- Account, sender ...*Account) {
	for _, acc := range sender {
		if (Account{}) != *acc { // check empty
			ch <- *acc
		}
		time.Sleep(time.Second)
	}
}

func CreateNewAccount(username string, name string, password string) *Account{
	return &Account{
		Id:     rand.Intn(1000)  ,
		Username: username,
		Name:     name,
		Password: password,
	}
}

func TestManagaeAccount(t *testing.T)  {
	ch := make(chan Account, 3) // create channel with buffer 3
	defer close(ch)

	account1 := CreateNewAccount("hikaru", "hikako", "wakaranai")
	account2 := CreateNewAccount("micchon", "micchon", "wakaranai")
	account3 := CreateNewAccount("kaguya", "kaguya", "wakaranai")
	account4 := CreateNewAccount("rem", "rem chan", "wakaranai") // will be wait until the data in the buffer is taken
	account := []*Account{account1, account2, account3,account4}
	go SendDataOnly(ch, account...)

	go ReadDataOnly(ch) // get data from buffer
	go ReadDataOnly(ch) // get data from buffer

	time.Sleep(time.Second)

}

type Anime struct {
	Title string
	Rating int
}

func BenchmarkWriteAndReceeiveManyData(t *testing.B) {
	// use for read and receive many data from channel
	// use range

	animeList := make(chan Anime, 1) // buffer 1

	go func() {
		for i := 0 ; i < 1000000; i++ {
			animeList <- Anime{
				Title:  "Anime iterate " + strconv.Itoa(i),
				Rating: i,
			}
		}

		close(animeList) // always close the channel after using it
	}()

	for anime := range animeList { // range channel
		log.Println(anime)
	}
}


func TestSelectChannel(t *testing.T) {
	// use for comunicate 2 channel or more
	channel1 := make(chan Account)
	channel2 := make(chan Account)
	defer close(channel1)
	defer close(channel2)

	account1 := CreateNewAccount("hikaru", "hikako", "wakaranai")
	account2 := CreateNewAccount("micchon", "micchon", "wakaranai")
	tempAccount1 := []*Account{account1}
	tempAccount2 := []*Account{account2}

	go SendDataOnly(channel1, tempAccount1...)
	go SendDataOnly(channel2, tempAccount2...)

	time.Sleep(time.Second)
	// if you want read data from channel use select with for loop

	counter := 0
	for  {
		select {
		case data := <- channel1:
			fmt.Println(data)
			counter++
		case data := <- channel2:
			fmt.Println(data)
			counter++
		default:
			fmt.Println("wait data...") // do something if the data doesn't exist
		}
		if counter == 2 {break}
	}
}





















