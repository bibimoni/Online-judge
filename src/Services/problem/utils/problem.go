package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"problem/models"
	"strconv"

	"github.com/antchfx/xmlquery"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
ParseProblem() - Read problem.xml

Current problems
- Only english statements are allowed

Workflow:
- Read problem.xml to get Name, ShortName, Tags and create a Problem
*/
func ParseProblemStruct(problemId uint64, xml *os.File) (models.Problem, error) {
	var problem models.Problem

	doc, err := xmlquery.Parse(xml)
	if err != nil {
		return problem, err
	}

	problem.ID = primitive.NewObjectID()
	problem.ProblemId = problemId
	problem.Name = xmlquery.Find(doc, "//problem/names/name[@language='english']/@value")[0].InnerText()
	problem.ShortName = xmlquery.Find(doc, "//problem/@short-name")[0].InnerText()
	tags := xmlquery.Find(doc, "//problem/tags/tag")
	for _, tag := range tags {
		problem.Tags = append(problem.Tags, tag.SelectAttr(("value")))
	}

	var str string = xmlquery.FindOne(doc, "//problem/judging/testset/test-count").InnerText()
	if val, err := strconv.Atoi(str); err != nil {
		return problem, err
	} else {
		problem.TestNum = uint64(val)
	}

	str = xmlquery.FindOne(doc, "//problem/judging/testset/time-limit").InnerText()
	if val, err := strconv.Atoi(str); err != nil {
		return problem, err
	} else {
		problem.TimeLimit = uint64(val)
	}

	str = xmlquery.FindOne(doc, "//problem/judging/testset/memory-limit").InnerText()
	if val, err := strconv.Atoi(str); err != nil {
		return problem, err
	} else {
		problem.MemoryLimit = uint64(val)
	}
	// ScoringMode

	// Parse problem's tests
	groups, err := xmlquery.QueryAll(doc, "//problem/judging/testset/groups/group")
	if err != nil {
		return problem, err
	}

	problem.ScoringMode = "OI"

	// Checking if ICPC type
	if len(groups) == 0 {
		tests, err := xmlquery.QueryAll(doc, "//problem/judging/testset/tests/test")
		if err != nil {
			return problem, err
		}
		isICPC := true
		for _, test := range tests {
			if test.SelectAttr("points") != "" {
				isICPC = false
				break
			}
		}
		if isICPC {
			problem.ScoringMode = "ICPC"
		}
	}

	if problem.ScoringMode == "OI" {
		// OI-style problem

		for _, group := range groups {
			var groupStruct models.ProblemGroup

			groupStruct.Name = group.SelectAttr("name")
			groupStruct.Scoring = group.SelectAttr("points-policy")

			// Dependencies
			if group.SelectElement("dependencies") == nil {
				continue
			}

			dependencies := group.SelectElement("dependencies")
			for _, dependency := range dependencies.SelectElements("dependency") {
				fmt.Printf("%s\n\n", dependency.SelectAttr("group"))

				groupStruct.Dependencies = append(groupStruct.Dependencies, dependency.SelectAttr("group"))
			}

			problem.TestGroups = append(problem.TestGroups, groupStruct)
		}

		tests, err := xmlquery.QueryAll(doc, "//problem/judging/testset/tests/test")
		if err != nil {
			return problem, err
		}
		for testIdx, test := range tests {
			groupName := test.SelectAttr("group")
			testPoint, _ := strconv.ParseFloat(test.SelectAttr("points"), 64)

			// handle cases where the doesn't belong to any group
			if groupName == "" {
				continue
			}

			for groupIdx := 0; groupIdx < len(problem.TestGroups); groupIdx++ {
				if problem.TestGroups[groupIdx].Name == groupName {
					problem.TestGroups[groupIdx].TestIndices = append(problem.TestGroups[groupIdx].TestIndices, testIdx)
					problem.TestGroups[groupIdx].MaxScore += testPoint
				}
			}
		}
	}

	return problem, nil
}

/*
SaveProblemToJson - Save Problem to a json file with specified filepath
*/
func SaveProblemToJson(problem models.Problem, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(problem); err != nil {
		return err
	}

	return nil
}
