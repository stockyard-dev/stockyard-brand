package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-brand/internal/server";"github.com/stockyard-dev/stockyard-brand/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./brand-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("brand: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Brand — Self-hosted audit trail and compliance log\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n",port,port)
log.Printf("brand: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
