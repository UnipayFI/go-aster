package futures

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-aster/v3/request"
)

// AnnouncementContent is one language variant of an announcement.
type AnnouncementContent struct {
	Language string `json:"language"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

// Announcement is a single direct (personal) announcement.
type Announcement struct {
	ID          int64                 `json:"id"`
	Contents    []AnnouncementContent `json:"contents"`
	Category    string                `json:"category"`
	PublishTime time.Time             `json:"publishTime,format:unixmilli"`
	JumpLink    string                `json:"jumpLink"`
}

// GetDirectAnnouncementsService -- GET /fapi/v3/announcement/direct (USER_DATA)
type GetDirectAnnouncementsService struct {
	c      *FuturesClient
	params map[string]string
}

func (c *FuturesClient) NewGetDirectAnnouncementsService() *GetDirectAnnouncementsService {
	return &GetDirectAnnouncementsService{c: c, params: map[string]string{}}
}

// SetPage sets the 1-based page number (default 1).
func (s *GetDirectAnnouncementsService) SetPage(page int) *GetDirectAnnouncementsService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetSize sets the number of results per page (default 20).
func (s *GetDirectAnnouncementsService) SetSize(size int) *GetDirectAnnouncementsService {
	s.params["size"] = strconv.Itoa(size)
	return s
}

func (s *GetDirectAnnouncementsService) Do(ctx context.Context) (*DirectAnnouncements, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/announcement/direct", s.params).WithSignature()
	return request.Do[DirectAnnouncements](req)
}

type DirectAnnouncements struct {
	Total int64          `json:"total"`
	Rows  []Announcement `json:"rows"`
}

// GetDirectAnnouncementByIdService -- GET /fapi/v3/announcement/directById (USER_DATA)
type GetDirectAnnouncementByIdService struct {
	c  *FuturesClient
	id int64
}

func (c *FuturesClient) NewGetDirectAnnouncementByIdService(id int64) *GetDirectAnnouncementByIdService {
	return &GetDirectAnnouncementByIdService{c: c, id: id}
}

func (s *GetDirectAnnouncementByIdService) Do(ctx context.Context) (*Announcement, error) {
	req := request.Get(ctx, s.c, "/fapi/v3/announcement/directById", map[string]string{
		"id": strconv.FormatInt(s.id, 10),
	}).WithSignature()
	return request.Do[Announcement](req)
}
