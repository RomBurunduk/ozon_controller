package fixtures

import (
	"pvz_controller/internal/model"
)

type ArticleBuilder struct {
	instance *model.Pickups
}

func Article() *ArticleBuilder {
	return &ArticleBuilder{instance: &model.Pickups{}}
}

func (b *ArticleBuilder) Address(v string) *ArticleBuilder {
	b.instance.Address = v
	return b
}

func (b *ArticleBuilder) Name(v string) *ArticleBuilder {
	b.instance.Name = v
	return b
}

func (b *ArticleBuilder) Contact(v string) *ArticleBuilder {
	b.instance.Contact = v
	return b
}

func (b *ArticleBuilder) P() *model.Pickups {
	return b.instance
}

func (b *ArticleBuilder) V() model.Pickups {
	return *b.instance
}

func (b *ArticleBuilder) Valid() *ArticleBuilder {
	return Article().Name("some").Address("some").Contact("some")
}
