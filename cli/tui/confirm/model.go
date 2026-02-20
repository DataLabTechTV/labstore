package confirm

type Model[T any, U any] struct {
	Prompt       string
	Width        int
	Height       int
	HoverConfirm bool
	ConfirmMsg   T
	CancelMsg    U
}
