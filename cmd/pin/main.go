package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/osm/pin"
)

func main() {
	// Define command line flags
	g := flag.Bool("generate", false, "generate a personal identity number")
	d := flag.String("date", "", "generate last four digits of a personal identity number based on the given date")
	v := flag.String("valid", "", "check whether the given personal identity number is valid")
	m := flag.String("male", "", "check if the given personal identity number is male")
	f := flag.String("female", "", "check if the given personal identity number is female")
	flag.Parse()

	// Local variables for val and errors.
	var val string
	var err error

	// Pass flags to the appropriate pin function.
	if *g {
		val, err = pin.Generate()
	} else if *d != "" {
		val, err = pin.GenerateFromDate(*d)
	} else if *v != "" {
		_, err = pin.IsValid(*v)
		if err == nil {
			goto end
		}
	} else if *m != "" {
		_, err = pin.IsMale(*m)
		if err == nil {
			goto end
		}
	} else if *f != "" {
		_, err = pin.IsFemale(*f)
		if err == nil {
			goto end
		}
	} else {
		flag.Usage()
	}

	// Make sure that err isn't nil.
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Write the value of val to stdout.
	fmt.Fprintf(os.Stdout, "%v\n", val)

end:
	// Exit successfully.
	os.Exit(0)
}
