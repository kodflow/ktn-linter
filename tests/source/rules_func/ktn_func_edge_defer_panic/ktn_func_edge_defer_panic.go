package badfuncdefer

import "fmt"

// Violations avec defer, panic, recover

// openResource sans documentation du defer
func openResource(name string) error {
	defer closeResource(name)

	if name == "" {
		panic("empty resource name")
	}

	return nil
}

func closeResource(name string) {
	fmt.Println("closing", name)
}

// processWithRecover sans doc du recover
func processWithRecover() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// Code qui peut paniquer
	riskyOperation()
	return nil
}

func riskyOperation() {
	panic("something went wrong")
}

// multipleDefers sans explication de l'ordre
func multipleDefers() {
	defer fmt.Println("first defer")
	defer fmt.Println("second defer")
	defer fmt.Println("third defer")

	fmt.Println("function body")
}

// deferInLoop anti-pattern non documenté
func deferInLoop(files []string) {
	for _, file := range files {
		f := openFile(file)
		defer closeFile(f)
	}
}

func openFile(name string) *File {
	return &File{name: name}
}

func closeFile(f *File) {
	fmt.Println("closing", f.name)
}

type File struct {
	name string
}

// complexDeferLogic defer avec logique complexe non commentée
func complexDeferLogic(shouldPanic bool) error {
	defer func() {
		if r := recover(); r != nil {
			if shouldPanic {
				fmt.Println("expected panic")
			} else {
				panic(r)
			}
		}
	}()

	if shouldPanic {
		panic("intentional panic")
	}

	return nil
}
