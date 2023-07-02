package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/suyuan32/simple-admin-common/orm/ent/mixins"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
	ent.Schema
}

// Fields of the Tag.
func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("Tag's name | 标签名称").
			Annotations(entsql.WithComments(true)),
		field.String("remark").Comment("The remark of tag | 标签的备注").
			Annotations(entsql.WithComments(true)),
	}
}

func (Tag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
		mixins.StatusMixin{},
	}
}

// Edges of the Tag.
func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
	}
}

func (Tag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
	}
}

func (Tag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fms_tags"}, // fms means Tag management service
	}
}
