package sys

import(
	"fmt"

	"iris-ticket/backend/app/models/basemodel"
)

func TableName(name string) string {
	return fmt.Sprintf("%s%s%s", basemodel.GetTablePrefix(),"sys_", name)
}