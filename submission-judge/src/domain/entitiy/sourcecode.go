package domain

import (
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SourceCode struct {
	Id           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	LanguageId   string        `json:"language,omitempty" bson:"language"`
	CreatedAt    time.Time     `json:"create_at" bson:"create_at"`
	FileSize     memory.Memory `json:"file_size,omitempty" bson:"file_size"`
	SourceCode   string        `json:"source_code,omitempty" bson:"source_code"`
	SubmissionId bson.ObjectID `json:"submission_id,omitempty" bson:"submission_id"`
}
