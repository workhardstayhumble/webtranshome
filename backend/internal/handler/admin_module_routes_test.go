package handler

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNewAdminModuleRoutesHavePermissions(t *testing.T) {
	b, err := os.ReadFile("routes.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(b)
	for _, path := range []string{
		"/admin/banners", "/admin/articles", "/admin/members", "/admin/members/search",
		"/admin/points", "/admin/qontak/contacts", "/admin/activity-logs",
		"/admin/activity-logs/options",
	} {
		if !strings.Contains(src, `Path: "`+path+`"`) {
			t.Errorf("missing route: %s", path)
		}
	}
	for _, p := range []string{
		"ContentBannerReadCode", "ContentBannerManageCode", "ContentArticleReadCode", "ContentArticleManageCode",
		"MembershipMemberReadCode", "MembershipMemberUpdateCode", "MembershipPointsReadCode", "MembershipPointsInputCode",
		"IntegrationQontakReadCode", "IntegrationQontakRetryCode", "AuditReadCode",
	} {
		if !strings.Contains(src, "serverCtx.RequirePermission(permission."+p+")") {
			t.Errorf("missing permission check: %s", p)
		}
	}
}
func TestAdminPaginationBounds(t *testing.T) {
	r := httptest.NewRequest("GET", "/admin/members?page=-1&per_page=1000", nil)
	p, n, offset := pageParams(r)
	if p != 1 || n != 100 || offset != 0 {
		t.Fatalf("unexpected pagination: %d %d %d", p, n, offset)
	}
}
func TestBannerDateOrder(t *testing.T) {
	_, _, ok := bannerDates(bannerInput{StartDate: "2026-10-12", EndDate: "2026-10-11"})
	if ok {
		t.Fatal("end date before start should fail")
	}
	_, _, ok = bannerDates(bannerInput{StartDate: "2026-10-12", EndDate: "2026-10-13"})
	if !ok {
		t.Fatal("valid date range rejected")
	}
}