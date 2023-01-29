package util

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type dsCondition struct {
	GreaterValue      string
	GreaterEqualValue string
	LesserValue       string
	LesserEqualValue  string
	EqualValue        string
	NotEqualValue     string //暂不启用
}

type tableSplice struct {
	MonthDsArray []string
	WeekDsArray  []string
	DayDsArray   []string
	HourDsArray  []string
}

func AnalyticSql(sql string) string {
	reg := regexp.MustCompile(`\{(.+?)\}`)
	results := reg.FindAllString(sql, -1)
	if len(results) == 0 {
		return sql
	}
	for _, v := range results {
		tableName, tableAlias, selectSql, dsCondition, whereSql := parseSqlFragment(v)
		newSql := fmt.Sprintf("%sFROM (%s) %s", selectSql, parseTable(selectSql, tableName, tableAlias, whereSql, dsCondition), tableAlias)
		sql = strings.Replace(sql, v, newSql, -1)
	}
	return sql
}

func parseSqlFragment(sqlFragment string) (string, string, string, dsCondition, string) {
	sqlFragment = strings.ReplaceAll(sqlFragment[1:len(sqlFragment)-1], "\n", " ")
	sqlFragment = strings.ReplaceAll(sqlFragment, "select ", "SELECT ")
	sqlFragment = strings.ReplaceAll(sqlFragment, "from ", "FROM ")
	sqlFragment = strings.ReplaceAll(sqlFragment, "where ", "WHERE ")
	sqlFragment = strings.ReplaceAll(sqlFragment, "and ", "AND ")
	sqlFragment = strings.ReplaceAll(sqlFragment, "or ", "OR ")
	sqlFragment = strings.ReplaceAll(sqlFragment, "in ", "IN ")
	//在sqlFragment尾部添加空格保证正则命中
	if !strings.HasSuffix(sqlFragment, " ") {
		sqlFragment = sqlFragment + " "
	}

	fromIndex := strings.Index(sqlFragment, "FROM")
	whereIndex := strings.Index(sqlFragment, "WHERE")
	selectSql := sqlFragment[:fromIndex]
	var tableName, tableAlias string
	for _, v := range strings.Split(strings.ReplaceAll(sqlFragment[fromIndex:whereIndex], "FROM ", ""), " ") {
		if v != "" && v != "AS" {
			if tableName == "" {
				tableName = v
			} else {
				tableAlias = v
			}
		}
	}
	var dsCondition dsCondition
	whereSql := sqlFragment[whereIndex:]
	parseDsConditionWhereSql(&dsCondition, &whereSql)
	if whereSql == "WHERE " {
		whereSql = ""
	}
	return tableName, tableAlias, selectSql, dsCondition, whereSql
}

func parseTable(selectSql, tableName, tableAlias, whereSql string, dsCondition dsCondition) string {
	if tableAlias != "" {
		selectSql = strings.ReplaceAll(selectSql, fmt.Sprintf("%s.", tableAlias), "")
	}
	var result string
	if dsCondition.EqualValue != "" {
		result = fmt.Sprintf("(%sFROM %s_h WHERE ds=%s)", selectSql, tableName, dsCondition.EqualValue)
	} else {
		var startDate, endDate string
		if dsCondition.GreaterValue == "" && dsCondition.GreaterEqualValue == "" {
			//startDate = viper.GetString("bi.startDate")
		} else if dsCondition.GreaterValue != "" {
			dateTime, _ := time.Parse("20060102", dsCondition.GreaterValue)
			startDate = dateTime.AddDate(0, 0, 1).Format("20060102")
		} else {
			startDate = dsCondition.GreaterEqualValue
		}
		if dsCondition.LesserValue == "" && dsCondition.LesserEqualValue == "" {
			endDate = time.Now().Format("20060102")
		} else if dsCondition.LesserValue != "" {
			dateTime, _ := time.Parse("20060102", dsCondition.LesserValue)
			endDate = dateTime.AddDate(0, 0, -1).Format("20060102")
		} else {
			endDate = dsCondition.LesserEqualValue
		}
		tableSplice := calcTableSplice(startDate, endDate)
		var tableSpliceItems []string
		addTableSpliceItems(&tableSpliceItems, selectSql, tableName, tableAlias, whereSql, "_m", tableSplice.MonthDsArray)
		addTableSpliceItems(&tableSpliceItems, selectSql, tableName, tableAlias, whereSql, "_w", tableSplice.WeekDsArray)
		addTableSpliceItems(&tableSpliceItems, selectSql, tableName, tableAlias, whereSql, "_d", tableSplice.DayDsArray)
		addTableSpliceItems(&tableSpliceItems, selectSql, tableName, tableAlias, whereSql, "_h", tableSplice.HourDsArray)
		result = strings.Join(tableSpliceItems, " UNION ALL ")

	}
	return result
}

func parseDsConditionWhereSql(dsConditions *dsCondition, whereSql *string) {
	//匹配公式：
	//(AND|OR)(空格)(别名.)?(ds)(空格)?(>|>=|<|<=|=|!=)(空格)(带引号或者不带引号的数字)(空格)
	//或
	//(WHERE)(空格)(别名.)?(ds)(空格)?(>|>=|<|<=|=|!=)(空格)(带引号或者不带引号的数字)(空格)(AND|OR)(空格)
	reg := regexp.MustCompile(`((AND|OR)\s(\w*\.)?ds(\s?>\s?['"]?\d+['"]?|\s?>=\s?['"]?\d+['"]?|\s?<\s?['"]?\d+['"]?|\s?<=\s?['"]?\d+['"]?|\s?=\s?['"]?\d+['"]?|\s?!=\s?['"]?\d+['"]?)\s)|(WHERE\s(\w*\.)?ds(\s?>\s?['"]?\d+['"]?|\s?>=\s?['"]?\d+['"]?|\s?<\s?['"]?\d+['"]?|\s?<=\s?['"]?\d+['"]?|\s?=\s?['"]?\d+['"]?|\s?!=\s?['"]?\d+['"]?)\s((AND|OR)\s)?)`)
	condition := reg.FindString(*whereSql)
	if condition != "" {
		if strings.HasPrefix(condition, "WHERE ") {
			*whereSql = strings.ReplaceAll(*whereSql, condition, "WHERE ")
		} else {
			*whereSql = strings.ReplaceAll(*whereSql, condition, "")
		}
		conditionReg := regexp.MustCompile(`>\s?\d+|>=\s?\d+|<\s?\d+|<=\s?\d+|=\s?\d+|!=\s?\d+`)
		trimCondition := conditionReg.FindString(strings.ReplaceAll(strings.ReplaceAll(condition, "'", ""), "\"", ""))
		if trimCondition != "" {
			conditionValueReg := regexp.MustCompile(`\d+`)
			clearSpaceCondition := strings.ReplaceAll(trimCondition, " ", "")
			conditionValue := conditionValueReg.FindString(clearSpaceCondition)
			if conditionValue != "" {
				conditionExp := strings.ReplaceAll(clearSpaceCondition, conditionValue, "")
				switch conditionExp {
				case ">":
					dsConditions.GreaterValue = conditionValue
				case ">=":
					dsConditions.GreaterEqualValue = conditionValue
				case "<":
					dsConditions.LesserValue = conditionValue
				case "<=":
					dsConditions.LesserEqualValue = conditionValue
				case "=":
					dsConditions.EqualValue = conditionValue
				case "!=":
					dsConditions.NotEqualValue = conditionValue
				}
				parseDsConditionWhereSql(dsConditions, whereSql)
			}
		}
	}
}

func calcTableSplice(startDate, endDate string) tableSplice {
	var tableSplice tableSplice
	if startDate == endDate {
		tableSplice.HourDsArray = append(tableSplice.HourDsArray, startDate)
	} else {
		startDateTime, _ := time.Parse("20060102", startDate)
		endDateTime, _ := time.Parse("20060102", endDate)
		today := time.Now()
		//结束日期包含今日，带上今日实时数据
		if endDate == today.Format("20060102") {
			tableSplice.HourDsArray = append(tableSplice.HourDsArray, endDate)
			dateTime, _ := time.Parse("20060102", endDate)
			endDateTime = dateTime.AddDate(0, 0, -1)
		}
		//结束日期包含昨日，避免凌晨t+1未生成对应日表，还是读取小时表
		yesterday := today.AddDate(0, 0, -1)
		if endDateTime.Year() == yesterday.Year() &&
			endDateTime.Month() == yesterday.Month() &&
			endDateTime.Day() == yesterday.Day() {
			tableSplice.HourDsArray = append(tableSplice.HourDsArray, endDateTime.Format("20060102"))
			endDateTime = endDateTime.AddDate(0, 0, -1)
		}
		//剩余日期取最小组合
		calcMinimumCombinationTableSplice(startDateTime, endDateTime, &tableSplice)
	}
	return tableSplice
}

func calcMinimumCombinationTableSplice(startDateTime, endDateTime time.Time, tableSplice *tableSplice) {
	if startDateTime.Unix() <= endDateTime.Unix() {
		//判断月
		monthStartDateTime := endDateTime.AddDate(0, 0, 1-endDateTime.Day())
		if startDateTime.Unix() <= monthStartDateTime.Unix() {
			tableSplice.MonthDsArray = append(tableSplice.MonthDsArray, endDateTime.Format("20060102"))
			endDateTime = monthStartDateTime.AddDate(0, 0, -1)
		} else {
			//判断周
			var weekMinusDays int
			switch endDateTime.Weekday() {
			case time.Monday:
				weekMinusDays = 0
			case time.Tuesday:
				weekMinusDays = 1
			case time.Wednesday:
				weekMinusDays = 2
			case time.Thursday:
				weekMinusDays = 3
			case time.Friday:
				weekMinusDays = 4
			case time.Saturday:
				weekMinusDays = 5
			case time.Sunday:
				weekMinusDays = 6
			}
			weekStartDateTime := endDateTime.AddDate(0, 0, -weekMinusDays)
			if startDateTime.Unix() <= weekStartDateTime.Unix() {
				tableSplice.WeekDsArray = append(tableSplice.WeekDsArray, endDateTime.Format("20060102"))
				endDateTime = weekStartDateTime.AddDate(0, 0, -1)
			} else {
				//不足一周
				tableSplice.DayDsArray = append(tableSplice.DayDsArray, endDateTime.Format("20060102"))
				endDateTime = endDateTime.AddDate(0, 0, -1)
			}
		}
		calcMinimumCombinationTableSplice(startDateTime, endDateTime, tableSplice)
	}
}

func addTableSpliceItems(tableSpliceItems *[]string, selectSql, tableName, tableAlias, whereSql, tableSuffix string, dsArray []string) {
	for _, v := range dsArray {
		if tableAlias != "" {
			selectSql = strings.ReplaceAll(selectSql, fmt.Sprintf("%s.", tableAlias), "")
			whereSql = strings.ReplaceAll(whereSql, fmt.Sprintf("%s.", tableAlias), "")
		}
		*tableSpliceItems = append(*tableSpliceItems, fmt.Sprintf("%sFROM %s%s_%s %s", selectSql, tableName, tableSuffix, v, whereSql))
	}
}
