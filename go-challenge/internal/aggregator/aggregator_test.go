package aggregator

import (
	"encoding/json"
	"fmt"
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/utils"
	"log"
	"testing"
)

const (
	groupMocks = "../mock/group/"
)

var (
	group1 = func() *[]dto.GroupDTO {
		var groupDTO []dto.GroupDTO

		data, err := utils.ReadFile(fmt.Sprintf("%s%s", groupMocks, "group_1.json"))
		if err != nil {
			log.Fatal(err)
		}
		_ = json.Unmarshal(data, &groupDTO)
		return &groupDTO
	}()
	associationsAggregator = func() *AssociationsAggregator {
		return NewAssociationsAggregator()
	}
)

func TestSearchingParentUUIDWithSourceAndDestinationUUIDsReturnParentUUID(t *testing.T) {
	t.Parallel()

	ag := associationsAggregator()

	for _, groupDTO := range *group1 {
		ag.AddAssociations(&groupDTO)
	}

	cases := map[string]struct {
		group dto.GroupDTO
		want  bool
	}{
		"with existing children source and destination id": {
			group: (*group1)[0],
			want:  true,
		},
	}

	for input, tc := range cases {
		t.Run(input, func(t *testing.T) {
			result := ag.FindParentByChildren(&tc.group)
			if result != tc.want {
				t.Fatalf("Expected: %v, Got: %v", true, result)
			}
		})
	}
}

func TestSearchingParentUUIDWithSourceAndDestinationUUIDsWithEmptyAggregatorReturnFalse(t *testing.T) {
	t.Parallel()

	ag := associationsAggregator()

	cases := map[string]struct {
		group dto.GroupDTO
		want  bool
	}{
		"with existing children source and destination id": {
			group: (*group1)[0],
			want:  false,
		},
	}

	for input, tc := range cases {
		t.Run(input, func(t *testing.T) {
			result := ag.FindParentByChildren(&tc.group)
			if result != tc.want {
				t.Fatalf("Expected: %v, Got: %v", true, result)
			}
		})
	}
}
