package openapis

import (
	"booking-service/internal/model"
	"booking-service/pkg/utils/errs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type edtronautAPI struct {
	Domain string
}

func NewEdtronautAPI(domain string) EdtronautAPI {
	return &edtronautAPI{
		Domain: domain,
	}
}

func (e *edtronautAPI) GetUserByID(id string) (model.EdtronautUser, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", e.Domain+"/users/"+id, nil)
	if err != nil {
		return model.EdtronautUser{}, fmt.Errorf("EdtronautWebAPI - GetUser - http.NewRequest: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return model.EdtronautUser{}, fmt.Errorf("EdtronautWebAPI - GetUser - httpClient.Do: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return model.EdtronautUser{}, errs.NotFoundError{Message: "User not found"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.EdtronautUser{}, fmt.Errorf("EdtronautWebAPI - GetUser - ioutil.ReadAll: %w", err)
	}
	var response model.EdtronautAPIResponse[model.EdtronautUser]
	err = json.Unmarshal(body, &response)
	if err != nil {
		return model.EdtronautUser{}, fmt.Errorf("EdtronautWebAPI - GetUser - json.Unmarshal: %w", err)
	}
	return response.Data, nil
}

func (e *edtronautAPI) GetUserByToken(token string) (model.EdtronautUser, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", e.Domain+"/auth/me", nil)
	if err != nil {
		return model.EdtronautUser{}, fmt.Errorf("EdtronautWebAPI - GetUser - http.NewRequest: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return model.EdtronautUser{}, fmt.Errorf("EdtronautWebAPI - GetUser - httpClient.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return model.EdtronautUser{}, errs.NotFoundError{Message: "User not found"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.EdtronautUser{}, fmt.Errorf("EdtronautWebAPI - GetUser - ioutil.ReadAll: %w", err)
	}
	var response model.EdtronautAPIResponse[model.EdtronautUser]
	err = json.Unmarshal(body, &response)
	if err != nil {
		return model.EdtronautUser{}, fmt.Errorf("EdtronautWebAPI - GetUser - json.Unmarshal: %w", err)
	}
	return response.Data, nil
}

func (e *edtronautAPI) GetUserCourses(userToken string) ([]model.EdtronautCourse, error) {
	curPage := 1
	var courses []model.EdtronautCourse
	for {
		httpClient := &http.Client{}
		req, err := http.NewRequest("GET", e.Domain+"/courses/my-courses?page="+fmt.Sprint(curPage), nil)
		if err != nil {
			return nil, fmt.Errorf("EdtronautWebAPI - GetUserCourses - http.NewRequest: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+userToken)
		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("EdtronautWebAPI - GetUserCourses - httpClient.Do: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, errs.NotFoundError{Message: "Courses not found"}
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("EdtronautWebAPI - GetUserCourses - ioutil.ReadAll: %w", err)
		}

		var response model.EdtronautPaginationResponse[model.EdtronautCourse]
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("EdtronautWebAPI - GetUserCourses - json.Unmarshal: %w", err)
		}
		courses = append(courses, response.Data.Courses...)
		if len(response.Data.Courses) == 0 {
			break
		}
		curPage++
	}
	return courses, nil
}
