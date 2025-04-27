/*
package reflect/helper 从包 reflect 的 Value 和 Type 分离出不同类型的反射接口转发。

TryTypeTo, TypeTo 从 [reflect.Type] 尝试构建指定类型的包装 Type，该反射接口转发仅包含适用于特定类型的反射方法集

	if st := new(SliceType); TryTypeTo(tp, &st) {
		fmt.Println(st.String())
	}

Is 检查一个包装 Type 是哪种类型
TypeFor 从 T 的类型构造一个包装 Type

	if t := TypeFor[[]int](); Is[*SliceType](t) {
		fmt.Println(t.String())
	}

TypeOf 从 v 中提取反射信息，并包装为其类型特定的反射接口转发。

	if n := TypeOf(nil); Is[Nil](n) {
		fmt.Println(n.String())
	}

TypeWrap 包装一个 [reflect.Type]，根据 Kind 可以显式转换为特定的反射接口转发。

	if t := TypeWrap(tp); Is[*SliceType](t) {
		fmt.Println(t.String())
	}

TypeCommon 包含 [reflect.Type] 的一些通用方法

	if !IsNilType(t) {
		t := TypeCom(t)
		...
	}

or

	if t := TypeCom(t); t != nil {
		...
	}

TypeProperty 附加属性（可选）

	if !IsNilType(t) {
		prop := PropFor(TypeCom(t))
		...
	}
*/
package helper
