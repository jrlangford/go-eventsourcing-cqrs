package server

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jrlangford/go-eventsourcing-cqrs/services/ui/internal/server/pbcommand"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/ui/internal/server/pbquery"
	"google.golang.org/grpc"
)

// A server manages the server's resources.
type server struct {
	commandClient pbcommand.CommandClient
	queryClient   pbquery.QueryClient
	template      *template.Template
}

// New constructs a new server.
func New() *server {
	return &server{}
}

// Run initilizes and runs the server.
func (srv *server) Run() {

	srv.template = template.Must(template.ParseGlob("web/templates/*"))

	commConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Println("Could not connect to command server:", err)
	}

	srv.commandClient = pbcommand.NewCommandClient(commConn)

	queryConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Println("Could not connect to query server:", err)
	}

	srv.queryClient = pbquery.NewQueryClient(queryConn)

	mux := srv.setupHandlers()
	if err := http.ListenAndServe("localhost:8088", mux); err != nil {
		log.Println(err)
	}

}

// setupHandlers configures defines handlers and thir routes.
func (srv *server) setupHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		resp, err := srv.queryClient.Run(r.Context(), &pbquery.QueryMessage{
			Query: &pbquery.QueryMessage_GetInventoryItems{
				GetInventoryItems: &pbquery.GetInventoryItemsRequest{},
			},
		})
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/unavailable", http.StatusSeeOther)
			return
		}

		assertedResponse, ok := resp.Response.(*pbquery.QueryResponse_GetInventoryItemsResponse)
		if !ok {
			log.Println("Response is not of the expected type: %+v", resp)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		items := assertedResponse.GetInventoryItemsResponse.ItemList

		sort.Slice(items, func(i, j int) bool {
			return items[i].Name < items[j].Name
		})

		err = srv.template.ExecuteTemplate(w, "index", items)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//id := uuid.New().String()
			_, err = srv.commandClient.Execute(r.Context(), &pbcommand.CommandMessage{
				Command: &pbcommand.CommandMessage_CreateInventoryItem{
					CreateInventoryItem: &pbcommand.CreateInventoryItem{
						Name: r.Form.Get("name"),
					},
				},
			})
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/unavailable", http.StatusSeeOther)
				return
			}

			redirectURL := "/"
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}

		err := srv.template.ExecuteTemplate(w, "add", nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/details/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.Split(r.URL.Path, "/")
		id := p[len(p)-1]

		resp, err := srv.queryClient.Run(r.Context(), &pbquery.QueryMessage{
			Query: &pbquery.QueryMessage_GetInventoryItemDetails{
				GetInventoryItemDetails: &pbquery.GetInventoryItemDetailsRequest{
					Uuid: id,
				},
			},
		})
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/unavailable", http.StatusSeeOther)
			return
		}

		assertedResponse, ok := resp.Response.(*pbquery.QueryResponse_GetInventoryItemDetailsResponse)
		if !ok {
			log.Println("Response is not of the expected type: %+v", resp)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		itemDetails := assertedResponse.GetInventoryItemDetailsResponse.ItemDetails

		err = srv.template.ExecuteTemplate(w, "details", itemDetails)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/deactivate/", func(w http.ResponseWriter, r *http.Request) {

		p := strings.Split(r.URL.Path, "/")
		id := p[len(p)-1]

		tdata := map[string]interface{}{
			"ID": id,
		}

		if r.Method == http.MethodPost {

			_, err := srv.commandClient.Execute(r.Context(), &pbcommand.CommandMessage{
				Command: &pbcommand.CommandMessage_DeactivateInventoryItem{
					DeactivateInventoryItem: &pbcommand.DeactivateInventoryItem{
						Uuid: id,
					},
				},
			})
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/unavailable", http.StatusSeeOther)
				return
			}

			redirectURL := "/"
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}

		err := srv.template.ExecuteTemplate(w, "deactivate", tdata)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/changename/", func(w http.ResponseWriter, r *http.Request) {

		p := strings.Split(r.URL.Path, "/")
		id := p[len(p)-1]

		tdata := map[string]interface{}{
			"ID": id,
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = srv.commandClient.Execute(r.Context(), &pbcommand.CommandMessage{
				Command: &pbcommand.CommandMessage_RenameInventoryItem{
					RenameInventoryItem: &pbcommand.RenameInventoryItem{
						Uuid:    id,
						NewName: r.Form.Get("name"),
					},
				},
			})
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/unavailable", http.StatusSeeOther)
				return
			}

			redirectURL := "/"
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}

		err := srv.template.ExecuteTemplate(w, "changename", tdata)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	mux.HandleFunc("/checkin/", func(w http.ResponseWriter, r *http.Request) {

		p := strings.Split(r.URL.Path, "/")
		id := p[len(p)-1]

		tdata := map[string]interface{}{
			"ID": id,
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			num, err := strconv.ParseInt(r.Form.Get("number"), 10, 64)
			if err != nil {
				http.Error(w, "Unable to read number.", http.StatusInternalServerError)
				return
			}

			_, err = srv.commandClient.Execute(r.Context(), &pbcommand.CommandMessage{
				Command: &pbcommand.CommandMessage_CheckInItemsToInventory{
					CheckInItemsToInventory: &pbcommand.CheckInItemsToInventory{
						Uuid:  id,
						Count: num,
					},
				},
			})
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/unavailable", http.StatusSeeOther)
				return
			}

			redirectURL := "/"
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}

		err := srv.template.ExecuteTemplate(w, "checkin", tdata)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	mux.HandleFunc("/remove/", func(w http.ResponseWriter, r *http.Request) {

		p := strings.Split(r.URL.Path, "/")
		id := p[len(p)-1]

		tdata := map[string]interface{}{
			"ID": id,
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			num, err := strconv.ParseInt(r.Form.Get("number"), 10, 64)
			if err != nil {
				http.Error(w, "Unable to read number.", http.StatusInternalServerError)
				return
			}

			_, err = srv.commandClient.Execute(r.Context(), &pbcommand.CommandMessage{
				Command: &pbcommand.CommandMessage_RemoveItemsFromInventory{
					RemoveItemsFromInventory: &pbcommand.RemoveItemsFromInventory{
						Uuid:  id,
						Count: num,
					},
				},
			})
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/unavailable", http.StatusSeeOther)
				return
			}

			redirectURL := "/"
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}

		err := srv.template.ExecuteTemplate(w, "remove", tdata)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		staticFile := r.URL.Path[len("/assets/"):]
		if len(staticFile) != 0 {
			f, err := http.Dir("web/assets/").Open(staticFile)
			if err == nil {
				content := io.ReadSeeker(f)
				http.ServeContent(w, r, staticFile, time.Now(), content)
				return
			}
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/unavailable", func(w http.ResponseWriter, r *http.Request) {
		err := srv.template.ExecuteTemplate(w, "unavailable", nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return mux
}
