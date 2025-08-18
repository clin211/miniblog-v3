package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
)

// validateRequired 验证必填字段
func validateRequired(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return fmt.Errorf("字段不能为空")
	}
	return nil
}

// validateEmail 验证邮箱格式
func validateEmail(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil // 空值跳过验证，除非同时有 required 标签
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("邮箱字段必须是字符串类型")
	}

	email := field.String()
	if email == "" {
		return nil // 空字符串跳过验证
	}

	if !govalidator.IsEmail(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	return nil
}

// validateAlphanum 验证字母数字
func validateAlphanum(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	// 允许下划线，因为用户名通常包含下划线
	if !govalidator.IsAlphanumeric(value) && !strings.Contains(value, "_") {
		return fmt.Errorf("只能包含字母、数字和下划线")
	}

	return nil
}

// validateNumeric 验证数字
func validateNumeric(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	if !govalidator.IsNumeric(value) {
		return fmt.Errorf("只能包含数字")
	}

	return nil
}

// validateAlpha 验证字母
func validateAlpha(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	if !govalidator.IsAlpha(value) {
		return fmt.Errorf("只能包含字母")
	}

	return nil
}

// validateASCII 验证ASCII字符
func validateASCII(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	if !govalidator.IsASCII(value) {
		return fmt.Errorf("只能包含ASCII字符")
	}

	return nil
}

// validateURL 验证URL格式
func validateURL(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	if !govalidator.IsURL(value) {
		return fmt.Errorf("URL格式不正确")
	}

	return nil
}

// validateIPv4 验证IPv4地址
func validateIPv4(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	if !govalidator.IsIPv4(value) {
		return fmt.Errorf("IPv4地址格式不正确")
	}

	return nil
}

// validateIPv6 验证IPv6地址
func validateIPv6(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	if !govalidator.IsIPv6(value) {
		return fmt.Errorf("IPv6地址格式不正确")
	}

	return nil
}

// validateUUID 验证UUID格式
func validateUUID(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	if !govalidator.IsUUID(value) {
		return fmt.Errorf("UUID格式不正确")
	}

	return nil
}

// validateJSON 验证JSON格式
func validateJSON(field reflect.Value, _ reflect.StructField) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	if !govalidator.IsJSON(value) {
		return fmt.Errorf("JSON格式不正确")
	}

	return nil
}

// validateLength 验证字符串长度
func validateLength(field reflect.Value, _ reflect.StructField, params []string) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()
	length := len(value)

	if len(params) == 0 {
		return fmt.Errorf("length规则需要参数")
	}

	if len(params) == 1 {
		// 固定长度
		expectedLength, err := strconv.Atoi(params[0])
		if err != nil {
			return fmt.Errorf("长度参数必须是数字")
		}

		if length != expectedLength {
			return fmt.Errorf("长度必须为 %d 个字符", expectedLength)
		}
	} else if len(params) == 2 {
		// 范围长度
		minLength, err := strconv.Atoi(params[0])
		if err != nil {
			return fmt.Errorf("最小长度参数必须是数字")
		}

		maxLength, err := strconv.Atoi(params[1])
		if err != nil {
			return fmt.Errorf("最大长度参数必须是数字")
		}

		if length < minLength || length > maxLength {
			return fmt.Errorf("长度必须在 %d 到 %d 个字符之间", minLength, maxLength)
		}
	}

	return nil
}

// validateRange 验证数值范围
func validateRange(field reflect.Value, _ reflect.StructField, params []string) error {
	if field.IsZero() {
		return nil
	}

	if len(params) != 2 {
		return fmt.Errorf("range规则需要两个参数")
	}

	min, err := strconv.ParseFloat(params[0], 64)
	if err != nil {
		return fmt.Errorf("最小值参数必须是数字")
	}

	max, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return fmt.Errorf("最大值参数必须是数字")
	}

	var value float64

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = float64(field.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = float64(field.Uint())
	case reflect.Float32, reflect.Float64:
		value = field.Float()
	case reflect.String:
		val, err := strconv.ParseFloat(field.String(), 64)
		if err != nil {
			return fmt.Errorf("字段值必须是数字")
		}
		value = val
	default:
		return fmt.Errorf("range规则只能用于数值类型")
	}

	// 检查是否在范围内（包含边界值）
	if value < min || value > max {
		return fmt.Errorf("值必须在 %v 到 %v 之间", min, max)
	}

	return nil
}

// validateMatches 验证正则表达式匹配
func validateMatches(field reflect.Value, _ reflect.StructField, params []string) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	if len(params) == 0 {
		return fmt.Errorf("matches规则需要正则表达式参数")
	}

	pattern := params[0]
	value := field.String()

	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return fmt.Errorf("正则表达式格式错误: %v", err)
	}

	if !matched {
		return fmt.Errorf("值不符合正则表达式模式: %s", pattern)
	}

	return nil
}

// validateIn 验证值是否在指定列表中
func validateIn(field reflect.Value, _ reflect.StructField, params []string) error {
	if field.IsZero() {
		return nil
	}

	if len(params) == 0 {
		return fmt.Errorf("in规则需要至少一个参数")
	}

	for _, param := range params {
		// 尝试类型转换
		switch field.Kind() {
		case reflect.String:
			if field.String() == param {
				return nil
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(param, 10, 64); err == nil && field.Int() == intVal {
				return nil
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintVal, err := strconv.ParseUint(param, 10, 64); err == nil && field.Uint() == uintVal {
				return nil
			}
		case reflect.Float32, reflect.Float64:
			if floatVal, err := strconv.ParseFloat(param, 64); err == nil && field.Float() == floatVal {
				return nil
			}
		default:
			return fmt.Errorf("in规则不支持该字段类型")
		}
	}

	return fmt.Errorf("值必须是以下之一: %s", strings.Join(params, ", "))
}

// validateWithGovalidator 使用 govalidator 的内置验证器
func validateWithGovalidator(field reflect.Value, _ reflect.StructField, rule string) error {
	if field.IsZero() {
		return nil
	}

	if field.Kind() != reflect.String {
		return fmt.Errorf("字段必须是字符串类型")
	}

	value := field.String()

	// 使用 govalidator 的内置函数
	switch rule {
	case "base64":
		if !govalidator.IsBase64(value) {
			return fmt.Errorf("必须是有效的Base64编码")
		}
	case "datauri":
		if !govalidator.IsDataURI(value) {
			return fmt.Errorf("必须是有效的Data URI")
		}
	case "port":
		if !govalidator.IsPort(value) {
			return fmt.Errorf("必须是有效的端口号")
		}
	case "dns":
		if !govalidator.IsDNSName(value) {
			return fmt.Errorf("必须是有效的DNS名称")
		}
	case "host":
		if !govalidator.IsHost(value) {
			return fmt.Errorf("必须是有效的主机名")
		}
	case "mac":
		if !govalidator.IsMAC(value) {
			return fmt.Errorf("必须是有效的MAC地址")
		}
	case "latitude":
		if !govalidator.IsLatitude(value) {
			return fmt.Errorf("必须是有效的纬度")
		}
	case "longitude":
		if !govalidator.IsLongitude(value) {
			return fmt.Errorf("必须是有效的经度")
		}
	case "ssn":
		if !govalidator.IsSSN(value) {
			return fmt.Errorf("必须是有效的SSN")
		}
	case "semver":
		if !govalidator.IsSemver(value) {
			return fmt.Errorf("必须是有效的语义化版本")
		}
	case "rfc3339":
		if !govalidator.IsRFC3339(value) {
			return fmt.Errorf("必须是有效的RFC3339时间格式")
		}
	case "ulid":
		if !govalidator.IsULID(value) {
			return fmt.Errorf("必须是有效的ULID")
		}
	case "yyyymmdd":
		// 使用正则表达式验证YYYYMMDD格式
		matched, _ := regexp.MatchString(`^\d{4}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])$`, value)
		if !matched {
			return fmt.Errorf("必须是有效的YYYYMMDD日期格式")
		}
	default:
		return fmt.Errorf("未知的验证规则: %s", rule)
	}

	return nil
}
