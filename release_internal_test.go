package engine

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The packages of the first slice, from the release v0.1.0 of altnovatrail/modules, which
// release-verify accepted on 2026-11-08 (docs/architecture.md, section 4): they stand in
// testdata/modules as the release has them, and a chain of digests holds them from then on. The
// digest of the release's module index is pinned here; the index names each package's digest,
// which is the SHA-256 of its inventory.json; and the inventory names the SHA-256 of server.json
// (the format's text, sections 25 and 26). Reading a package checks every link.
//
// The translation of server.json into the engine's types is this test's, not api's or sim's: the
// golden vectors do not catch a difference between them.

const moduleIndexDigest = "ed4d8bb325fdcb3ea54da439e08c58d927f773ba0f097338f7767ad76be8c90e"

// The three modules of the first slice, by full ID.
const (
	releaseOlympiad = "novatrail/olympiad-math"
	releaseGrade3   = "novatrail/grade-3-math"
	releaseEssays   = "novatrail/english-essays"
)

// release is a module of the release as the engine gets it: its graph, its pool and its lessons,
// each task and lesson with the age categories it suits beside, for a scenario to keep only what
// api would give a learner.
type release struct {
	module  Module
	ages    []string            // the module's age categories, its audience
	taskAge map[string][]string // a task's own age categories, by full ID; none, the module's
	lessAge map[string][]string // a lesson's own, by full ID; none, the module's
}

type serverPart struct {
	Manifest struct {
		ID       string `json:"id"`
		Audience struct {
			AgeCategories []string `json:"age_categories"`
		} `json:"audience"`
		Bands []struct {
			ID     string  `json:"id"`
			Grades []int32 `json:"grades"`
		} `json:"bands"`
	} `json:"manifest"`
	SkillGraph struct {
		Skills map[string]struct {
			Assessability string   `json:"assessability"`
			Prerequisites []string `json:"prerequisites"`
			Teaches       []string `json:"teaches"`
		} `json:"skills"`
	} `json:"skill_graph"`
	SkillTiers map[string]int32 `json:"skill_tiers"`
	Tasks      map[string]struct {
		Band          string     `json:"band"`
		Difficulty    int32      `json:"difficulty"`
		AgeCategories []string   `json:"age_categories"`
		Skills        skillTags  `json:"skills"`
		Items         []taskItem `json:"items"`
	} `json:"tasks"`
	Lessons map[string]struct {
		Skills        []string `json:"skills"`
		AgeCategories []string `json:"age_categories"`
		Segments      []struct {
			Task string `json:"task"`
		} `json:"segments"`
	} `json:"lessons"`
}

type skillTags struct {
	Primary    []string `json:"primary"`
	Supporting []string `json:"supporting"`
}

type taskItem struct {
	Skills      *skillTags   `json:"skills"`
	AnswerSlots []answerSlot `json:"answer_slots"`
	Line        []answerSlot `json:"line"` // the parts of a line of a workbook, its slots among them
}

type answerSlot struct {
	Options []string `json:"options"`
	Check   *struct {
		Kind       string   `json:"kind"`
		Options    []string `json:"options"`
		TrapValues []struct {
			Trap string `json:"trap"`
		} `json:"trap_values"`
	} `json:"check"`
}

// releaseRead holds the release once read: the tests run one at a time, and every test that
// reads it gets the same.
var releaseRead map[string]release

// readRelease reads the three packages of the release, by module ID, checking the chain of
// digests from the pinned index; once a run.
func readRelease(t *testing.T) map[string]release {
	t.Helper()
	if releaseRead != nil {
		return releaseRead
	}
	dir := filepath.Join("testdata", "modules")
	index := readFile(t, filepath.Join(dir, "module-index.json"))
	if got := fmt.Sprintf("%x", sha256.Sum256(index)); got != moduleIndexDigest {
		t.Fatalf("module-index.json has the digest %s, not %s", got, moduleIndexDigest)
	}
	var idx struct {
		Modules map[string]map[string]struct {
			PackageDigest string `json:"package_digest"`
		} `json:"modules"`
	}
	if err := json.Unmarshal(index, &idx); err != nil {
		t.Fatal(err)
	}
	out := map[string]release{}
	for _, id := range []string{releaseOlympiad, releaseGrade3, releaseEssays} {
		digest := strings.TrimPrefix(idx.Modules[id]["0.1.0"].PackageDigest, "sha256:")
		files := untar(t, filepath.Join(dir, "sha256-"+digest+".tar.gz"))
		inventory, server := files["inventory.json"], files["server.json"]
		if got := fmt.Sprintf("%x", sha256.Sum256(inventory)); got != digest {
			t.Fatalf("%s: inventory.json has the digest %s, not the package's %s", id, got, digest)
		}
		var inv struct {
			ServerPart struct {
				Digest string `json:"digest"`
			} `json:"server_part"`
		}
		if err := json.Unmarshal(inventory, &inv); err != nil {
			t.Fatal(err)
		}
		if got := "sha256:" + fmt.Sprintf("%x", sha256.Sum256(server)); got != inv.ServerPart.Digest {
			t.Fatalf("%s: server.json has the digest %s, not the inventory's %s", id, got, inv.ServerPart.Digest)
		}
		out[id] = translate(t, server)
	}
	releaseRead = out
	return out
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// untar gives the files of a package's archive, by name. The archive is a gzip of a tar of
// 512-byte blocks: a header — the name in its first 100 bytes, the size in octal at 124, the type at
// 156 — then the file's blocks; the packages name their files short, with no ustar prefix. The test
// reads it itself, as archive/tar would pull os/user, whose cgo the rules test cannot type-check for
// every target.
func untar(t *testing.T, path string) map[string][]byte {
	t.Helper()
	zr, err := gzip.NewReader(bytes.NewReader(readFile(t, path)))
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if _, err := b.ReadFrom(zr); err != nil {
		t.Fatal(err)
	}
	data := b.Bytes()
	files := map[string][]byte{}
	for len(data) >= 512 && data[0] != 0 {
		h := data[:512]
		name := string(bytes.TrimRight(h[:100], "\x00"))
		size, err := strconv.ParseInt(string(bytes.Trim(h[124:136], " \x00")), 8, 64)
		next := 512 + (size+511)/512*512
		if err != nil || size < 0 || int64(len(data)) < next {
			t.Fatalf("%s: a broken header or a short block at %q", path, name)
		}
		if h[156] == '0' || h[156] == 0 {
			files[name] = data[512 : 512+size]
		}
		data = data[next:]
	}
	return files
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// translate reads a server part into the engine's types: full IDs, an entry's module, its kind
// and its local ID joined by / (the format's text, section 10); tiers from skill_tiers; bands by
// place from the manifest; a slot's kind of check, and a choice's number of options; a check by a
// rubric gives no slot with a verdict.
func translate(t *testing.T, server []byte) release {
	t.Helper()
	var s serverPart
	if err := json.Unmarshal(server, &s); err != nil {
		t.Fatal(err)
	}
	module := s.Manifest.ID
	full := func(kind, id string) string {
		if strings.HasPrefix(id, anchorPrefix) {
			return id
		}
		return module + "/" + kind + "/" + id
	}
	fulls := func(kind string, ids []string) []string {
		var out []string
		for _, id := range ids {
			out = append(out, full(kind, id))
		}
		return out
	}

	var bands []learnermodel.Band
	place := map[string]int32{}
	for k, b := range s.Manifest.Bands {
		bands = append(bands, learnermodel.Band{Place: int32(k + 1), ID: b.ID, Grades: b.Grades})
		place[b.ID] = int32(k + 1)
	}
	var skills []learnermodel.Skill
	for _, id := range sortedKeys(s.SkillGraph.Skills) {
		sk := s.SkillGraph.Skills[id]
		a := map[string]learnermodel.Assessability{
			"verdict": learnermodel.AssessabilityVerdict, "rubric": learnermodel.AssessabilityRubric, "none": learnermodel.AssessabilityNone,
		}[sk.Assessability]
		if a == learnermodel.AssessabilityUnspecified {
			t.Fatalf("%s: skill %s has the assessability %q", module, id, sk.Assessability)
		}
		skills = append(skills, learnermodel.Skill{
			ID: full("skill", id), Assessability: a, Prerequisites: fulls("skill", sk.Prerequisites),
			Teaches: sk.Teaches, Tier: s.SkillTiers[id],
		})
	}

	r := release{ages: s.Manifest.Audience.AgeCategories, taskAge: map[string][]string{}, lessAge: map[string][]string{}}
	var tasks []learnermodel.Task
	var traps [][]ItemTraps
	byID := map[string]learnermodel.Task{}
	for _, id := range sortedKeys(s.Tasks) {
		st := s.Tasks[id]
		task := learnermodel.Task{
			ID: full("task", id), Band: learnermodel.BandRef{Place: place[st.Band], ID: st.Band}, Difficulty: st.Difficulty,
			Primary: fulls("skill", st.Skills.Primary), Supporting: fulls("skill", st.Skills.Supporting),
		}
		var its []ItemTraps
		for k, it := range st.Items {
			item := learnermodel.TaskItem{Place: int32(k + 1)}
			if it.Skills != nil {
				item.OwnSkills, item.Primary, item.Supporting = true, fulls("skill", it.Skills.Primary), fulls("skill", it.Skills.Supporting)
			}
			var caught []string
			for _, sl := range slices.Concat(it.AnswerSlots, it.Line) {
				if sl.Check == nil || sl.Check.Kind == "rubric" {
					continue
				}
				kind := map[string]learnermodel.CheckKind{
					"number": learnermodel.CheckKindNumber, "fraction": learnermodel.CheckKindFraction,
					"choice": learnermodel.CheckKindChoice, "predicate": learnermodel.CheckKindPredicate,
				}[sl.Check.Kind]
				if kind == learnermodel.CheckKindUnspecified {
					t.Fatalf("%s: task %s has a check of the kind %q", module, id, sl.Check.Kind)
				}
				slot := learnermodel.Slot{Check: kind}
				if kind == learnermodel.CheckKindChoice {
					slot.Options, slot.HasOptions = int32(len(slices.Concat(sl.Options, sl.Check.Options))), true
				}
				item.Slots = append(item.Slots, slot)
				for _, tv := range sl.Check.TrapValues {
					caught = append(caught, full("trap", tv.Trap))
				}
			}
			task.Items = append(task.Items, item)
			if len(caught) > 0 {
				its = append(its, ItemTraps{Place: item.Place, Traps: slices.Compact(sortedCopy(caught))})
			}
		}
		tasks = append(tasks, task)
		traps = append(traps, its)
		byID[id] = task
		r.taskAge[task.ID] = st.AgeCategories
	}
	g, err := learnermodel.NewGraph(module, bands, skills, nil, tasks)
	if err != nil {
		t.Fatalf("%s: %v", module, err)
	}
	r.module.Graph = g
	for k, task := range tasks {
		r.module.Pool = append(r.module.Pool, Candidate{Task: task, Traps: traps[k]})
	}
	for _, id := range sortedKeys(s.Lessons) {
		sl := s.Lessons[id]
		l := Lesson{ID: full("lesson", id), Skills: fulls("skill", sl.Skills)}
		for _, seg := range sl.Segments {
			if seg.Task == "" {
				continue
			}
			task, ok := byID[seg.Task]
			if !ok {
				t.Fatalf("%s: lesson %s names task %s, which the package has not", module, id, seg.Task)
			}
			l.Tasks = append(l.Tasks, task)
		}
		r.module.Lessons = append(r.module.Lessons, l)
		r.lessAge[l.ID] = sl.AgeCategories
	}
	return r
}

// forAge gives the module as api gives it a learner of an age category: the tasks and the lessons
// that suit them — their own age categories, or else the module's —, a lesson only when each of
// its tasks does too (architecture 3.7 of api).
func (r release) forAge(age string) Module {
	suits := func(ages []string) bool {
		if len(ages) == 0 {
			ages = r.ages
		}
		return slices.Contains(ages, age)
	}
	m := r.module
	m.Pool = slices.DeleteFunc(slices.Clone(m.Pool), func(c Candidate) bool { return !suits(r.taskAge[c.Task.ID]) })
	m.Lessons = slices.DeleteFunc(slices.Clone(m.Lessons), func(l Lesson) bool {
		return !suits(r.lessAge[l.ID]) || slices.ContainsFunc(l.Tasks, func(t learnermodel.Task) bool { return !suits(r.taskAge[t.ID]) })
	})
	return m
}

// The packages read through the chain of digests, and their translation holds what the examples
// of docs/architecture.md, sections 3.1–3.3, say of them: the skills, their tiers and
// prerequisites, the bands, the lessons and the anchor skills taught.
func TestReleasePackages(t *testing.T) {
	rel := readRelease(t)
	skillIDs := func(m Module) []string {
		var out []string
		for _, s := range graphOrder(m.Graph.Skills) {
			out = append(out, fmt.Sprintf("%s:%d", strings.TrimPrefix(s.ID, m.Graph.Module+"/skill/"), s.Tier))
		}
		return out
	}
	for _, tc := range []struct {
		module  string
		skills  []string
		bands   int
		tasks   int
		lessons int
	}{
		{releaseOlympiad, []string{"combinatorics-enumeration:1", "counting-gaps:1", "pigeonhole-basic:2"}, 3, 27, 9},
		{releaseGrade3, []string{"times-tables:1", "division-with-remainder:2", "fractions-of-a-whole:2", "comparing-fractions:3"}, 1, 12, 4},
		{releaseEssays, []string{"choosing-a-first-sentence:1", "reading-model-essays:1", "writing-a-first-sentence:2", "paragraphs:3",
			"endings:4", "improving-an-essay:4"}, 2, 8, 6},
	} {
		m := rel[tc.module].module
		if got := skillIDs(m); !slices.Equal(got, tc.skills) {
			t.Errorf("%s: skills %v, want %v", tc.module, got, tc.skills)
		}
		if len(m.Graph.Bands) != tc.bands || len(m.Pool) != tc.tasks || len(m.Lessons) != tc.lessons {
			t.Errorf("%s: %d bands, %d tasks, %d lessons; want %d, %d, %d", tc.module, len(m.Graph.Bands), len(m.Pool), len(m.Lessons), tc.bands, tc.tasks, tc.lessons)
		}
	}
	grade := rel[releaseGrade3].module
	i := slices.IndexFunc(grade.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == sid(releaseGrade3, "fractions-of-a-whole") })
	if !slices.Contains(grade.Graph.Skills[i].Teaches, fractionsAnchor) {
		t.Errorf("fractions-of-a-whole teaches %v, want %s among them", grade.Graph.Skills[i].Teaches, fractionsAnchor)
	}
	if child := rel[releaseEssays].forAge("under_13"); len(child.Lessons) != 5 || len(child.Pool) != 7 {
		t.Errorf("a child gets %d lessons and %d tasks of essays, want 5 and 7", len(child.Lessons), len(child.Pool))
	}
	// A task or a lesson with no age categories of its own takes the module's: school maths for
	// grade 3 is under_13 alone.
	if adult := rel[releaseGrade3].forAge("adult"); len(adult.Lessons) != 0 || len(adult.Pool) != 0 {
		t.Errorf("an adult gets %d lessons and %d tasks of school maths for grade 3, want none", len(adult.Lessons), len(adult.Pool))
	}
	l := lessonNamed(rel[releaseOlympiad].module, "counting-gaps-grades-3-4")
	var ids []string
	for _, task := range l.Tasks {
		ids = append(ids, task.ID)
	}
	if want := []string{tid(releaseOlympiad, "gaps-34-d1-1"), tid(releaseOlympiad, "gaps-34-d3-1"), tid(releaseOlympiad, "gaps-34-d5-1")}; !slices.Equal(ids, want) {
		t.Errorf("lesson counting-gaps-grades-3-4 holds %v, want %v", ids, want)
	}
}
