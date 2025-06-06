// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package partials

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func InfoModal() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div x-data=\"{open: false}\"><button class=\"btn btn-soft text-primary\" x-on:click=\"open = true\">Infos</button> <dialog class=\"modal modal-middle\" x-bind:open=\"open\"><div class=\"modal-box bg-base-100 rounded-lg shadow-lg\"><div class=\"prose text-base-content\"><p>Official repository:</p><p><a href=\"https://github.com/noahfraiture/nexzap\" class=\"link link-primary\" target=\"_blank\">Github repository</a></p><p>If you have any questions, please contact me:</p><p><a href=\"mailto:contact@nexzap.app\" class=\"link link-primary\">contact@nexzap.app</a></p></div><div class=\"modal-action mt-6\"><form method=\"dialog\"><button x-on:click=\"open = false\" class=\"btn btn-outline btn-primary\">Close</button></form></div></div><form method=\"dialog\" class=\"modal-backdrop\"><button x-on:click=\"open = false\"></button></form></dialog></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
