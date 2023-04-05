package materialize

import (
	"fmt"
	"strings"
)

type TableLoadgen struct {
	Name  string
	Alias string
}

type CounterOptions struct {
	TickInterval   string
	ScaleFactor    float64
	MaxCardinality int
}

type AuctionOptions struct {
	TickInterval string
	ScaleFactor  float64
	Table        []TableLoadgen
}

type TPCHOptions struct {
	TickInterval string
	ScaleFactor  float64
	Table        []TableLoadgen
}

type SourceLoadgenBuilder struct {
	Source
	clusterName       string
	size              string
	loadGeneratorType string
	counterOptions    CounterOptions
	auctionOptions    AuctionOptions
	tpchOptions       TPCHOptions
}

func NewSourceLoadgenBuilder(sourceName, schemaName, databaseName string) *SourceLoadgenBuilder {
	return &SourceLoadgenBuilder{
		Source: Source{sourceName, schemaName, databaseName},
	}
}

func (b *SourceLoadgenBuilder) ClusterName(c string) *SourceLoadgenBuilder {
	b.clusterName = c
	return b
}

func (b *SourceLoadgenBuilder) Size(s string) *SourceLoadgenBuilder {
	b.size = s
	return b
}

func (b *SourceLoadgenBuilder) LoadGeneratorType(l string) *SourceLoadgenBuilder {
	b.loadGeneratorType = l
	return b
}

func (b *SourceLoadgenBuilder) CounterOptions(c CounterOptions) *SourceLoadgenBuilder {
	b.counterOptions = c
	return b
}

func (b *SourceLoadgenBuilder) AuctionOptions(a AuctionOptions) *SourceLoadgenBuilder {
	b.auctionOptions = a
	return b
}

func (b *SourceLoadgenBuilder) TPCHOptions(t TPCHOptions) *SourceLoadgenBuilder {
	b.tpchOptions = t
	return b
}

func (b *SourceLoadgenBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE SOURCE %s`, b.QualifiedName()))

	if b.clusterName != "" {
		q.WriteString(fmt.Sprintf(` IN CLUSTER %s`, QuoteIdentifier(b.clusterName)))
	}

	q.WriteString(fmt.Sprintf(` FROM LOAD GENERATOR %s`, b.loadGeneratorType))

	// Optional Parameters
	var p []string

	for _, t := range []string{b.counterOptions.TickInterval, b.auctionOptions.TickInterval, b.tpchOptions.TickInterval} {
		if t != "" {
			p = append(p, fmt.Sprintf(`TICK INTERVAL %s`, QuoteString(t)))
		}
	}

	for _, t := range []float64{b.counterOptions.ScaleFactor, b.auctionOptions.ScaleFactor, b.tpchOptions.ScaleFactor} {
		if t != 0 {
			p = append(p, fmt.Sprintf(`SCALE FACTOR %.2f`, t))
		}
	}

	if b.counterOptions.MaxCardinality != 0 {
		s := fmt.Sprintf(`MAX CARDINALITY %d`, b.counterOptions.MaxCardinality)
		p = append(p, s)
	}

	if len(p) != 0 {
		p := strings.Join(p[:], ", ")
		q.WriteString(fmt.Sprintf(` (%s)`, p))
	}

	// Table Mapping
	if b.loadGeneratorType == "COUNTER" {
		// Tables do not apply to COUNTER
	} else if len(b.auctionOptions.Table) > 0 || len(b.tpchOptions.Table) > 0 {

		var ot []TableLoadgen
		if len(b.auctionOptions.Table) > 0 {
			ot = b.auctionOptions.Table
		} else {
			ot = b.tpchOptions.Table
		}

		var tables []string
		for _, t := range ot {
			if t.Alias == "" {
				t.Alias = t.Name
			}
			s := fmt.Sprintf(`%s AS %s`, t.Name, t.Alias)
			tables = append(tables, s)
		}
		o := strings.Join(tables[:], ", ")
		q.WriteString(fmt.Sprintf(` FOR TABLES (%s)`, o))
	} else {
		q.WriteString(` FOR ALL TABLES`)
	}

	// Size
	if b.size != "" {
		q.WriteString(fmt.Sprintf(` WITH (SIZE = %s)`, QuoteString(b.size)))
	}

	q.WriteString(`;`)
	return q.String()
}

func (b *SourceLoadgenBuilder) Rename(newName string) string {
	n := QualifiedName(b.DatabaseName, b.SchemaName, newName)
	return fmt.Sprintf(`ALTER SOURCE %s RENAME TO %s;`, b.QualifiedName(), n)
}

func (b *SourceLoadgenBuilder) UpdateSize(newSize string) string {
	return fmt.Sprintf(`ALTER SOURCE %s SET (SIZE = %s);`, b.QualifiedName(), QuoteString(newSize))
}

func (b *SourceLoadgenBuilder) Drop() string {
	return fmt.Sprintf(`DROP SOURCE %s;`, b.QualifiedName())
}

func (b *SourceLoadgenBuilder) ReadId() string {
	return ReadSourceId(b.SourceName, b.SchemaName, b.DatabaseName)
}
