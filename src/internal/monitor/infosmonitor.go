package monitor

type InfosMonitor struct {
	Id       int    `json:"id"`
	Type     int    `json:"type"`
	Location int    `json:"location"`
	Data     string `json:"data"`
}

func NewInfosMonitor(columns []string, row []interface{}) (*InfosMonitor, error) {
	rowData := &InfosMonitor{}
	for i, col := range columns {
		switch col {
		case "id":
			rowData.Id = int(row[i].(float64))
		case "type":
			rowData.Type = int(row[i].(float64))
		case "location":
			rowData.Location = int(row[i].(float64))
		case "data":
			rowData.Data = row[i].(string)
		}
	}
	return rowData, nil
}
