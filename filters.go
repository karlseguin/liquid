package liquid

import (
	"github.com/karlseguin/liquid/core"
	"github.com/karlseguin/liquid/filters"
)

func init() {
	core.RegisterFilter("capitalize", filters.CapitalizeFactory)
	core.RegisterFilter("downcase", filters.DowncaseFactory)
	core.RegisterFilter("upcase", filters.UpcaseFactory)
	core.RegisterFilter("first", filters.FirstFactory)
	core.RegisterFilter("last", filters.LastFactory)
	core.RegisterFilter("join", filters.JoinFactory)
	core.RegisterFilter("debug", filters.DebugFactory)
	core.RegisterFilter("plus", filters.PlusFactory)
	core.RegisterFilter("minus", filters.MinusFactory)
	core.RegisterFilter("size", filters.SizeFactory)
	core.RegisterFilter("times", filters.TimesFactory)
	core.RegisterFilter("divideby", filters.DivideByFactory)
	core.RegisterFilter("prepend", filters.PrependFactory)
	core.RegisterFilter("append", filters.AppendFactory)
	core.RegisterFilter("strip_newlines", filters.StripNewLinesFactory)
	core.RegisterFilter("strip_html", filters.StripHtmlFactory)
	core.RegisterFilter("replace", filters.ReplaceFactory)
	core.RegisterFilter("replace_first", filters.ReplaceFirstFactory)
	core.RegisterFilter("remove", filters.RemoveFactory)
	core.RegisterFilter("remove_first", filters.RemoveFirstFactory)
	core.RegisterFilter("newline_to_br", filters.NewLineToBrFactory)
	core.RegisterFilter("split", filters.SplitFactory)
	core.RegisterFilter("modulo", filters.ModuloFactory)
	core.RegisterFilter("truncate", filters.TruncateFactory)
	core.RegisterFilter("truncatewords", filters.TruncateWordsFactory)
	core.RegisterFilter("escape", filters.EscapeFactory)
	core.RegisterFilter("escape_once", filters.EscapeOnceFactory)
	core.RegisterFilter("sort", filters.SortFactory)
	core.RegisterFilter("default", filters.DefaultFactory)
	core.RegisterFilter("reverse", filters.ReverseFactory)
	core.RegisterFilter("date", filters.DateFactory)
}
