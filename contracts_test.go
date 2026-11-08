package engine_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/altnovatrail/engine"
)

// The engine's ReasonCode, ChoiceKind and RunEnd carry the names and numbers of ReasonCode,
// ChoiceKind and RunEndReason in cause.proto, in spec's contracts, which api translates them by
// (docs/architecture.md, section 3.6; R09). The file is read in the spec submodule, at the tag of
// spec's release this repository pins (R64 spec): it is no part of the engine's Go module. Where
// the submodule is not checked out — the shared go.yml checks out none, and a consumer runs the
// tests from the module cache — the test is skipped; the job Contracts of ci.yml checks the
// submodule out and sets ENGINE_CONTRACTS=1, and there a missing file fails.
const causeProto = "spec/contracts/novatrail/api/v1/cause.proto"

var contractEnums = map[string]map[string]int32{
	"ReasonCode": {
		"REASON_CODE_UNSPECIFIED":   int32(engine.ReasonUnspecified),
		"REASON_CODE_CONSOLIDATION": int32(engine.ReasonConsolidation),
		"REASON_CODE_NEW_SKILL":     int32(engine.ReasonNewSkill),
		"REASON_CODE_QUEUE":         int32(engine.ReasonQueue),
		"REASON_CODE_REPETITION":    int32(engine.ReasonRepetition),
		"REASON_CODE_PROBE":         int32(engine.ReasonProbe),
		"REASON_CODE_SKILL_CHOSEN":  int32(engine.ReasonSkillChosen),
		"REASON_CODE_AUTHOR_ORDER":  int32(engine.ReasonAuthorOrder),
		"REASON_CODE_OPEN_RUN":      int32(engine.ReasonOpenRun),
		"REASON_CODE_ROTATION":      int32(engine.ReasonRotation),
	},
	"ChoiceKind": {
		"CHOICE_KIND_UNSPECIFIED": int32(engine.ChoiceUnspecified),
		"CHOICE_KIND_LESSON":      int32(engine.ChoiceLesson),
		"CHOICE_KIND_TASK":        int32(engine.ChoiceTask),
		"CHOICE_KIND_NOTHING":     int32(engine.ChoiceNothing),
	},
	"RunEndReason": {
		"RUN_END_REASON_UNSPECIFIED":  int32(engine.RunOpen),
		"RUN_END_REASON_PASSED":       int32(engine.RunPassed),
		"RUN_END_REASON_UNCHECKED":    int32(engine.RunUnchecked),
		"RUN_END_REASON_NOT_SUITABLE": int32(engine.RunNotSuitable),
		"RUN_END_REASON_ABANDONED":    int32(engine.RunAbandoned),
	},
}

func TestContracts(t *testing.T) {
	data, err := os.ReadFile(filepath.FromSlash(causeProto))
	if errors.Is(err, os.ErrNotExist) && os.Getenv("ENGINE_CONTRACTS") != "1" {
		t.Skipf("%s is not here: the spec submodule is not checked out", causeProto)
	}
	if err != nil {
		t.Fatal(err)
	}
	enums, err := protoEnums(string(data))
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range contractEnums {
		got, ok := enums[name]
		if !ok {
			t.Errorf("%s has no enum %s", causeProto, name)
			continue
		}
		for value, n := range want {
			if m, ok := got[value]; !ok {
				t.Errorf("%s: %s has no %s, which the engine names %d", causeProto, name, value, n)
			} else if m != n {
				t.Errorf("%s: %s is %d there and %d in the engine", causeProto, value, m, n)
			}
		}
		for value, m := range got {
			if _, ok := want[value]; !ok {
				t.Errorf("%s: %s has %s = %d, which the engine does not name", causeProto, name, value, m)
			}
		}
	}
}

// protoEnums reads the values of every enum of a proto3 file, by the enum's name: the lines
// "NAME = N;" between "enum Name {" and its "}", past comments. It reads as much of proto3 as
// cause.proto needs, and refuses a line in an enum it cannot read rather than skip it.
func protoEnums(src string) (map[string]map[string]int32, error) {
	enums := map[string]map[string]int32{}
	var current string
	for i, line := range strings.Split(src, "\n") {
		line, _, _ = strings.Cut(line, "//")
		line = strings.Trim(line, " \t\r")
		switch {
		case current == "" && strings.HasPrefix(line, "enum "):
			name, ok := strings.CutSuffix(strings.TrimPrefix(line, "enum "), " {")
			if !ok || name == "" || strings.ContainsAny(name, " \t{}") {
				return nil, errors.New("line " + strconv.Itoa(i+1) + ": an enum this test cannot read: " + line)
			}
			current = name
			enums[current] = map[string]int32{}
		case current == "" || line == "":
		case line == "}":
			current = ""
		default:
			value, number, ok := strings.Cut(strings.TrimSuffix(line, ";"), " = ")
			n, err := strconv.ParseInt(number, 10, 32)
			if !ok || !strings.HasSuffix(line, ";") || err != nil || slices.Contains([]string{"", "option", "reserved"}, value) {
				return nil, errors.New("line " + strconv.Itoa(i+1) + ": a line of enum " + current + " this test cannot read: " + line)
			}
			enums[current][value] = int32(n)
		}
	}
	if current != "" {
		return nil, errors.New("enum " + current + " has no end")
	}
	return enums, nil
}

// protoEnums reads what it should and refuses what it cannot read, so that a change of the file
// that the test would misread fails rather than passes.
func TestProtoEnums(t *testing.T) {
	got, err := protoEnums("syntax = \"proto3\";\n// enum Commented {\nenum A {\n  // a comment\n  A_ZERO = 0; // trailing\n\n  A_ONE = 1;\r\n}\nmessage M {\n  int32 x = 1;\n}\nenum B {\n  B_ZERO = 0;\n}\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || len(got["A"]) != 2 || got["A"]["A_ONE"] != 1 || got["B"]["B_ZERO"] != 0 {
		t.Errorf("protoEnums = %v", got)
	}
	for _, src := range []string{
		"enum A {\n  A_ZERO = 0;\n",                        // no end
		"enum A {\n  reserved 2;\n}\n",                     // a line it cannot read
		"enum A {\n  A_ZERO = zero;\n}\n",                  // not a number
		"enum A {\n  A_ZERO = 0\n}\n",                      // no semicolon
		"enum A {\n  option allow_alias = true;\n}\n",      // an option
		"enum A { A_ZERO = 0; }\n",                         // all on one line
		"enum {\n}\n",                                      // no name
		"enum A {\n  A_ZERO = 0 [deprecated = true];\n}\n", // an option of a value
	} {
		if _, err := protoEnums(src); err == nil {
			t.Errorf("protoEnums(%q) read it", src)
		}
	}
}
