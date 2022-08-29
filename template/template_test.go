package template

import (
	"fmt"
	"testing"
	"time"
)

type _Test struct {
	Status         string    `json:"status" gorm:"column:status"`
	Business       string    `json:"business" gorm:"column:business"`
	Service        string    `json:"service" gorm:"column:service"`
	Cluster        string    `json:"cluster" gorm:"column:cluster"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:deadline"`
}

func (c _Test) StatusCN(status string) string {
	return "中文 " + status
}

var _date = _Test{
	Status:         "new",
	Business:       "测试业务",
	Service:        "test-service",
	Cluster:        "测试集群",
	UpdatedAt:       time.Now(),

}

func TestRender(t *testing.T) {
	templateContent := `
Business：{{.Business}}
Service：{{.Service}}
{{ if .Cluster }}
Cluster：{{.Cluster}}
{{ end }}
StatusCN：{{.StatusCN .Status}}
UpdatedAt：{{.UpdatedAt.Format "2006-01-02 15:04:05"}}`

	content, err := Render(templateContent, _date)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(content)
}
