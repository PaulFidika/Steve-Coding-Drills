package monkey

func AddStuff(num1 int, num2 int) int {
	var _ = multiplyStuff(num1, num2)

	return num1 + num2
}

func multiplyStuff(num1 int, num2 int) int {
	return num1 * num2
}
