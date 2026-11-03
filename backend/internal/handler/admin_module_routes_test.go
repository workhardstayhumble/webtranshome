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
		"/admin/activity-logs/options", "/admin/services", "/admin/services/:id",
	} {
		if !strings.Contains(src, `Path: "`+path+`"`) {
			t.Errorf("missing route: %s", path)
		}
	}
	for _, p := range []string{
		"ContentBannerReadCode", "ContentBannerManageCode", "ContentArticleReadCode", "ContentArticleManageCode",
		"MembershipMemberReadCode", "MembershipMemberUpdateCode", "MembershipPointsReadCode", "MembershipPointsInputCode",
		"IntegrationQontakReadCode", "IntegrationQontakRetryCode", "AuditReadCode",
		"ContentServiceReadCode", "ContentServiceManageCode",
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
func TestServiceValidation(t *testing.T) {
	area := strings.Repeat("a", 101)
	in := serviceInput{Name: "  ", Slug: "Bukan Slug", Summary: strings.Repeat("x", 161), ServiceArea: &area, CategoryIDs: []int64{0}}
	errs := validateService(&in)
	for _, field := range []string{"name", "slug", "summary", "service_area", "category_ids"} {
		if _, ok := errs[field]; !ok {
			t.Errorf("expected error for %s", field)
		}
	}
	in = serviceInput{Name: " Jasa pengecatan ", Slug: "pengecatan", Summary: "Cat dinding.", CategoryIDs: []int64{5, 5, 2}}
	if errs := validateService(&in); len(errs) > 0 {
		t.Fatalf("valid service rejected: %v", errs)
	}
	if in.Name != "Jasa pengecatan" || len(in.CategoryIDs) != 2 {
		t.Fatalf("input not normalised: %q %v", in.Name, in.CategoryIDs)
	}
}
func TestServiceChangesOnlyListsChangedFields(t *testing.T) {
	before := map[string]any{"id": 1, "name": "A", "is_active": true, "category_ids": []int64{1}}
	after := map[string]any{"id": 1, "name": "A", "is_active": false, "category_ids": []int64{1}}
	got := serviceChanges(before, after)
	if _, ok := got["is_active"]; !ok {
		t.Fatal("is_active change missing")
	}
	if _, ok := got["category_ids"]; ok {
		t.Fatal("unchanged category_ids logged")
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