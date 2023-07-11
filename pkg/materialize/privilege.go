package materialize

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

var Permissions = map[string]string{
	"r": "SELECT",
	"a": "INSERT",
	"w": "UPDATE",
	"d": "DELETE",
	"C": "CREATE",
	"U": "USAGE",
	"R": "CREATEROLE",
	"B": "CREATEDB",
	"N": "CREATECLUSTER",
}

type ObjectType struct {
	Permissions []string
}

// https://materialize.com/docs/sql/grant-privilege/#details
var ObjectPermissions = map[string]ObjectType{
	"DATABASE": {
		Permissions: []string{"U", "C"},
	},
	"SCHEMA": {
		Permissions: []string{"U", "C"},
	},
	"TABLE": {
		Permissions: []string{"a", "r", "w", "d"},
	},
	"VIEW": {
		Permissions: []string{"r"},
	},
	"MATERIALIZED VIEW": {
		Permissions: []string{"r"},
	},
	"INDEX": {
		Permissions: []string{},
	},
	"TYPE": {
		Permissions: []string{"U"},
	},
	"SOURCE": {
		Permissions: []string{"r"},
	},
	"SINK": {
		Permissions: []string{},
	},
	"CONNECTION": {
		Permissions: []string{"U"},
	},
	"SECRET": {
		Permissions: []string{"U"},
	},
	"CLUSTER": {
		Permissions: []string{"U", "C"},
	},
	"SYSTEM": {
		Permissions: []string{"R", "B", "N"},
	},
}

func ParsePrivileges(privileges string) map[string][]string {
	o := map[string][]string{}

	privileges = strings.TrimPrefix(privileges, "{")
	privileges = strings.TrimSuffix(privileges, "}")

	for _, p := range strings.Split(privileges, ",") {
		e := strings.Split(p, "=")

		roleId := e[0]
		roleprivileges := strings.Split(e[1], "/")[0]

		privilegeMap := []string{}
		for _, rp := range strings.Split(roleprivileges, "") {
			v := Permissions[rp]
			privilegeMap = append(privilegeMap, v)
		}

		o[roleId] = privilegeMap
	}

	return o
}

func HasPrivilege(privileges []string, checkPrivilege string) bool {
	for _, v := range privileges {
		if v == checkPrivilege {
			return true
		}
	}
	return false
}

type PrivilegeObjectStruct struct {
	Type         string
	Name         string
	SchemaName   string
	DatabaseName string
}

func GetPrivilegeObjectStruct(databaseName string, schemaName string, v interface{}) PrivilegeObjectStruct {
	var p PrivilegeObjectStruct
	u := v.([]interface{})[0].(map[string]interface{})

	if v, ok := u["type"]; ok {
		p.Type = v.(string)
	}

	if v, ok := u["name"]; ok {
		p.Name = v.(string)
	}

	if v, ok := u["schema_name"]; ok && v.(string) != "" {
		p.SchemaName = v.(string)
	}

	if v, ok := u["database_name"]; ok && v.(string) != "" {
		p.DatabaseName = v.(string)
	}

	return p
}

func (i *PrivilegeObjectStruct) QualifiedName() string {
	p := []string{}

	if i.DatabaseName != "" {
		p = append(p, i.DatabaseName)
	}

	if i.SchemaName != "" {
		p = append(p, i.SchemaName)
	}

	p = append(p, i.Name)
	return QualifiedName(p...)
}

// DDL
type PrivilegeBuilder struct {
	ddl       Builder
	role      string
	privilege string
	object    PrivilegeObjectStruct
}

func NewPrivilegeBuilder(conn *sqlx.DB, role, privilege string, object PrivilegeObjectStruct) *PrivilegeBuilder {
	return &PrivilegeBuilder{
		ddl:       Builder{conn, Privilege},
		role:      role,
		privilege: privilege,
		object:    object,
	}
}

// https://materialize.com/docs/sql/grant-privilege/#compatibility
func objectCompatibility(objectType string) string {
	compatibility := []string{"SOURCE", "VIEW", "MATERIALIZED VIEW"}

	for _, c := range compatibility {
		if c == objectType {
			return "TABLE"
		}
	}
	return objectType
}

func (b *PrivilegeBuilder) Grant() error {
	t := objectCompatibility(b.object.Type)
	q := fmt.Sprintf(`GRANT %s ON %s %s TO %s;`, b.privilege, t, b.object.QualifiedName(), b.role)
	return b.ddl.exec(q)
}

func (b *PrivilegeBuilder) Revoke() error {
	t := objectCompatibility(b.object.Type)
	q := fmt.Sprintf(`REVOKE %s ON %s %s FROM %s;`, b.privilege, t, b.object.QualifiedName(), b.role)
	return b.ddl.exec(q)
}

func PrivilegeId(conn *sqlx.DB, object PrivilegeObjectStruct, roleId, privilege string) (string, error) {
	var id string

	switch t := object.Type; t {
	case "DATABASE":
		o := ObjectSchemaStruct{Name: object.Name}
		i, err := DatabaseId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "SCHEMA":
		o := ObjectSchemaStruct{Name: object.Name, DatabaseName: object.DatabaseName}
		i, err := SchemaId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "TABLE":
		o := ObjectSchemaStruct{Name: object.Name, SchemaName: object.SchemaName, DatabaseName: object.DatabaseName}
		i, err := TableId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "VIEW":
		o := ObjectSchemaStruct{Name: object.Name, SchemaName: object.SchemaName, DatabaseName: object.DatabaseName}
		i, err := ViewId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "MATERIALIZED VIEW":
		o := ObjectSchemaStruct{Name: object.Name, SchemaName: object.SchemaName, DatabaseName: object.DatabaseName}
		i, err := MaterializedViewId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "TYPE":
		o := ObjectSchemaStruct{Name: object.Name, SchemaName: object.SchemaName, DatabaseName: object.DatabaseName}
		i, err := TypeId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "SOURCE":
		o := ObjectSchemaStruct{Name: object.Name, SchemaName: object.SchemaName, DatabaseName: object.DatabaseName}
		i, err := SourceId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "CONNECTION":
		o := ObjectSchemaStruct{Name: object.Name, SchemaName: object.SchemaName, DatabaseName: object.DatabaseName}
		i, err := ConnectionId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "SECRET":
		o := ObjectSchemaStruct{Name: object.Name, SchemaName: object.SchemaName, DatabaseName: object.DatabaseName}
		i, err := SecretId(conn, o)
		if err != nil {
			return "", err
		}
		id = i

	case "CLUSTER":
		o := ObjectSchemaStruct{Name: object.Name}
		i, err := ClusterId(conn, o)
		if err != nil {
			return "", err
		}
		id = i
	}

	f := fmt.Sprintf(`GRANT|%s|%s|%s|%s`, object.Type, id, roleId, privilege)
	return f, nil
}

func ScanPrivileges(conn *sqlx.DB, objectType, objectId string) (string, error) {
	var params string

	switch t := objectType; t {
	case "DATABASE":
		p, err := ScanDatabase(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "SCHEMA":
		p, err := ScanSchema(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "TABLE":
		p, err := ScanTable(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "VIEW":
		p, err := ScanView(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "MATERIALIZED VIEW":
		p, err := ScanMaterializedView(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "TYPE":
		p, err := ScanType(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "SOURCE":
		p, err := ScanSource(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "CONNECTION":
		p, err := ScanConnection(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "SECRET":
		p, err := ScanSecret(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String

	case "CLUSTER":
		p, err := ScanCluster(conn, objectId)
		if err != nil {
			return "", err
		}
		params = p.Privileges.String
	}

	return params, nil
}
