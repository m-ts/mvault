package helpnote

import (
	"fmt"
	"os"
	"text/tabwriter"
)

var helpNote = `
 __       __  __     __   ______   __    __  __     ________ 
|  \     /  \|  \   |  \ /      \ |  \  |  \|  \   |        \
| $$\   /  $$| $$   | $$|  $$$$$$\| $$  | $$| $$    \$$$$$$$$
| $$$\ /  $$$| $$   | $$| $$__| $$| $$  | $$| $$      | $$   
| $$$$\  $$$$ \$$\ /  $$| $$    $$| $$  | $$| $$      | $$   
| $$\$$ $$ $$  \$$\  $$ | $$$$$$$$| $$  | $$| $$      | $$   
| $$ \$$$| $$   \$$ $$  | $$  | $$| $$__/ $$| $$_____ | $$   
| $$  \$ | $$    \$$$   | $$  | $$ \$$    $$| $$     \| $$   
 \$$      \$$     \$     \$$   \$$  \$$$$$$  \$$$$$$$$ \$$   

` + "\t" + `designed in 2020 by awesome people with luv

Use: mvault [-opts]

Common:
encode file` + "\t" + `mvault -e -file=file_to_encode -u username
decode file` + "\t" + `mvault -d -file=file_to_decode -u username

Options:
-d, -decrypt` + "\t" + `Decode given file
-e, -encrypt` + "\t" + `Encode given file
-file[=path]` + "\t" + `Recieve file path
-h, -help` + "\t" + `Show docs
-l, -local` + "\t" + `Do not interact with server
-u, -user` + "\t" + `Define username`

/*GetHelp returns help note*/
func GetHelp() string {
	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	result := fmt.Sprintln(w, helpNote)
	w.Flush()
	// FIXME: result stores 68 characters of crap and payload. WHY?
	return result[68:]
}
