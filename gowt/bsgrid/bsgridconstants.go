package bsgrid

import (
	"fmt"
	"strings"
)

//Row is for new rows in the grid
const Row = "row"

//VerySmall for phones
const VerySmall = "xs"

//Small for tablets
const Small = "sm"

//Medium for desktops
const Medium = "md"

//Large for larger desktops
const Large = "lg"

//Cell returns a bootstrap grid cell class (each row has 12 columns)
func Cell(colspan int, size string) string {

	return strings.Join([]string{"col", size, fmt.Sprint(colspan)}, "-")

}
