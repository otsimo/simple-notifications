package notification
import (
	"log"

	"golang.org/x/net/context"
	"fmt"
)

type Server struct {}

func (s *Server) SendMessage(ctx context.Context, in*Message) (*SendMessageResponse, error) {
	log.Println("SendMessage", in.Template)
	fmt.Printf("%d of targets found\n", len(in.Targets))

	for _, t := range (in.Targets) {
		if t.Email != nil {
			fmt.Println("target: Email", t.Email)
		}else if t.Push != nil {
			fmt.Println("target: Push", t.Push)
		}else if t.Sms != nil {
			fmt.Println("target: Sms", t.Sms)
		}
	}
	return &SendMessageResponse{}, nil
}
