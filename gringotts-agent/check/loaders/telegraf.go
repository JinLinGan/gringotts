package loaders

import (
	"time"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"github.com/jinlingan/gringotts/gringotts-agent/check"
	"github.com/jinlingan/gringotts/gringotts-agent/check/checks"
	"github.com/jinlingan/gringotts/gringotts-agent/model"
	"github.com/jinlingan/gringotts/gringotts-agent/telegraf"
	"github.com/jinlingan/gringotts/gringotts-agent/telegraf/inputs"
	"github.com/jinlingan/gringotts/pkg/log"
	"github.com/pkg/errors"
)

type TelegrafLoader struct{}

func (t *TelegrafLoader) Loade(jobConfig *model.JobConfig) (check.Check, error) {
	//panic("implement me")
	creator, ok := inputs.Inputs[jobConfig.RunnerModule]
	if !ok {
		return nil, errors.Errorf("未找到 Telegraf 采集插件 %q", jobConfig.RunnerModule)
	}
	input := creator()

	//Telegraf 原生在这里会插入解析函数，目前我们暂不实现

	table, err := toml.Parse([]byte(jobConfig.Config))

	// 处理 Interval
	if node, ok := table.Fields["interval"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				dur, err := time.ParseDuration(str.Value)
				if err != nil {
					log.Errorf("任务 %q 的配置文件中发现了 interval 字段值为 %q，解析失败，忽略该字段: %v", jobConfig.JobID, str, err)
				}
				newInterval := int32(dur.Seconds())
				log.Infof("任务 %q 的配置文件中发现了 interval 字段值为 %q 解析为 %d 秒，使用该值替换任务中的值", jobConfig.JobID, str, newInterval)
				//TODO： 强制类型转换，可能会有问题
				jobConfig.Interval = newInterval
			}
		}
	}

	delete(table.Fields, "interval")

	newCheck, err := t.GetTelegrafChecker(jobConfig.JobID, input, table)

	if err != nil {
		return newCheck, err
	}

	if err := toml.UnmarshalTable(table, input); err != nil {
		return nil, err
	}

	return newCheck, nil

}

func (t *TelegrafLoader) GetTelegrafChecker(jobID string, input telegraf.Input, tbl *ast.Table) (*checks.TelegrafCheck, error) {
	cp := &checks.TelegrafCheck{
		JobID:         jobID,
		TelegrafInput: input,
	}
	if node, ok := tbl.Fields["name_prefix"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				cp.MeasurementPrefix = str.Value
			}
		}
	}

	if node, ok := tbl.Fields["name_suffix"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				cp.MeasurementSuffix = str.Value
			}
		}
	}

	if node, ok := tbl.Fields["name_override"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				cp.NameOverride = str.Value
			}
		}
	}

	//if node, ok := tbl.Fields["alias"]; ok {
	//	if kv, ok := node.(*ast.KeyValue); ok {
	//		if str, ok := kv.Value.(*ast.String); ok {
	//			cp.Alias = str.Value
	//		}
	//	}
	//}

	cp.Tags = make(map[string]string)
	if node, ok := tbl.Fields["tags"]; ok {
		if subtbl, ok := node.(*ast.Table); ok {
			if err := toml.UnmarshalTable(subtbl, cp.Tags); err != nil {
				log.Errorf("Tag，解析失败，忽略该字段: %v")
			}
		}
	}

	delete(tbl.Fields, "name_prefix")
	delete(tbl.Fields, "name_suffix")
	delete(tbl.Fields, "name_override")
	delete(tbl.Fields, "alias")
	delete(tbl.Fields, "tags")
	var err error
	cp.Filter, err = buildFilter(tbl)
	if err != nil {
		return cp, err
	}
	return cp, nil
}

func buildFilter(tbl *ast.Table) (telegraf.Filter, error) {
	f := telegraf.Filter{}

	if node, ok := tbl.Fields["namepass"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						f.NamePass = append(f.NamePass, str.Value)
					}
				}
			}
		}
	}

	if node, ok := tbl.Fields["namedrop"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						f.NameDrop = append(f.NameDrop, str.Value)
					}
				}
			}
		}
	}

	fields := []string{"pass", "fieldpass"}
	for _, field := range fields {
		if node, ok := tbl.Fields[field]; ok {
			if kv, ok := node.(*ast.KeyValue); ok {
				if ary, ok := kv.Value.(*ast.Array); ok {
					for _, elem := range ary.Value {
						if str, ok := elem.(*ast.String); ok {
							f.FieldPass = append(f.FieldPass, str.Value)
						}
					}
				}
			}
		}
	}

	fields = []string{"drop", "fielddrop"}
	for _, field := range fields {
		if node, ok := tbl.Fields[field]; ok {
			if kv, ok := node.(*ast.KeyValue); ok {
				if ary, ok := kv.Value.(*ast.Array); ok {
					for _, elem := range ary.Value {
						if str, ok := elem.(*ast.String); ok {
							f.FieldDrop = append(f.FieldDrop, str.Value)
						}
					}
				}
			}
		}
	}

	if node, ok := tbl.Fields["tagpass"]; ok {
		if subtbl, ok := node.(*ast.Table); ok {
			for name, val := range subtbl.Fields {
				if kv, ok := val.(*ast.KeyValue); ok {
					tagfilter := &telegraf.TagFilter{Name: name}
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								tagfilter.Filter = append(tagfilter.Filter, str.Value)
							}
						}
					}
					f.TagPass = append(f.TagPass, *tagfilter)
				}
			}
		}
	}

	if node, ok := tbl.Fields["tagdrop"]; ok {
		if subtbl, ok := node.(*ast.Table); ok {
			for name, val := range subtbl.Fields {
				if kv, ok := val.(*ast.KeyValue); ok {
					tagfilter := &telegraf.TagFilter{Name: name}
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								tagfilter.Filter = append(tagfilter.Filter, str.Value)
							}
						}
					}
					f.TagDrop = append(f.TagDrop, *tagfilter)
				}
			}
		}
	}

	if node, ok := tbl.Fields["tagexclude"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						f.TagExclude = append(f.TagExclude, str.Value)
					}
				}
			}
		}
	}

	if node, ok := tbl.Fields["taginclude"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						f.TagInclude = append(f.TagInclude, str.Value)
					}
				}
			}
		}
	}
	if err := f.Compile(); err != nil {
		return f, err
	}

	delete(tbl.Fields, "namedrop")
	delete(tbl.Fields, "namepass")
	delete(tbl.Fields, "fielddrop")
	delete(tbl.Fields, "fieldpass")
	delete(tbl.Fields, "drop")
	delete(tbl.Fields, "pass")
	delete(tbl.Fields, "tagdrop")
	delete(tbl.Fields, "tagpass")
	delete(tbl.Fields, "tagexclude")
	delete(tbl.Fields, "taginclude")
	return f, nil
}

func newTelegrafLoader() check.Loader {
	return &TelegrafLoader{}
}

func init() {
	check.RegisterLoader(model.RunnerTypeTelegraf, newTelegrafLoader())
}
