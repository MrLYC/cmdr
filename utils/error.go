package utils

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func CallClose(closer interface {
	Close() error
}) {
	CheckError(closer.Close())
}
