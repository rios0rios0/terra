package entities

type OSLinux struct{}

func GetOS() *OSLinux {
	return &OSLinux{}
}
