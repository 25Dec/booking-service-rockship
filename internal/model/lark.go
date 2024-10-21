package model

type LarkAppToken struct {
	AppAccessToken    string `json:"app_access_token"`
	Code              int    `json:"code"`
	Expire            int64  `json:"expire"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}

type LarkResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
type LarkEventData struct {
	Event LarkEvent `json:"event"`
}

type LarkEvent struct {
	EventID             string `json:"event_id"`
	OrganizerCalendarID string `json:"organizer_calendar_id"`
	Summary             string `json:"summary"`
	Description         string `json:"description"`
	NeedNotification    bool   `json:"need_notification"`
	StartTime           struct {
		Date      string `json:"date"`
		Timestamp string `json:"timestamp"`
		Timezone  string `json:"timezone"`
	} `json:"start_time"`
	EndTime struct {
		Date      string `json:"date"`
		Timestamp string `json:"timestamp"`
		Timezone  string `json:"timezone"`
	} `json:"end_time"`
	Vchat struct {
		VcType      string `json:"vc_type"`
		IconType    string `json:"icon_type"`
		Description string `json:"description"`
		MeetingURL  string `json:"meeting_url"`
	} `json:"vchat"`
	CreateTime string `json:"create_time"`
}

type LarkEventRequest struct {
	Summary         string             `json:"summary"`
	Description     string             `json:"description"`
	StartTime       LarkEventTime      `json:"start_time"`
	EndTime         LarkEventTime      `json:"end_time"`
	Vchat           LarkVchat          `json:"vchat"`
	Visibility      string             `json:"visibility"`
	AttendeeAbility string             `json:"attendee_ability"`
	FreeBusyStatus  string             `json:"free_busy_status"`
	Color           int                `json:"color"`
	Reminders       []LarkEventRemider `json:"reminders"`
}

func NewLarkEventRequest(parameter map[string]any) LarkEventRequest {
	//setup defualt value
	if parameter["reminders"] == nil {
		parameter["reminders"] = []LarkEventRemider{
			{
				Minutes: 5,
			},
		}
	}
	if parameter["color"] == nil {
		parameter["color"] = -1
	}
	if parameter["visibility"] == nil {
		parameter["visibility"] = "default"
	}
	if parameter["attendee_ability"] == nil {
		parameter["attendee_ability"] = "can_see_others"
	}
	if parameter["free_busy_status"] == nil {
		parameter["free_busy_status"] = "busy"
	}
	if parameter["vchat"] == nil {
		parameter["vchat"] = LarkVchat{
			VcType:      "vc",
			IconType:    "vc",
			Description: "description",
		}
	}

	return LarkEventRequest{
		Summary:         parameter["summary"].(string),
		Description:     parameter["description"].(string),
		StartTime:       parameter["start_time"].(LarkEventTime),
		EndTime:         parameter["end_time"].(LarkEventTime),
		Vchat:           parameter["vchat"].(LarkVchat),
		Visibility:      parameter["visibility"].(string),
		AttendeeAbility: parameter["attendee_ability"].(string),
		FreeBusyStatus:  parameter["free_busy_status"].(string),
		Color:           parameter["color"].(int),
		Reminders:       parameter["reminders"].([]LarkEventRemider),
	}
}

type LarkEventRemider struct {
	Minutes int `json:"minutes"`
}

type LarkVchat struct {
	VcType      string `json:"vc_type"`
	IconType    string `json:"icon_type"`
	Description string `json:"description"`
}

type LarkEventTime struct {
	Timestamp string `json:"timestamp"`
	TimeZone  string `json:"timezone"`
}

type LarkEventAttendee struct {
	Type            string `json:"type"`
	ThirdPartyEmail string `json:"third_party_email"`
}
