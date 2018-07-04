package dnspod

import (
	"fmt"
	"strconv"
)

type Record struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Line          string `json:"line,omitempty"`
	LineID        string `json:"line_id,omitempty"`
	Type          string `json:"type,omitempty"`
	TTL           string `json:"ttl,omitempty"`
	Value         string `json:"value,omitempty"`
	MX            string `json:"mx,omitempty"`
	Enabled       string `json:"enabled,omitempty"`
	Status        string `json:"status,omitempty"`
	MonitorStatus string `json:"monitor_status,omitempty"`
	Remark        string `json:"remark,omitempty"`
	UpdateOn      string `json:"updated_on,omitempty"`
	UseAQB        string `json:"use_aqb,omitempty"`
	SubDomain string `json:"sub_domain,omitempty"`
	RecordType string `json:"record_type,omitempty"`
	RecordLine string `json:"record_line,omitempty"`
	RecordLineID string `json:"record_line_id,omitempty"`
}

type RecordsInfo struct {
	SubDomains int `json:"sub_domains"`
	RecordTotal int `json:"record_total"`
}

func (r Record) String() string {
	bs, _ := json.Marshal(r)
	return string(bs)
}

type RecordLine struct {
  Line string `json:"line"`
  LineID string `json:"line_id"`
}

type RecordQuery struct {
	DomainID string	 `json:"domainID"`
	Domain string	 `json:"domain,omitempty"`
	CurrentPage int	 `json:"currentPage,omitempty"`
	PageSize int	 `json:"pageSize,omitempty"`
	SubDomain string `json:"subDomain,omitempty"`
	Keyword string	 `json:"keyword,omitempty"`
}

type PaginationRecordList struct {
	CurrentPage int `json:"currentPage"`
	PageSize int `json:"pageSize"`
	Total int `json:"total"`
	List []Record `json:"list"`
}

type linesWrapper struct {
	Lines []string `json:"lines"`
	LineIDs map[string]interface{} `json:"line_ids"`
	Status Status `json:"status"`
}

type recordsWrapper struct {
	Status  Status     `json:"status"`
	Info    RecordsInfo `json:"info"`
	Records []Record   `json:"records"`
}

type recordWrapper struct {
	Status Status     `json:"status"`
	Info   DomainInfo `json:"info"`
	Record Record     `json:"record"`
}

// recordAction generates the resource path for given record that belongs to a domain.
func recordAction(action string) string {
	if len(action) > 0 {
		return fmt.Sprintf("Record.%s", action)
	}
	return "Record.List"
}

// List the domain records.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-list

func (s *DomainsService) ListRecords(query RecordQuery) (PaginationRecordList, *Response, error) {
	path := recordAction("List")

	payload := newPayLoad(s.client.CommonParams)

	if query.DomainID != "" {
		payload.Add("domain_id", query.DomainID)
	}
	if query.Domain != "" {
		payload.Add("domain", query.Domain)
	}
	if query.PageSize != 0 {
		payload.Add("offset", strconv.Itoa(query.CurrentPage))
		payload.Add("length", strconv.Itoa(query.PageSize))
	}

	if query.SubDomain != "" {
		payload.Add("sub_domain", query.SubDomain)
	}
	if query.Keyword != "" {
		payload.Add("keyword", query.Keyword)
	}

	wrappedRecords := recordsWrapper{}

	res, err := s.client.post(path, payload, &wrappedRecords)
	if err != nil {
		return PaginationRecordList{}, res, err
	}

	if wrappedRecords.Status.Code != "1" {
		return PaginationRecordList{}, nil, fmt.Errorf("Could not get domains: %s", wrappedRecords.Status.Message)
	}

	return PaginationRecordList{
		CurrentPage: query.CurrentPage,
		PageSize: query.PageSize,
		Total: wrappedRecords.Info.RecordTotal,
		List: wrappedRecords.Records,
	}, res, nil
}

// CreateRecord creates a domain record.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-create
func (s *DomainsService) CreateRecord(domain string, recordAttributes Record) (Record, *Response, error) {
	path := recordAction("Create")

	payload := newPayLoad(s.client.CommonParams)

	payload.Add("domain_id", domain)

	if recordAttributes.Name != "" {
		payload.Add("sub_domain", recordAttributes.Name)
	}

	if recordAttributes.Type != "" {
		payload.Add("record_type", recordAttributes.Type)
	}

	if recordAttributes.Line != "" {
		payload.Add("record_line", recordAttributes.Line)
	}

	if recordAttributes.LineID != "" {
		payload.Add("record_line_id", recordAttributes.LineID)
	}

	if recordAttributes.Value != "" {
		payload.Add("value", recordAttributes.Value)
	}

	if recordAttributes.MX != "" {
		payload.Add("mx", recordAttributes.MX)
	}

	if recordAttributes.TTL != "" {
		payload.Add("ttl", recordAttributes.TTL)
	}

	if recordAttributes.Status != "" {
		payload.Add("status", recordAttributes.Status)
	}

	returnedRecord := recordWrapper{}

	res, err := s.client.post(path, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("Could not get domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// GetRecord fetches the domain record.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-info
func (s *DomainsService) GetRecord(domain string, recordID string) (Record, *Response, error) {
	path := recordAction("Info")

	payload := newPayLoad(s.client.CommonParams)

	payload.Add("domain_id", domain)
	payload.Add("record_id", recordID)

	returnedRecord := recordWrapper{}

	res, err := s.client.post(path, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("Could not get domains: %s", returnedRecord.Status.Message)
	}
	record := returnedRecord.Record
	if record.Type == "" && record.RecordType != "" {
		record.Type = record.RecordType
	}
	if record.Line == "" && record.RecordLine != "" {
		record.Line = record.RecordLine
	}
	if record.LineID == "" && record.RecordLineID != "" {
		record.LineID = record.RecordLineID
	}
	if record.Name == "" && record.SubDomain != "" {
		record.Name = record.SubDomain
	}
	return record, res, nil
}

// UpdateRecord updates a domain record.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-modify
func (s *DomainsService) UpdateRecord(domain string, recordID string, recordAttributes Record) (Record, *Response, error) {
	path := recordAction("Modify")

	payload := newPayLoad(s.client.CommonParams)

	payload.Add("domain_id", domain)

	payload.Add("record_id", recordID)

	if recordAttributes.Name != "" {
		payload.Add("sub_domain", recordAttributes.Name)
	}

	if recordAttributes.Type != "" {
		payload.Add("record_type", recordAttributes.Type)
	}

	if recordAttributes.Line != "" {
		payload.Add("record_line", recordAttributes.Line)
	}

	if recordAttributes.LineID != "" {
		payload.Add("record_line_id", recordAttributes.LineID)
	}

	if recordAttributes.Value != "" {
		payload.Add("value", recordAttributes.Value)
	}

	if recordAttributes.MX != "" {
		payload.Add("mx", recordAttributes.MX)
	}

	if recordAttributes.TTL != "" {
		payload.Add("ttl", recordAttributes.TTL)
	}

	if recordAttributes.Status != "" {
		payload.Add("status", recordAttributes.Status)
	}

	returnedRecord := recordWrapper{}

	res, err := s.client.post(path, payload, &returnedRecord)
	if err != nil {
		return Record{}, res, err
	}

	if returnedRecord.Status.Code != "1" {
		return returnedRecord.Record, nil, fmt.Errorf("Could not get domains: %s", returnedRecord.Status.Message)
	}

	return returnedRecord.Record, res, nil
}

// DeleteRecord deletes a domain record.
//
// dnspod API docs: https://www.dnspod.cn/docs/records.html#record-remove
func (s *DomainsService) DeleteRecord(domain string, recordID string) (*Response, error) {
	path := recordAction("Remove")

	payload := newPayLoad(s.client.CommonParams)

	payload.Add("domain_id", domain)
	payload.Add("record_id", recordID)

	returnedRecord := recordWrapper{}

	res, err := s.client.post(path, payload, &returnedRecord)
	if err != nil {
		return res, err
	}

	if returnedRecord.Status.Code != "1" {
		return nil, fmt.Errorf("could not delete record: %s", returnedRecord.Status.Message)
	}

	return res, nil
}

func (s *DomainsService) UpdateRecordStatus(domainID string, recordID string, status string) (*Response, error) {
	path := recordAction("Status")
	payload := newPayLoad(s.client.CommonParams)
	payload.Add("domain_id", domainID)
	payload.Add("record_id", recordID)
	payload.Add("status", status)

	returnedRecord := recordWrapper{}

	res, err := s.client.post(path, payload, &returnedRecord)
	if err != nil {
		return res, err
	}
	if returnedRecord.Status.Code != "1" {
		return nil, fmt.Errorf("could not change record status: %s", returnedRecord.Status.Message)
	}
	return res, nil
}

func (s *DomainsService) GetRecordLine(domainGrade string, domainID string) ([]RecordLine, *Response, error) {
	path := recordAction("Line")
	payload := newPayLoad(s.client.CommonParams)
	payload.Set("domain_grade", domainGrade)
	payload.Set("domain_id", domainID)
	lines := linesWrapper{}
	res, err := s.client.post(path, payload, &lines)
	if err != nil {
		return []RecordLine{}, res, err
	}
	if lines.Status.Code != "1" {
		return []RecordLine{}, nil, fmt.Errorf("could not get record line: %s", lines.Status.Message)
	}
	var ret []RecordLine
	for k, v := range lines.LineIDs {
		s, ok := v.(string)
		if !ok && k == "默认" {
			ret = append(ret, RecordLine{
				Line: k,
				LineID: "0",
			})
		} else {
			ret = append(ret, RecordLine{
				Line: k,
				LineID: s,
			})
		}
	}
	return ret, res, nil
}

