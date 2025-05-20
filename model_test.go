package sqlorm_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm/v2"
)

func Test_Model(t *testing.T) {
	now := time.Now()
	model := &sqlorm.ModelSnakeCase{
		ID:        uuid.New(),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	abc, err := json.Marshal(model)
	require.Nil(t, err)
	fmt.Println(string(abc))
}
