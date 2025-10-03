package main

import (
	"testing"
)

func TestUnpack1(t *testing.T) {
	testString := "a4bc2d5e"
	expected := "aaaabccddddde"

	result, _ := UnpackString(testString)

	if result != expected {
		t.Errorf("UnpackString(\"%s\") returned: \"%s\", expeted \"%s\"", testString, result, expected)
	}
}

func TestUnpack2(t *testing.T) {
	testString := "abcd"
	expected := "abcd"

	result, _ := UnpackString(testString)

	if result != expected {
		t.Errorf("UnpackString(\"%s\") returned: \"%s\", expeted \"%s\"", testString, result, expected)
	}
}

func TestUnpack3(t *testing.T) {
	testString := "42"
	expected := ""

	result, err := UnpackString(testString)

	if err == nil || result != expected {
		t.Errorf("UnpackString(\"%s\") returned: \"%s\", expeted \"%s\"", testString, result, expected)
	}
}

func TestUnpack4(t *testing.T) {
	testString := ""
	expected := ""

	result, _ := UnpackString(testString)
	if result != expected {
		t.Errorf("UnpackString(\"%s\") returned: \"%s\", expeted \"%s\"", testString, result, expected)
	}
}

func TestUnpack5(t *testing.T) {
	testString := "s1g1h2v3"
	expected := "sghhvvv"

	result, _ := UnpackString(testString)
	if result != expected {
		t.Errorf("UnpackString(\"%s\") returned: \"%s\", expeted \"%s\"", testString, result, expected)
	}
}

func TestUnpack6(t *testing.T) {
	testString := "qwe\\4\\5"
	expected := "qwe45"

	result, _ := UnpackString(testString)
	if result != expected {
		t.Errorf("UnpackString(\"%s\") returned: \"%s\", expeted \"%s\"", testString, result, expected)
	}
}

func TestUnpack7(t *testing.T) {
	testString := "qwe\\45"
	expected := "qwe44444"

	result, _ := UnpackString(testString)
	if result != expected {
		t.Errorf("UnpackString(\"%s\") returned: \"%s\", expeted \"%s\"", testString, result, expected)
	}
}
