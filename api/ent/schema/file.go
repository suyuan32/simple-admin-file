package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/suyuan32/simple-admin-common/orm/ent/mixins"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("uuid").Comment("File's UUID | 文件的UUID"),
		field.String("name").Comment("File's name | 文件名称"),
		field.Uint8("file_type").Comment("File's type | 文件类型"),
		field.Uint64("size").Comment("File's size"),
		field.String("path").Comment("File's path"),
		field.String("user_uuid").Comment("User's UUID | 用户的 UUID"),
		field.String("md5").Comment("The md5 of the file | 文件的 md5"),
	}
}

func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.BaseMixin{},
		mixins.StatusMixin{},
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return nil
}

func (File) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fms_files"}, // fms means file management service
	}
}
