package mr

func mrCoordinator() {
	c := Coordinator{}
	c.Router()
	c.Run()
}
