#反射 The Laws of Reflection

##介绍
在这篇文章中我们将讲解golang的reflect是如何工作的，每一门语言的反射模型都不同，并且语多语言也不支持
##Types and interface
由于reflect是在类型(type)系统上建立的，所以我们先从类型开始复习.
Go是一个静态语言，每一个变量都有一个静态类型, 如 int, float32, *MyType, []byte等，如果我们声名如下的代码
```go
type MyInt int

var i int
var j MyInt 
```

则i的类型为int, j的类型为MyInt. 虽然变量i和j的底层类型都是int, 但变量i和j的类型是有明显的区别的, 它们彼此不能直接进行赋值，必须通过转换.

非常重要的一个类型种类为interface类型，它表示固定的方法集。一个interface变量可以存放作何具体的值, 只要它实现了了interface包含的方法. 众所周知的一个例子是io package的 
io.Reader和io.Writer

```go

// Reader is the interface that wraps the basic Read method.
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

任何类型，只要它实现了上面的Read方法签名(或者Write方法), 则实现了io.Reader或者io.Writer, 举例说明如下
```go
var r io.Reader
r = os.Stdin
r = bufio.NewReader(r)
r = new(bytes.Buffer)
// and so on
```
无论变量r的具体的值是什么，它的类型为io.Reader. Go语言是静态类型，r的静态类型为io.Reader

别一人上非常重要的概念是空interface

```go
interface{}
```
它表示一个空的方法集，而任何一个值都有零个方法或多个方法，所以空接口可以满足作何值的保存

有人说go的接口是动态类型，这是错误的，interface也是静态类型: 一个interface的变量在作何都是相同的静态作型(interface), 即使在运行时，保存的值改变类型，但是在怎么转变，值还是满足接口

###接口的描述
Russ Cox写过一篇 [博客](http://research.swtch.com/2009/12/go-data-structures-interfaces.html) 详细描述了go的接口。在这里我们只是简单的概括一下

接口类型的变量，存储了两部分，一个是分配给这个变量的具体值，一个是值的类型的描述器。更确切的说，值是底层具体的数据项，而类型描述了数据项的完赖类型, 比如
```go

var r io.Reader
tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
if err != nil {
    return nil, err
}
r = tty
```

r 包含一对(value, type), (tty, *os.File). 注意， *os.File实现了不仅实现了Read方法，还实现了其它方法，虽然这个接口变量只提供了Read 方法, 但是它的值包含了整个类型信息，所以我们可以按照下面的方法，转换为Writer

```go
var w io.Writer
w = r.(io.Writer)
``` 
这里使用了go的类型断言，它表示，变量r里面的项，也实现了io.Writer. 所以我们可以把它分配给w, 在分配完成后，w将包含（tty, *os.File). 这跟r是同一对数据. interface的静态
类型取决于可调用的方法，而不在乎具体值实现了多少方法.

```go
var empty interface{}
empty = w
```
empty将包含(tty, *os.File), 在这里我们没有使用类型断言，是因为w肯定满足于empty interface. Reader到writer，是因为writer的方法不是Reader的子集, 所以需要断言来测试

需要注意的一个细节是interface的一对值的形式是（value, concrete type) 不而不(value, interface type). Interfaces不能保存interface的值

接下来讲reflect

##第一条Reflect规则
1. 从interface value反射到对像
根源上来说， reflection的原理就是检查interface中保存的一对值和类型, 所以在reflect包中，有两个类型我们需要记住， Type和Value两个类型. 通过这两个类型，我们可以访问一个
接口变量的内容. 调用reflect.ValueOf和reflect.TypeOf可以检索出一个interface的值和具体类型. 当然通过reflect.Value我们也可以很空易的获得reflect.Type， 但我们在这还是先将这两个概念独立开.

接下来我们从TypeOf开始

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    var x float64 = 3.4
    fmt.Println("type:", reflect.TypeOf(x))
}
```
输出
```go
type: float64
```

你可能很想知道，我们上面所以说的interface在哪， 由于这个函数只是将float64传递给变量x, 而不是interface值传递给reflect.TypeOf. 为什么？因为reflect.TypeOf的参数为interface{}空接口

````go
// TypeOf returns the reflection Type of the value in the interface{}.
func TypeOf(i interface{}) Type
````
当我们调用reflect.TypeOf(x), x首先存储在一个空接口上，然后在作为参数传递给TypeOf; Reflect.TypeOf解压这个空接口, 接收类型信息

同理，reflect.ValueOf是一个的，只是它接收到的是一个Value

```go
var x float64 = 3.4
fmt.Println("value:", reflect.ValueOf(x).String())
```
打印
```go
value: <float64 Value>
```

在这里，我们使用了String()方法，是因为默认的fmt包，会打印出reflect.Value的具体值3.4. 而 String方法不是

reflect.Type和reflect.Value都有大量的方法让我们检查和操作它们。一个重要的例子是 Value有一个Type方法，它反回reflect.Value的Type类型. 虽一个是Type和Value
都有一个Kind方法，这两个方法返回一个常量，表示interface存储的项是什么,比如Uint, Float64, Slice等等。同样value还有类似于int, Float方法，让我们检索interface中具体的值

```go
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())
fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
fmt.Println("value:", v.Float())
```
打印
```go
type: float64
kind is float64: true
value: 3.4
```

这里同样有SetInt和SetFloat等方法，但是在使用它们之前，我们需要理解什么是settability, 这将在关于reflect第三条法则中讨论 

reflect包中有一个很重要的特性， 第一个为,为了保持API的简单，Value类型的"getter"和"setter"方法都是在最大类型上操作，比如int64，保存所有有符号的整数, 
举例来说， Int方法返回的一是个int64，而SetInt接受的是一个int64位的参数 。它可能需要通过以下的方法转换为实现了类型

```go
var x uint8 = 'x'
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())                            // uint8.
fmt.Println("kind is uint8: ", v.Kind() == reflect.Uint8) // true.
x = uint8(v.Uint()) 
```

第二个特性是Kind描述的是reflection对像的底层类型，而不是静态类型。 假如一个reflection对像包含了一个用户自定义的静态类型

```go
type MyInt int
var x MyInt = 7
v := reflect.ValueOf(x)
```
v 的Kind方法依然为reflect.Int, 即使x的静态类型为MyInt, 而不是int. 换句话说，从Kind上不能够区分int和MyInt. 但是Type可以

##第二条法则 
2. 将一个reflection对像转换为interface值
像照镜子一样的物理反射一样，在Go语言中的反射它是生成一个自已的反转
通过reflect.Value的Interface方法，我们可以获得一个Interface值。实际上这个方法将一个type和value打包回interface

```go
// Interface returns v's value as an interface{}.
func (v Value) Interface() interface{}
```
结果为
```go
y := v.Interface().(float64) // y will have type float64.
fmt.Println(y)
```

打印的float64的值表示是reflection对像v
我们还可以更进一步， fmt.Println, fmt.Printf等方法的参数可以是一个空interface值，我们可以像上面那样，让fmt包内部解压出y的值.

```go
fmt.Println(v.Interface())
```
(其实也可以是fmt.Println(v)),  由于我们的value是一个float64, 所以我们还可以指定float格式
```go
fmt.Printf("value is %7.1e\n", v.Interface())
```
输出
```go
3.4e+00
```
再次说明，这里不需要使用断言，将v.Interface转换为float64. 一个空interface的值拥有一个具体的值的具体类型信息，而printf可以重新取得这些信息

简单的来说，Interface方法是ValueOf方法的反过程, 不同的是它的结果总是返回一个静态类型interface{}

> 重申一下： Reflection可以从interface至reflection对像，也可以从reflection对像到interface

##第三条法则
3. 为了修改reflection 对像，它的值必须是可settable.
第三条是最难理解的，但如果从第一条法则开始，其它它还是很好理解
以下的代码不是有用的代码，但是值得学习，帮助我们理解
```go
var x float64 = 3.4
v := reflect.ValueOf(x)
x.SetFloat(7.1) //Error: will panic
```
如果你运行这段代码，它将抛出一个异常
```go
panic: reflect.Value.SetFloat using unaddressable value //reflect.Value.SetFloat使用了一个不可寻址的值
```
这个问题不是说值7.1是不可寻址的，而是v是不可设置的(settable). Settability(可设置)是reflection Value的一个属性, 并不是所有的reflection Values都拥有这个属性.

我们可以通过CanSet方法，测试value是否可测试
```go
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println("settability of v:", v.CanSet())
```
打印
```go
settability of v: false
```
在一个不可设置的(false)的Value上调用set方法将报错. 但什么是可设置的呢？

Settability有点类似于是否可寻址(addressability), 但严格来说，它表示的是reflection对像是否可以修改创建这个reflection对像的实际值. 可设置取决于reflection对像所持有的原始值.

```go
var x float64 = 3.4
v := reflect.ValueOf(x)
```
在这里，x作为参数传递时，首先复制x. 所以reflect.ValueOf中的interface值是x的复制，而不是x本身. 所以，如果下面的语句
```go
v.SetFloat(7.1)
```
是被允许的，它也不会更新x, 而是更新x的复制, 所以x本身是没有影响的，这很混乱，也是无用的，所以它是非法的.

 如果你不是很清楚，我们将通过下面的方式，进一步向你解释
 
 ```go
f(x)
```
我们可能不希望f函数修改x, 所以我们传递的是x的复制，而不是x它本身, 如果我们想要f修改，则可以传递x的地址（一个指向x的指针)

```go
f(&x)
```
这非常简单明了，reflection也是使用相同的方式, 如果我们想要通过reflection修改x, 我们必须传递一个x的指针。
代码如下所示
```go
var x float64 = 3.4
p := reflect.ValueOf(&x) // Note: take the address of x.
fmt.Println("type of p:", p.Type())
fmt.Println("settability of p:", p.CanSet())
```
打印
```go
type of p: *float64
settability of p: false
```
reflection对像p还是不可设置,但其时我们也不是想设置p，而针指所指向的值. 为了获取p所指向的值, 我们调用value的ELem方法。

```go
v := p.Elem()
fmt.Println("settability of v:", v.CanSet())
```
现在v是一个可设置的reflection以像
```go
settability of v: true
```
由于它代表的是变量x， 我们可以通过v.SetFloat来改变x的值
```go
v.SetFloat(7.1)
fmt.Println(v.Interface())
fmt.Println(x)
```
结果如我们期望的那样
```go
7.1
7.1
```
###struct
在我们上一个例子中，v本身不是一个指针，它仅是从一个指针推导而来。对于要修改值的情况，常常用于修改一个struct的字段，只要有struct的地址，我们就可以修改它。
下面是一个简单的例子，用来分析一个struct的值. 我们创建了一个relection对像，它是一个struct的地址.因为我们在这之后需要修改它。

```go
type T struct {
    A int
    B string
}
t := T{23, "skidoo"}
s := reflect.ValueOf(&t).Elem()
typeOfT := s.Type()
for i := 0; i < s.NumField(); i++ {
    f := s.Field(i)
    fmt.Printf("%d: %s %s = %v\n", i,
        typeOfT.Field(i).Name, f.Type(), f.Interface())
}
```
输出
```go
0: A int = 23
1: B string = skidoo
```
这里我们需要强调一下，只有大写开头的字段(Export)才可设置

```go
s.Field(0).SetInt(77)
s.Field(1).SetString("Sunset Strip")
fmt.Println("t is now", t)
```
打印
```go
t is now {77 Sunset Strip}
```
##总结

reflection的三条法则
1. 从interface值到reflection 对像(ValueOf, Typeof)
2. reflection对像到interface 
3. 更改reflection对像，它的值必须是可更改的

