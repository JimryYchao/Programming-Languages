package gostd

import (
	. "reflect"
	"testing"
	"time"
)

//? go test -v -run=^$
/*
! reflect.Type 是 Go 类型的反射；并非所有方法都适用于所有类型，需要使用 Kind 确定其种类。
	Name 返回类型名称；匿名, 未定义类型返回 ""
	Align 返回该值的类型的对齐
	Size 返回存储给定类型的值时需要的字节大小
	String 返回类型的字符串形式（短名称）
	Kind 返回该类型的种类
	PkgPath 返回命名类型的包路径；T：Defined-Type
	FieldAlign 返回该值作为结构体字段时的对齐；T：Struct.Field

	Bits 返回类型的位大小；T：sized or unsized Int, Uint, Float, Complex 类型，则 panic
	ChanDir 返回通道类型的方向；T：Chan
	Elem 返回类型的元素类型；T：Array, Chan, Map, Pointer, Slice
	Len 返回数组类型的长度；T：Array
	Key 返回映射的键类型；T：Map

	Implements 报告该类型是否实现接口 u
	AssignableTo 报告该类型的值是否可分配给类型 u
	ConvertibleTo 报告该类型的值是否可以转换为类型 u；但仍可能 panic，例如 []T to *[N]T，len([]T) < N
	Comparable 报告该类型的值是否可比较，但仍可能 panic；例如接口的值可比较，但它们的动态类型不可比

	NumMethod 返回类型方法集的数目
	Method 返回其方法集中第 i 个方法，i ∈ [0, type.NumMethod())
	MethodByName 返回指定名称的方法

T：Function
	IsVariadic 报告函数类型是否是可变参 ... 函数；如是，t.In(t.NumIn()-1) 返回可变参的隐式类型 T[]
	NumIn 返回函数类型的参数数目
	In 返回函数类型的第 i 个输入参数类型；i ∈ [0, type.NumIn())
	NumOut 返回函数类型的输出参数数目
	Out 返回函数类型的第 i 个输出参数类型
T：Struct
	NumField 返回结构类型的字段数目
	Field 返回结构类型的第 i 个字段 StructField
		! field.IsExported 报告是否可导出
		! StructTag 字段的标记字符串
	! VisibleFields 返回结构的所有可见字段
	FieldByIndex 返回结构类型索引序列对应的嵌套字段
	FieldByName 返回结构类型指定名称的字段；如果字段是提升字段，则 StructField 的 Offset 是嵌入字段的结构体中的偏移量
	FieldByNameFunc 返回结构字段，其名称满足 match 匹配函数；优先考虑结构本身的字段，然后是嵌入字段
! ArrayOf 返回具有给定长度和元素类型的数组类型
! ChanOf 返回给定方向和元素类型的通道类型
! FuncOf 返回具有给定参数和结果类型的函数类型
! MapOf 返回具有给定键和元素类型的映射类型
! PointerTo 返回带有元素 t 的指针类型
! SliceOf 返回元素类型为 t 的切片类型
! StructOf 返回包含字段的结构类型
! TypeFor 返回表示类型参数 T 的 Type
! TypeOf 返回表示 i 的动态类型的反射 Type
*/

func TestTypeOf(t *testing.T) {
	tp := TypeOf(myInt(0))
	logfln(`the typeOf(myInt(0)): 
	Name: %s
	Align: %d
	Size: %d
	String: %s
	Kind: %s
	PkgPath: %s`, tp.Name(), tp.Align(), tp.Size(), tp.String(), tp.Kind(), tp.PkgPath())

	array := ArrayOf(time.Now().Nanosecond()%256, TypeFor[myInt]())
	logfln("the array type is [%d]%s", array.Len(), array.Elem())

	ch := ChanOf(RecvDir, TypeFor[bool]())
	logfln("the chan type is %s, ChanDir: %s, Elem: %s", ch.String(), ch.ChanDir(), ch.Elem())

	f := FuncOf([]Type{ChanOf(SendDir, TypeFor[bool]()), TypeFor[[]int]()}, []Type{TypeFor[string](), TypeFor[error]()}, true)
	logfln("%s : In(%d), Out(%d)", f.String(), f.NumIn(), f.NumOut())
	logfln("IsVariadic: %t, VariadicType: %s", f.IsVariadic(), f.In(f.NumIn()-1).String())

	m := MapOf(TypeFor[string](), TypeFor[any]())
	logfln("%s : <Key(%s), Elem(%s)>", m.String(), m.Key(), m.Elem())

	p := PointerTo(TypeOf(TestTypeOf))
	logfln("%s: Elem(%s)", p.String(), p.Elem())

	slice := SliceOf(TypeFor[byte]())
	logfln("%s: Elem(%s)", slice.String(), slice.Elem())
}

func TestStructOf(t *testing.T) {
	var fieldInfo func(Type)
	fieldInfo = func(tp Type) {
		n := tp.NumField()
		for i := range n {
			f := tp.Field(i)
			output := f.Name
			if v, ok := f.Tag.Lookup("k2"); ok {
				output += ": k2=" + v
			} else if f.Tag != "" {
				output += ": tag=" + string(f.Tag)
			}
			log(output)
			if f.Anonymous {
				if f.Type.Kind() == Struct {
					fieldInfo(f.Type)
				}
			}
		}
	}
	methodInfo := func(tp Type) {
		n := tp.NumMethod()
		logfln("methods in %s", tp.Name())
		for i := range n {
			log(tp.Method(i).Name)
		}
	}

	fieldInfo(TypeFor[myS]())
	methodInfo(TypeFor[myS]())
	methodInfo(TypeFor[*myS]())
}

type myS struct {
	V       myInt `tag of V`
	nestedS `k1:"v1" k2:"v2"`
}

func (myS) FunA() bool  { return true }
func (myS) FunB(...any) {}
func (*myS) FunC()      {}

type nestedS struct {
	V myInt
	myInt
}

func (nestedS) FunC()  {}
func (nestedS) FunD()  {}
func (*nestedS) FunE() {}

/*
! reflect.Value 是 Go 值（interface）的反射；零值表示没有值。IsValid() 返回 false；Kind() 返回 Invalid
	√ Type 返回值的 Type
	√ Cap 获取或设置 v 的容量: Array, Chan, Slice, *Array
	√ Len 返回 v 的长度: Array, Chan, Map, Slice, String, *Array
	√ Clear 清除切片或映射的内容：Map, Slice
	√ Index 返回 v 的第 i 个元素: Array, Slice, String
	√ Elem 返回接口类型或指针类型包含的值: Interface, Pointer
type getter, setter:
	 Addr, CanAddr: S.Field, Slice[i]
	√ Set, CanSet, SetZero
	UnsafeAddr 返回一个指向 v 数据的指针，最好使用 uintptr(Value.Addr().UnsafePointer())
	Bool, CanBool, SetBool
	Bytes, SetBytes: []byte, [L]byte
	√String, SetString
	Call, CallSlice: Func
	Complex, CanComplex, OverflowComplex, SetComplex
	Convert, CanConvert
	Float, CanFloat, OverflowFloat, SetFloat
	Int, CanInt, OverflowInt, SetInt
	Uint, CanUint, OverflowUint, SetUint
	Interface, CanInterface
	√ SetIterKey 将 iter 当前 map 元素的键赋值给 v；等价于 v.Set(iter.Key())
	√ SetIterValue 将 iter 当前 map 元素的值赋值给 v；等价于 v.Set(iter.Value())
checkers:
	Comparable, Equal
	IsNil : Chan, Slice, Pointer, Func, Interface, Map
	IsZero, IsValid, Kind
slice:
	SetLen, SetCap 设置切片的长度和容量
	Grow 容量不足时，尽可能增加切片的容量, 以满足至少有 n+1 的富余
	Slice 返回 s[i:j]; T: Array, Slice, String, *Array
	Slice3 返回 s[i:j:c]; T: Array, Slice, *Array
map:
	√ MapIndex, MapKeys 返回 map[key] 或 []keys
	√ MapRange 返回映射 v 的范围迭代器 *MapIter
	√ SetMapIndex 设置 map[key] = elem; 当 elem 为 zero Value 时删除键 key
chan:
	Close 关闭通道; Chan, <-Chan
	Recv 阻塞并接受一个值; Chan,<-Chan
	TryRecv 不会阻塞
	Send 发送一个值; Chan, Chan<-
	TrySend 不会阻塞
method:
	Method, MethodByName, NumMethod
struct:
	Field, FieldByIndex, FieldByIndexErr, FieldByName, FieldByNameFunc 返回结构字段
	NumField 返回结构字段数目
pointer:
	SetPointer 将 unsafe.Pointer 值 v 设置为 x；T：UnsafePointer
	Pointer 返回 v 的 uintptr: Chan, Func, Map, Pointer, Slice, UnsafePointer
	UnsafePointer 返回 v 的 unsafe.Pointer: Chan, Func, Map, Pointer, Slice, UnsafePointer
! Append, AppendSlice 将值 x... 或切片 t 追加到切片 s 上
! Indirect 间接返回 v 指向的值；v 非指针则返回 v
! MakeChan 创建一个具有指定类型和缓冲区大小的通道
! MakeFunc 创建返回包装函数 fn(args)results 的给定 Type 的新函数, Value.Call 方法调用该类型化函数
! MakeMap, MakeMapWithSize 创建一个指定类型的新 map
! MakeSlice 为指定的切片类型、长度和容量创建一个新的零初始化切片值
! New 返回指定类型的零值的指针 Value，即 PointerTo(type)
! NewAt 返回指定类型值的指针
! Select 执行由 cases 列表描述的选择操作；它会阻塞直到至少有一个 case 可以继续，做出统一的伪随机选择，然后执行该情况。它返回所选情况的索引；
! ValueOf 返回一个新的 Value；它被初始化为存储在接口 i 中的具体值；ValueOf(nil) 返回零值，代表空
! Zero 返回一个指定类型的零值而非 ValueOf(nil)
*/

// func TestSliceType(t *testing.T) { // 从 Value 包装一个 non-panic 的 Slice
// 	sp0 := helper.SliceFor[int]()              // is []int
// 	log(sp0.Elem().String())                   // int
// 	sp1, _ := helper.SliceOf(TypeFor[[]int]()) // is [][]int
// 	log(sp1.Elem().String())                   // []int

// 	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
// 	if s, err := helper.SliceFrom(&ints); err == nil {
// 		log(s.Cap(), s.Len(), s.ElemString(), s) // 10, 10, int, [1, 2, 3, 4, 5, 6, 7, 8, 9, 0]
// 		s.SetIndex(5, 10)
// 		log(s)                            // String: [1, 2, 3, 4, 5, 10, 7, 8, 9, 0]
// 		logfln("%#v\n", s)                // GoString: []int{1, 2, 3, 4, 5, 10, 7, 8, 9, 0}
// 		log(s.Index(5))                   // Index(5): 10
// 		log(s.Grow(10), s.Len(), s.Cap()) // <nil> , ok: 10, 20
// 		log(s.Interface().([]int))        // Interface : [1 2 3 4 5 10 7 8 9 0]

// 		sij, _ := s.Slice(0, 5)
// 		log(sij.Interface())             // Slice(i,j) : [1 2 3 4 5]
// 		log(s.Grow(5), s.Len(), s.Cap()) // cap - len > 5
// 		sij.SetIndex(0, 20)              // set sij[0] = 20; >> set s[i:j][0] = 20
// 		log(s)                           // [20, 2, 3, 4, 5, 10, 7, 8, 9, 0]

// 		sijk, _ := s.Slice3(0, 5, 7)
// 		log(sijk, sijk.Len(), sijk.Cap()) // [20, 2, 3, 4, 5], 5, 7

// 		s, _ = sijk.Append((int)(100), 200, 300)
// 		log(s)                                                   //[20, 2, 3, 4, 5, 100, 200, 300]
// 		log(s.AppendSlice(s.Value()))                            // [20, 2, 3, 4, 5, 100, 200, 300, 20, 2, 3, 4, 5, 100, 200, 300]
// 		log(s.AppendSlice(ValueOf([]int{0, 0, 0, 0, 0, 10086}))) // [20, 2, 3, 4, 5, 100, 200, 300, 0, 0, 0, 0, 0, 10086]
// 		s, _ = s.Slice(0, 5)
// 		s.Clear()
// 		log(s, s.Len(), s.Cap()) // [0, 0, 0, 0, 0], 5, 14
// 	}

// 	s, _ := helper.SliceFrom([]int{1, 2, 3, 4, 5, 6})
// 	log(s.ElemString(), s) // int [1, 2, 3, 4, 5, 6]

// 	log(helper.SliceFromValue(ValueOf(0))) // err: s is not a slice: kind=int

// 	log(helper.SliceFromValue(ValueOf([]int{9, 8, 7, 6}))) // [9, 8, 7, 6]
// }
