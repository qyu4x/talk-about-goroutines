package goroutines

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// mutex use for solve problem race condition in golang
// example
func IncrementValue(value *int, m *sync.Mutex, wg *sync.WaitGroup ) {
	m.Lock() // acquire  lock

	*value++ // do condition

	m.Unlock() // release
	wg.Done() // done

}
func TestSolveProblemRaceConditidionWithMuatex(t *testing.T) {
	var wg = new(sync.WaitGroup) // auto pointer
	var mutex = new(sync.Mutex)

	value := 0
	for i := 1 ; i <= 1000; i++ {
		for k := 1 ; k <= 100; k++ {
			wg.Add(1)
			go IncrementValue(&value, mutex, wg)
		}
	}

	wg.Wait()

	fmt.Println(value)

}

// RWmutex use for lock read and write data in same time
// componen Lock, Unlock, Rlock , RUnlock

type BankAccount struct {
	RWmutec sync.RWMutex
	Balance int
}

func (account *BankAccount)AddBalance(amount int) {
	account.RWmutec.Lock()
	account.Balance = account.Balance + amount
	account.RWmutec.Unlock()
}

func (account *BankAccount)GetBalance() int {
	account.RWmutec.RLock()
	balance := account.Balance
	account.RWmutec.RUnlock()
	return balance
}

func TestBankAccountRWmutex(t *testing.T) {
	account := BankAccount{}
	for i := 1 ; i <= 100; i++ {
		go func() {
			for z := 1 ; z <= 100; z++ {
				account.AddBalance(1)
				fmt.Println(account.GetBalance())
			}
		}()
	}

	time.Sleep(1* time.Second)
	fmt.Println("final balance ", account.GetBalance())
}


var counter = 0
func Incremment() {
	counter += 1
}

func TestSyncOnce(t *testing.T) {
	once := sync.Once{}
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		go func() {
			wg.Add(1)
			once.Do(Incremment) // only execute once(sekali doang di eksekusinya)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("value ", counter)
}

type Person struct { // dummy struct
	Name string
}
func TestSyncPool(t *testing.T) {
	pool := sync.Pool{
		// default value if data in pool nil
		New: func() interface{} {
			return new(Person)
		},
	}

	wg := sync.WaitGroup{}
	// add data to pool use pool.Put(data)
	// get data from pool use pool.Get()

	// get data from pool
	newPerson := pool.Get().(*Person)
	newPerson.Name = "hikaru"

	newPerson2 := pool.Get().(*Person)
	newPerson2.Name = "kaguya"


	// return data from pool
	pool.Put(newPerson)
	pool.Put(newPerson2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			data := pool.Get()
			fmt.Println(data)

			// return data from pool
			// reusable by another routine
			pool.Put(data)
			wg.Done()
		}()
	}
	wg.Wait()
}

// sync map, like a normal map but safe from race conditions
func TestSyncMap(t *testing.T) {
	data := sync.Map{}

	account := []Person{
		{Name: "hikaru"},
		{Name: "kaguya chan"},
		{Name: "micchon"},
		{Name: "rem chan"},
	}

	// load acoount  into map
	for i := 0; i < len(account); i ++ {
		data.Store("account" + strconv.Itoa(i), account[i]) // store data account
	}

	// get all data account
	data.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})

	// use load to get 1 data. use delete to delete data
	if _, ok := data.Load("account1"); ok {
		data.Delete("account1") // delete
	}
}

// this code has a bug
type UserOtp interface {
	Verify(otpCode int, locker *sync.Cond) bool
}

type CurrentUser struct {
	ProccessRegister VerifyUser
	Status bool
}
func (user *CurrentUser)Verify(otpCode int, locker *sync.Cond) bool {
	if user.ProccessRegister.CodeOtp == otpCode {
		locker.Broadcast() // next execution
		user.Status = true
		return true
	} else {
		locker.Signal()
		return false
	}
}

// sync cond
type VerifyUser struct {
	Name string
	CodeOtp int
}

type NewUser struct {
	User VerifyUser
}

func (verifyUser *VerifyUser)VerifyRegisterUser(ch chan NewUser, verify *VerifyUser, wg *sync.WaitGroup, otpCode int) {
	mu := sync.Mutex{}
	locker := sync.Cond{ // set locker cond L as mu -> sync mutex
		L: &mu,
	}

	currentUser := CurrentUser{
		ProccessRegister:VerifyUser{
			Name:    verify.Name,
			CodeOtp: verify.CodeOtp,
		},
		Status:           false,
	}

	locker.L.Lock()
	if ok := currentUser.Verify(otpCode, &locker); ok {
		locker.Wait()
		if okk := currentUser.Status; okk {
			ch <- NewUser{User: VerifyUser{
				Name:    verify.Name,
				CodeOtp: verify.CodeOtp,
			}}
		}
	}

	locker.L.Unlock()
	ch <- NewUser{}
	wg.Done()
}

func TestVerifyAccount(t *testing.T)  {
	chAccount := make(chan NewUser)
	verifyNewAcccount := VerifyUser{
		Name:    "hikaru",
		CodeOtp: 59121,
	}

	otpCodeReceiver := 59121

	wg := sync.WaitGroup{}
	wg.Add(1)

	go verifyNewAcccount.VerifyRegisterUser(chAccount, &verifyNewAcccount, &wg, otpCodeReceiver )


	data := <-chAccount
	if chAccount != nil {
		fmt.Println(data)
	}
	wg.Wait()
}

// sync atomic use for manipulating data primitive(save) without exception race condition, like number in goroutines read doc
func TestAtomic(t *testing.T) {
	var v int32 = 0
	wg := sync.WaitGroup{}
	for i := 1; i <= 10000; i++ {
		wg.Add(1)
		go func() {
			atomic.AddInt32(&v, 1)
			fmt.Println(v)
			wg.Done()
		}()

	}

	wg.Wait()
}