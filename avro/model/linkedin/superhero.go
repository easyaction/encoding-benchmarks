package model

// Superhero is a model with custom hydration and extraction methods required to work with linkedin/goavro maps.
type Superhero struct {
	ID            int32
	AffiliationID int32
	Name          string
	Life          float32
	Energy        float32
	Powers        []*Superpower
}

// ToMap extracts model data to a map required to perform encoding.
func (s *Superhero) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"id":             s.ID,
		"affiliation_id": s.AffiliationID,
		"name":           s.Name,
		"life":           s.Life,
		"energy":         s.Energy,
		"powers":         make([]interface{}, len(s.Powers)),
	}

	for i, p := range s.Powers {
		data["powers"].([]interface{})[i] = p.ToMap()
	}

	return data
}

// FromMap hydrates model with data from a provided map.
func (s *Superhero) FromMap(data map[string]interface{}) {
	s.ID = mapGet(data, []string{"id"}, 0).(int32)
	s.AffiliationID = mapGet(data, []string{"affiliation_id"}, 0).(int32)
	s.Name = mapGet(data, []string{"name"}, "").(string)
	s.Life = mapGet(data, []string{"life"}, 0).(float32)
	s.Energy = mapGet(data, []string{"energy"}, 0).(float32)

	powers := mapGet(data, []string{"powers"}, []interface{}{}).([]interface{})
	s.Powers = make([]*Superpower, len(powers))
	for i, p := range powers {
		power := &Superpower{}
		power.FromMap(p.(map[string]interface{}))
		s.Powers[i] = power
	}
}

// Superpower is a model with custom hydration and extraction methods required to work with linkedin/goavro maps.
type Superpower struct {
	ID      int32
	Name    string
	Damage  float32
	Energy  float32
	Passive bool
}

// ToMap extracts model data to a map required to perform encoding.
func (s *Superpower) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      s.ID,
		"name":    s.Name,
		"damage":  s.Damage,
		"energy":  s.Energy,
		"passive": s.Passive,
	}
}

// FromMap hydrates model with data from a provided map.
func (s *Superpower) FromMap(data map[string]interface{}) {
	s.ID = mapGet(data, []string{"id"}, 0).(int32)
	s.Name = mapGet(data, []string{"name"}, "").(string)
	s.Damage = mapGet(data, []string{"damage"}, 0).(float32)
	s.Energy = mapGet(data, []string{"energy"}, 0).(float32)
	s.Passive = mapGet(data, []string{"passive"}, false).(bool)
}

// mapGet is a helper function allowing safe access of an undefined map that is a result of decoding by linkedin/goavro.
func mapGet(m map[string]interface{}, path []string, empty interface{}) interface{} {
	var temp = m
	var v interface{}
	var ok bool

	for i, k := range path {
		v, ok = temp[k]
		if !ok {
			return empty
		}

		if i == len(path)-1 {
			break
		}

		temp = v.(map[string]interface{})
	}

	return v
}
