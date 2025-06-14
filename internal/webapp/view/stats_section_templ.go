// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package view

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/deseteral/resistere/internal/metrics"
	"github.com/invopop/ctxi18n"
	"github.com/invopop/ctxi18n/i18n"
	"time"
)

func StatsSection(m *metrics.Registry) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<fieldset id=\"stats-section\" hx-get=\"/view/stats\" hx-trigger=\"load delay:5s\" hx-swap-oob=\"true\"><legend>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(i18n.T(ctx, "statistics.tab_title"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/webapp/view/stats_section.templ`, Line: 18, Col: 47}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</legend><div><style>\n\t\t\t\tme {\n\t\t\t\t\tdisplay: flex;\n\t\t\t\t\tflex-direction: column;\n\t\t\t\t}\n\t\t\t</style><div class=\"field-array\"><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(i18n.T(ctx, "statistics.power_production"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/webapp/view/stats_section.templ`, Line: 27, Col: 53}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, ":</div><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%.2f", m.LatestFrame.PowerProductionWatts/1000.0))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/webapp/view/stats_section.templ`, Line: 28, Col: 75}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, " kW</div><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(i18n.T(ctx, "statistics.power_consumption"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/webapp/view/stats_section.templ`, Line: 29, Col: 54}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, ":</div><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%.2f", m.LatestFrame.PowerConsumptionWatts/1000.0))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/webapp/view/stats_section.templ`, Line: 30, Col: 76}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, " kW</div></div><div><style>\n\t\t\t\t\tme {\n\t\t\t\t\t\tdisplay: flex;\n\t\t\t\t\t\tflex-direction: row;\n\t\t\t\t\t\talign-items: center;\n\t\t\t\t\t\talign-self: flex-end;\n\n\t\t\t\t\t\tmargin-top: 8px;\n\n\t\t\t\t\t\tfont-size: 13px;\n\t\t\t\t\t\tfont-variant-numeric: tabular-nums;\n\n\t\t\t\t\t\ti {\n\t\t\t\t\t\t\tmargin-right: 4px;\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t</style><i class=\"ph-bold ph-arrows-clockwise\"></i> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(i18n.T(ctx, "statistics.last_updated"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/webapp/view/stats_section.templ`, Line: 51, Col: 44}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, ":&nbsp; <time datetime=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(m.LatestFrame.Timestamp.Format(time.RFC3339))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/webapp/view/stats_section.templ`, Line: 52, Col: 65}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "\"></time><script>\n\t\t\t\t\t(function() {\n\t\t\t\t\t\tconst rtf = new Intl.RelativeTimeFormat(")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var9, templ_7745c5c3_Err := templruntime.ScriptContentOutsideStringLiteral(ctxi18n.Locale(ctx).Code().Base())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/webapp/view/stats_section.templ`, Line: 55, Col: 82}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var9)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, ", { style: \"long\" });\n\n\t\t\t\t\t\tconst el = me(\"time\");\n\t\t\t\t\t\tconst time = new Date(el.getAttribute(\"datetime\"));\n\n\t\t\t\t\t\tfunction updateTimer() {\n\t\t\t\t\t\t\tconst diffInSeconds = -Math.floor((Date.now() - time) / 1000);\n\t\t\t\t\t\t\tel.innerHTML = rtf.format(diffInSeconds, \"second\");\n\t\t\t\t\t\t}\n\n\t\t\t\t\t\tif (window.statsTimerIntervalId) {\n\t\t\t\t\t\t\tclearInterval(window.statsTimerIntervalId);\n\t\t\t\t\t\t}\n\n\t\t\t\t\t\tupdateTimer();\n\t\t\t\t\t\twindow.statsTimerIntervalId = setInterval(updateTimer, 500);\n\t\t\t\t\t})();\n\t\t\t\t</script></div></div></fieldset>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
