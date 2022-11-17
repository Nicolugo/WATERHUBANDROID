package window

// Window is
type Window interface {
	Render(imageData []byte)
	Run(run fun