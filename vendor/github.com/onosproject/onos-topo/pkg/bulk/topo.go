// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bulk

import (
	"fmt"

	configlib "github.com/onosproject/onos-lib-go/pkg/config"
	"github.com/onosproject/onos-topo/api/topo"
)

var topoConfig *TopoConfig

// TopoConfig - the top level object
type TopoConfig struct {
	TopoKinds     []TopoKind
	TopoEntities  []TopoEntity
	TopoRelations []TopoRelation
}

// TopoKind - required to get around the "oneof" Obj
type TopoKind struct {
	ID         topo.ID
	Type       topo.Object_Type
	Obj        *topo.Object_Kind
	Attributes *map[string]string
}

// TopoKindToTopoObject - convert to Object
func TopoKindToTopoObject(topoKind *TopoKind) *topo.Object {
	return &topo.Object{
		ID:         topoKind.ID,
		Type:       topoKind.Type,
		Obj:        topoKind.Obj,
		Attributes: *topoKind.Attributes,
	}
}

// TopoEntity - required to get around the "oneof" Obj
type TopoEntity struct {
	ID         topo.ID
	Type       topo.Object_Type
	Obj        *topo.Object_Entity
	Attributes *map[string]string
}

// TopoEntityToTopoObject - convert to Object
func TopoEntityToTopoObject(topoEntity *TopoEntity) *topo.Object {
	return &topo.Object{
		ID:         topoEntity.ID,
		Type:       topoEntity.Type,
		Obj:        topoEntity.Obj,
		Attributes: *topoEntity.Attributes,
	}
}

// TopoRelation - required to get around the "oneof" Obj
type TopoRelation struct {
	ID         topo.ID
	Type       topo.Object_Type
	Obj        *topo.Object_Relation
	Attributes *map[string]string
}

// TopoRelationToTopoObject - convert to Object
func TopoRelationToTopoObject(topoRelation *TopoRelation) *topo.Object {
	return &topo.Object{
		ID:         topoRelation.ID,
		Type:       topoRelation.Type,
		Obj:        topoRelation.Obj,
		Attributes: *topoRelation.Attributes,
	}
}

// ClearTopo - reset the config - needed for tests
func ClearTopo() {
	topoConfig = nil
}

// GetTopoConfig gets the onos-topo configuration
func GetTopoConfig(location string) (TopoConfig, error) {
	if topoConfig == nil {
		topoConfig = &TopoConfig{}
		if err := configlib.LoadNamedConfig(location, topoConfig); err != nil {
			return TopoConfig{}, err
		}
		if err := TopoChecker(topoConfig); err != nil {
			return TopoConfig{}, err
		}
	}
	return *topoConfig, nil
}

// TopoChecker - check everything is within bounds
func TopoChecker(config *TopoConfig) error {
	if len(config.TopoKinds) == 0 {
		return fmt.Errorf("no kinds found")
	}

	for _, kind := range config.TopoKinds {
		topoKind := kind // pin
		if topoKind.Type != topo.Object_KIND {
			return fmt.Errorf("unexpected type %v for TopoKind", topoKind.Type)
		} else if topoKind.ID == topo.NullID {
			return fmt.Errorf("empty ref for TopoKind")
		} else if topoKind.Obj.Kind.GetName() == "" {
			return fmt.Errorf("empty name for TopoKind")
		}
	}

	if len(config.TopoEntities) == 0 {
		return fmt.Errorf("no entities found")
	}

	for _, entity := range config.TopoEntities {
		topoEntity := entity // pin
		if topoEntity.Type != topo.Object_ENTITY {
			return fmt.Errorf("unexpected type %v for TopoEntity", topoEntity.Type)
		} else if topoEntity.ID == topo.NullID {
			return fmt.Errorf("empty ref for TopoEntity")
		}
	}

	for _, relation := range config.TopoRelations {
		topoRelation := relation // pin
		if topoRelation.Type != topo.Object_RELATION {
			return fmt.Errorf("unexpected type %v for TopoRelation", topoRelation.Type)
		} else if topoRelation.ID == topo.NullID {
			return fmt.Errorf("null id for TopoRelation")
		} else if topoRelation.Obj.Relation.SrcEntityID == "" {
			return fmt.Errorf("null source entity id for TopoRelation")
		} else if topoRelation.Obj.Relation.TgtEntityID == "" {
			return fmt.Errorf("null target entity id for TopoRelation")
		}
	}

	return nil
}
