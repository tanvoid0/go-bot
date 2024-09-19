package scratch_pad

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/tanvoid0/dev-bot/util"
	"log"
)

func DemoPrint() {
	markdown := util.GetMarkdownTemplate()
	//html := markdown.ToHTML([]byte(util.GetMarkdownTemplate()), nil, nil)
	//fmt.Printf("\033[1m%s\033[0m\n", string(html)).
	//res := blackfriday.Run([]byte(markdown))
	//fmt.Println(string(res))

	rendered, err := glamour.Render(markdown, "dark")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(rendered)

}
