package sender

import "livebets/analazer/internal/entity"

type GeneralFilter struct {
	Bookmakers []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	Name     string `json:"name"`
	Live     Live   `json:"live"`
	Prematch Live   `json:"prematch"`
}

type Live struct {
	Filter bool     `json:"filter"`
	Sports []string `json:"sports"`
}

func (s *Sender) Filter(data []entity.ResponsePair, filter GeneralFilter) []entity.ResponsePair {
	var newData []entity.ResponsePair
	for _, val := range data {

		// Bookmakers
		flag := false
		for _, bookmaker := range filter.Bookmakers {
			if bookmaker.Name == val.Second.Bookmaker {
				
				if !bookmaker.Live.Filter && val.IsLive {
					continue
				}
		
				if !bookmaker.Prematch.Filter && !val.IsLive {
					continue
				}
		
				// Live
				if bookmaker.Live.Filter && val.IsLive {
					flag = false
					for _, sport := range bookmaker.Live.Sports {
						if sport == val.SportName {
							flag = true
							break
						}
					}
					if !flag {
						continue
					}
				}
		
				// Prematch
				if bookmaker.Prematch.Filter && !val.IsLive {
					flag = false
					for _, sport := range bookmaker.Prematch.Sports {
						if sport == val.SportName {
							flag = true
							break
						}
					}
					if !flag {
						continue
					}
				}
				
				flag = true
				break
			}
		}
		if !flag {
			continue
		}

		newData = append(newData, val)
	}

	return newData
}
