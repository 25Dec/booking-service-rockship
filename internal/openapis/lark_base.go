package openapis

import (
	"booking-service/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type larkBaseAPI struct {
	AppID          string `json:"app_id"`
	SecretKey      string `json:"secrect_key"`
	BaseToken      string `json:"base_token"`
	appToken       model.LarkAppToken
	larkBaseTables map[string]string
}

func NewLarkBaseAPI(appID string, secretKey string, baseToken string, tables map[string]string) LarkBaseAPI {
	return &larkBaseAPI{
		AppID:          appID,
		SecretKey:      secretKey,
		BaseToken:      baseToken,
		appToken:       model.LarkAppToken{},
		larkBaseTables: tables,
	}
}

func (l *larkBaseAPI) ObtainAppToken() (string, error) {
	if l.appToken.Expire < time.Now().Unix() || l.appToken.AppAccessToken == "" {
		httpClient := &http.Client{}
		//https://open.larksuite.com/open-apis/auth/v3/tenant_access_token/internal
		req, err := http.NewRequest("POST", "https://open.larksuite.com/open-apis/auth/v3/tenant_access_token/internal", nil)
		if err != nil {
			return "", fmt.Errorf("LarkCalendarWebAPI - ObtainAppToken - http.NewRequest: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"app_id":"%s","app_secret":"%s"}`, l.AppID, l.SecretKey)))

		req.Body = io.NopCloser(reqBody)
		resp, err := httpClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("LarkCalendarWebAPI - ObtainAppToken - httpClient.Do: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("http status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("LarkCalendarWebAPI - ObtainAppToken - ioutil.ReadAll: %w", err)
		}
		var token model.LarkAppToken
		err = json.Unmarshal(body, &token)
		if err != nil {
			return "", fmt.Errorf("LarkCalendarWebAPI - ObtainAppToken - json.Unmarshal: %w", err)
		}

		token.Expire = time.Now().Unix() + token.Expire - 300
		l.appToken = token
	}

	return l.appToken.TenantAccessToken, nil
}

func (l *larkBaseAPI) AppendRecord(record map[string]interface{}, tableName string) (string, error) {
	token, err := l.ObtainAppToken()
	if err != nil {
		return "", fmt.Errorf("LarkBaseWebAPI - AppendRecord - l.ObtainAppToken: %w", err)
	}
	tableID := l.larkBaseTables[tableName]
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://open.larksuite.com/open-apis/bitable/v1/apps/%s/tables/%s/records", l.BaseToken, tableID), nil)
	if err != nil {
		return "", fmt.Errorf("LarkBaseWebAPI - AppendRecord - http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	jsonRecord, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("LarkBaseWebAPI - AppendRecord - json.Marshal: %w", err)
	}

	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"fields":%s}`, jsonRecord)))

	req.Body = io.NopCloser(reqBody)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("LarkBaseWebAPI - AppendRecord - httpClient.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status code 1: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("LarkBaseWebAPI - AppendRecord - ioutil.ReadAll: %w", err)
	}

	return string(body), nil
}
