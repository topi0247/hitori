package domain_test

import (
	"testing"

	"github.com/topi0247/hitori/domain"
)

// Judge はカードの並び替え結果を判定する
// 正解順（card_number昇順）と照らし合わせてマッチポイントと正解率を算出する

func TestJudge_ExactMatch(t *testing.T) {
	cards := []domain.CardEntry{
		{UUID: "a", CardNumber: 10},
		{UUID: "b", CardNumber: 30},
		{UUID: "c", CardNumber: 50},
	}
	answers := []domain.GameAnswer{
		{UUID: "a", Order: 1},
		{UUID: "b", Order: 2},
		{UUID: "c", Order: 3},
	}

	result, err := domain.Judge(cards, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range result.Cards {
		if !r.IsCorrect {
			t.Errorf("UUID=%s: 全正解のはずが IsCorrect=false", r.UUID)
		}
		if r.MatchPoints != domain.PointsExact {
			t.Errorf("UUID=%s: MatchPoints=%d, want %d", r.UUID, r.MatchPoints, domain.PointsExact)
		}
	}
	if result.CorrectRate != 100.0 {
		t.Errorf("CorrectRate=%v, want 100.0", result.CorrectRate)
	}
}

func TestJudge_AdjacentMatch(t *testing.T) {
	// 正解順: a(10), b(30), c(50) → ユーザーは b,a,c の順で提出
	cards := []domain.CardEntry{
		{UUID: "a", CardNumber: 10},
		{UUID: "b", CardNumber: 30},
		{UUID: "c", CardNumber: 50},
	}
	answers := []domain.GameAnswer{
		{UUID: "b", Order: 1},
		{UUID: "a", Order: 2},
		{UUID: "c", Order: 3},
	}

	result, err := domain.Judge(cards, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// a: 正解order=1, ユーザーorder=2 → 隣
	// b: 正解order=2, ユーザーorder=1 → 隣
	// c: 正解order=3, ユーザーorder=3 → ぴったり
	for _, r := range result.Cards {
		switch r.UUID {
		case "a", "b":
			if r.IsCorrect {
				t.Errorf("UUID=%s: 隣なのに IsCorrect=true", r.UUID)
			}
			if r.MatchPoints != domain.PointsAdjacent {
				t.Errorf("UUID=%s: MatchPoints=%d, want %d", r.UUID, r.MatchPoints, domain.PointsAdjacent)
			}
		case "c":
			if !r.IsCorrect {
				t.Errorf("UUID=%s: 正解なのに IsCorrect=false", r.UUID)
			}
			if r.MatchPoints != domain.PointsExact {
				t.Errorf("UUID=%s: MatchPoints=%d, want %d", r.UUID, r.MatchPoints, domain.PointsExact)
			}
		}
	}
}

func TestJudge_NoMatch(t *testing.T) {
	// 正解順: a(10), b(30), c(50), d(70) → 全て外れ
	cards := []domain.CardEntry{
		{UUID: "a", CardNumber: 10},
		{UUID: "b", CardNumber: 30},
		{UUID: "c", CardNumber: 50},
		{UUID: "d", CardNumber: 70},
	}
	answers := []domain.GameAnswer{
		{UUID: "c", Order: 1},
		{UUID: "d", Order: 2},
		{UUID: "a", Order: 3},
		{UUID: "b", Order: 4},
	}

	result, err := domain.Judge(cards, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range result.Cards {
		if r.IsCorrect {
			t.Errorf("UUID=%s: 不正解のはずが IsCorrect=true", r.UUID)
		}
		if r.MatchPoints != 0 {
			t.Errorf("UUID=%s: MatchPoints=%d, want 0", r.UUID, r.MatchPoints)
		}
	}
	if result.CorrectRate != 0.0 {
		t.Errorf("CorrectRate=%v, want 0.0", result.CorrectRate)
	}
}

func TestJudge_CorrectRate(t *testing.T) {
	// 3枚中1枚正解 → 33.33%
	cards := []domain.CardEntry{
		{UUID: "a", CardNumber: 10},
		{UUID: "b", CardNumber: 30},
		{UUID: "c", CardNumber: 50},
	}
	answers := []domain.GameAnswer{
		{UUID: "a", Order: 1},
		{UUID: "c", Order: 2},
		{UUID: "b", Order: 3},
	}

	result, err := domain.Judge(cards, answers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 33.33
	if result.CorrectRate != want {
		t.Errorf("CorrectRate=%v, want %v", result.CorrectRate, want)
	}
}

func TestJudge_InvalidAnswer(t *testing.T) {
	cards := []domain.CardEntry{
		{UUID: "a", CardNumber: 10},
	}
	// 存在しないUUIDを渡す
	answers := []domain.GameAnswer{
		{UUID: "unknown", Order: 1},
	}

	_, err := domain.Judge(cards, answers)
	if err == nil {
		t.Error("存在しないUUIDに対してエラーが返るべき")
	}
}
