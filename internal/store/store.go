package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type AuditEntry struct {
	ID string `json:"id"`
	Action string `json:"name"`
	Actor string `json:"actor"`
	Resource string `json:"resource"`
	Details string `json:"details"`
	IPAddress string `json:"ip_address"`
	Source string `json:"source"`
	Severity string `json:"severity"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"brand.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS audit_entries(id TEXT PRIMARY KEY,name TEXT NOT NULL,actor TEXT DEFAULT '',resource TEXT DEFAULT '',details TEXT DEFAULT '',ip_address TEXT DEFAULT '',source TEXT DEFAULT '',severity TEXT DEFAULT 'info',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *AuditEntry)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO audit_entries(id,name,actor,resource,details,ip_address,source,severity,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Action,e.Actor,e.Resource,e.Details,e.IPAddress,e.Source,e.Severity,e.CreatedAt);return err}
func(d *DB)Get(id string)*AuditEntry{var e AuditEntry;if d.db.QueryRow(`SELECT id,name,actor,resource,details,ip_address,source,severity,created_at FROM audit_entries WHERE id=?`,id).Scan(&e.ID,&e.Action,&e.Actor,&e.Resource,&e.Details,&e.IPAddress,&e.Source,&e.Severity,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]AuditEntry{rows,_:=d.db.Query(`SELECT id,name,actor,resource,details,ip_address,source,severity,created_at FROM audit_entries ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []AuditEntry;for rows.Next(){var e AuditEntry;rows.Scan(&e.ID,&e.Action,&e.Actor,&e.Resource,&e.Details,&e.IPAddress,&e.Source,&e.Severity,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *AuditEntry)error{_,err:=d.db.Exec(`UPDATE audit_entries SET name=?,actor=?,resource=?,details=?,ip_address=?,source=?,severity=? WHERE id=?`,e.Action,e.Actor,e.Resource,e.Details,e.IPAddress,e.Source,e.Severity,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM audit_entries WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM audit_entries`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]AuditEntry{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["source"];ok&&v!=""{where+=" AND source=?";args=append(args,v)}
    if v,ok:=filters["severity"];ok&&v!=""{where+=" AND severity=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,actor,resource,details,ip_address,source,severity,created_at FROM audit_entries WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []AuditEntry;for rows.Next(){var e AuditEntry;rows.Scan(&e.ID,&e.Action,&e.Actor,&e.Resource,&e.Details,&e.IPAddress,&e.Source,&e.Severity,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}
