package main

func YourName(name string, aliases ...string) {
	if len(aliases) > 0 {
		println(aliases[0])
	}
}

func YourNameInvoke() {
	YourName("Tony")
	YourName("Tony", "Tony1")
	YourName("Tony", "Tony2", "Tony1")
}
