package main

import (
	"errors"
)

type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

type UsersWrapper struct {
	Data []User `json:"data"`
}

var SumSubError = errors.New("an error occurred while searching for personal information on SumSub")

func main() {
	//err1 := errors.New("err1")
	//err2 := errors.New("err2")
	//err3 := errors.New("err3")
	//
	//err := errors.Join(err1, err2, err3)
	//fmt.Println(err.Error())

	//var g errgroup.Group
	//
	//g.Go(func() error {
	//	return errors.New("test Group")
	//})
	//
	//g.Go(func() error {
	//	time.Sleep(1 * time.Second)
	//	return nil
	//})
	//
	//if err := g.Wait(); err != nil {
	//	fmt.Println(err.Error())
	//}
	//e := fmt.Errorf("%w:can not find any information about occupation in the sumsub", SumSubError)
	//
	//fmt.Println(e.Error())
	//
	////fmt.Println(errors.Is(e, SumSubError))
	//timeStr := "2023-12-14T22:03:05.000000Z"
	//parsedTime, _ := time.Parse(time.RFC3339, timeStr)
	//timestamp := timestamppb.New(parsedTime)
	//fmt.Println(parsedTime)
	//fmt.Println(timestamp)
}
