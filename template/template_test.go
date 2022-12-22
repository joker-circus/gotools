package template

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type _Test struct {
	Status    string    `json:"status" gorm:"column:status"`
	Business  string    `json:"business" gorm:"column:business"`
	Service   string    `json:"service" gorm:"column:service"`
	Cluster   string    `json:"cluster" gorm:"column:cluster"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:deadline"`
}

func (c _Test) StatusCN(status string) string {
	return "中文 " + status
}

var _date = _Test{
	Status:    "new",
	Business:  "测试业务",
	Service:   "test-service",
	Cluster:   "测试集群",
	UpdatedAt: time.Now(),
}

func TestRender(t *testing.T) {
	templateContent := `
Business：{{.Business}}
Service：{{Title .Service}}
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

	result := fmt.Sprintf(`
Business：%s
Service：%s
Cluster：%s
StatusCN：%s
UpdatedAt：%s`, _date.Business, strings.Title(_date.Service), _date.Cluster,
		_date.StatusCN(_date.Status), _date.UpdatedAt.Format("2006-01-02 15:04:05"))
	assert.Equal(t, result, content, "the should be equal")
}
