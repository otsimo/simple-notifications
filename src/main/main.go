package main
import "fmt"
import (
	"text/template"
	"bytes"
)

func main() {
	fmt.Println("Hello World!!")

	const letter =
`Dear {{.name}}, And {{.gift}},
{{.attended}}
Best wishes,
SD
`
	t := template.Must(template.New("letter").Parse(letter))

	commits := map[string]interface{}{
		"name": "Sercan Degirmenci",
		"attended":   true,
		"gift": "xyaz",
	}

	var doc bytes.Buffer
	err := t.Execute(&doc, commits)
	if err == nil {
		s := doc.String()
		fmt.Println(s)
	}else {
		panic(err)
	}
}