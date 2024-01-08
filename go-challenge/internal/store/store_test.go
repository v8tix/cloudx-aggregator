package store

import (
	"encoding/json"
	"fmt"
	"github.com/cloudx-labs/challenge/internal/model/dto"
	"github.com/cloudx-labs/challenge/internal/utils"
	"github.com/google/uuid"
	"log"
	"testing"
)

const (
	groupMocks = "../../mock/group/"
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
	associationsStore = func() *AssociationsStore {
		return NewAssociationsStore()
	}
	responseStore = func() *ResponseStore {
		return NewResponseStore()
	}()
	parentGroupArr = func() []ParentGroup {
		return []ParentGroup{
			{
				Source:      uuid.NewString(),
				Destination: uuid.NewString(),
			},
			{
				Source:      uuid.NewString(),
				Destination: uuid.NewString(),
			},
		}
	}()
)

func TestSearchingAssociationStoreByParentUUIDWithSourceAndDestinationUUIDsReturnParentUUID(t *testing.T) {
	t.Parallel()

	ag := associationsStore()

	for _, groupDTO := range *group1 {
		ag.Add(&groupDTO)
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
			_, result := ag.FindParentsByChildren(&tc.group)
			if result != tc.want {
				t.Fatalf("Expected: %v, Got: %v", true, result)
			}
		})
	}
}

func TestSearchingAssociationStoreByParentUUIDWithSourceAndDestinationUUIDsWithEmptyStoreReturnFalse(t *testing.T) {
	t.Parallel()

	ag := associationsStore()

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
			_, result := ag.FindParentsByChildren(&tc.group)
			if result != tc.want {
				t.Fatalf("Expected: %v, Got: %v", true, result)
			}
		})
	}
}

func TestGettingResponseStoreResponseWithNonEmptyStoreReturnsResponse(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input []ParentGroup
		want  bool
	}{
		"with non empty store": {
			input: parentGroupArr,
			want:  false,
		},
	}

	for input, tc := range cases {
		t.Run(input, func(t *testing.T) {
			for _, parentGroup := range tc.input {
				responseStore.Add(parentGroup)
			}

			bytes, err := responseStore.ToResponse()
			if err != nil {
				t.Fatalf("Error reading store")
			}

			if bytes == nil {
				t.Fatalf("Expected non nil response")
			}
		})
	}
}
