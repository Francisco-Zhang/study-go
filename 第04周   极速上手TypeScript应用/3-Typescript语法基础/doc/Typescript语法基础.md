## 1、 基本数据类型

数据类型、逻辑控制、枚举类型、数组类型、对象类型、函数、函数式编程，这些都是 javascript 和 typescript 共有的部分。

Typescript 类型相关：接口、类、泛型

可以在下面 www.typescriptlang.org 这个网站，点击 try ts now -> in your browser ,打开Playgroud,进行ts代码编写

字符串既可以单引号也可以双引号，没有强制规定，但是在一个项目里要统一。

根据微信小程序的习惯，字符串 采用单引号，语句之间没有分号

在javascript 中定义了const 类型变量，可能过了很多行之后自己忘记了，又重新赋值，结果运行后报错，而不能在开发过程中提示。

var 最好不要用，只用 let 和 const

```typescript
let anExampleVariable = 'Hello World'
anExampleVariable = 123 //在 ts 中报错，但是编程成 js 后可以正常运行。
console.log(anExampleVariable)

let anExampleVariable = 'Hello World'
anExample = 123 //在 ts 中报错，会提前发现，但是js中没有提示，运行后才能发现错误。
console.log(anExampleVariable)
```

自己指定类型：

```typescript
let anExampleVariable:number = 1000  //就只能赋值数字，不能是字符串，反之也可以指定 :string
```

```typescript
#实际写代码的时候，能让编译器推断就不用手写
let anExampleVariable :string= 'Hello World'
let anExampleNum:number = 123
let anExampleBool: boolean =true
console.log(anExampleVariable,anExampleNum,anExampleBool)
```

literal type 类似于枚举，规定只能是固定的几种值

```typescript
let answer :'yes'|'no'|'maybe'='maybe'
let httpStatus :200|404|500|'200'|'404' = 200
```



## 2、 基本数据类型

typescript,虽然类型错误也能转成javascript。

```typescript
function t(s:200|'200'){
  //let status:string =s   此种写法错误，因为 s 有两种类型，union of types
  let status:string|number =s
}

let a: any = 'a'  //一旦加了any，就放弃了ts所有的类型检查，a 可以赋值任何类型值。
let b: undefined | undefined  //undefined 类型的变量值只能是 undefined
//undefined用法，一个变量的值，如下，除了几个固定的值，还可能是undefined
let answer :'yes'|'no'|'maybe'| undefined = undefined
```

## 3、 逻辑控制

```typescript
function processHttpstatus(s:200|404|'200'|'404'){
      //一律使用 === 以及 !==
      if (s===200){
        console.log("ok")
      }else if(s===404){
        console.log("not found")
      }else if(s==='200'){
        console.log("ok")
      }else if(s==='404'){
        console.log("not found")
      }else{
        //s 只有可能是这四个值，不需要写过多的 else,这个else 不需要写
      }
  }
processHttpstatus(200)
```

类型转换，这样就只用判断一种类型, 字符串转number 使用 parseInt()函数。parseInt(1.1) 返回结果1，parseFloat(1.1) 返回 1.1

```typescript
function processHttpstatus(s:200|404|'200'|'404'){
  let status=''
  if (typeof s === 'number'){
    status=s.toString()
  }else{
    status=s  //此处已经排除了number，可以直接赋值，但是没办法在status定义处赋值，因为类型不符
  }
  if(status==='200'){
    console.log("ok")
  }else if(status==='404'){
    console.log("not found")
  }
}
```

进一步简化

```typescript
function processHttpstatus(s:200|404|'200'|'404'){
  //不会变的值，尽可能使用const
  const status= typeof s==='string'? parseInt(s):s
	if(status===200){
    console.log("ok")
  }else if(status===404){
    console.log("not found")
  }
}
processHttpstatus('404')
```

改成switch写法

```typescript
function processHttpstatus(s:200|404|'200'|'404'){
     let status= typeof s==='string'? parseInt(s):s
     switch(status){
         case 200:
         console.log("ok")
         break; //一定要写break,否则会向下执行，会打印出多行数据。
         default:
          console.log("not found")
         break;
     }
  }
```

循环

```typescript
let sum=0
for(let i=0;i<100;i++){
  sum+=i
}
console.log(sum)


let sum=0
let i=1
while (i<=100){
  sum+=i
  i++
}
console.log(sum)
```

try catch语法：

```typescript
let sum=0
for(let i=0;i<100;i++){
  try{
    sum+=i
    if(i%17==0){
      throw `bad number ${i}`  //用反引号可以直接打印变量，不需要拼接。
    }
  }catch(err){
    console.error(err)
  }
}
console.log(sum)
```

## 4、 枚举类型
