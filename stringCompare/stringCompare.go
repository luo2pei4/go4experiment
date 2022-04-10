package main

import "fmt"

func main() {

	str1 := "abc"
	str2 := "a1c"

	if str1 > str2 {
		fmt.Printf("%s is bigger than %s\n", str1, str2) // print
	} else {
		fmt.Printf("%s is less than %s\n", str1, str2)
	}

	str1 = "abc"
	str2 = "aac"

	if str1 > str2 {
		fmt.Printf("%s is bigger than %s\n", str1, str2) // print
	} else {
		fmt.Printf("%s is less than %s\n", str1, str2)
	}

	str1 = "abc"
	str2 = "bbc"

	if str1 > str2 {
		fmt.Printf("%s is bigger than %s\n", str1, str2)
	} else {
		fmt.Printf("%s is less than %s\n", str1, str2) // print
	}

	str1 = "abc"
	str2 = "aec"

	if str1 > str2 {
		fmt.Printf("%s is bigger than %s\n", str1, str2)
	} else {
		fmt.Printf("%s is less than %s\n", str1, str2) // print
	}

	str1 = "abc"
	str2 = "abcd"

	if str1 > str2 {
		fmt.Printf("%s is bigger than %s\n", str1, str2)
	} else {
		fmt.Printf("%s is less than %s\n", str1, str2) // print
	}

	str1 = "aaa"
	str2 = "aaaa"

	if str1 > str2 {
		fmt.Printf("%s is bigger than %s\n", str1, str2)
	} else {
		fmt.Printf("%s is less than %s\n", str1, str2) // print
	}

	str1 = "bbb"
	str2 = "aaaa"

	if str1 > str2 {
		fmt.Printf("%s is bigger than %s\n", str1, str2) // print
	} else {
		fmt.Printf("%s is less than %s\n", str1, str2)
	}

	str1 = "baa1"
	str2 = "abaaaaa"

	if str1 > str2 {
		fmt.Printf("%s is bigger than %s\n", str1, str2) // print
	} else {
		fmt.Printf("%s is less than %s\n", str1, str2)
	}
}
