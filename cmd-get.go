package main

func get(name string, filter []string) {
	v := load(name, false, filter)
	if v == nil {
		return
	}

	v.AddNames(true)
	v.Update()
}
