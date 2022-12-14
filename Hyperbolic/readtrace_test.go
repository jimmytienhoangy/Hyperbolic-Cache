package cache

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"testing"
// 	"time"
// )

// func trace_test(t *testing.T) {
// 	file, err := os.Open("test.tr")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// set the capacity
// 	capacity := 100 //(?)
// 	hyperbolic_cache := NewHyperbolicCache(capacity)

// 	defer file.Close()
// 	scanner := bufio.NewScanner(file)
// 	// optionally, resize scanner's capacity for lines over 64K, see next example
// 	for scanner.Scan() {

// 		text := scanner.Text()

// 		parsed_text := strings.Split(text, " ")

// 		_, key, size := parsed_text[0], parsed_text[1], parsed_text[2]

// 		_, ok := hyperbolic_cache.Get(key)

// 		if !ok {
// 			value, _ := strconv.Atoi(size)
// 			hyperbolic_cache.Set(key, value)
// 			time.Sleep(1 * time.Millisecond)
// 		}

// 		//fmt.Println(_ + "|" + key + "|" + size)

// 		//fmt.Println("test")
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(hyperbolic_cache.Stats())
// }
