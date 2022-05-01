## 1、 接口

### 为什么使用接口

在js中使用{}直接生成对象，那么相同类型的对象，怎么保证属性不拼写错误呢，就需要通过接口定义一个类型

```typescript
const emp1 = {
  name:"john",
  salary:8000
}

//定义第二个相同类型的对象，怎么保证类型相同，属性名以及类型不定义错误呢？
const emp2 = {
  
}
```

### 使用接口定义类型

```typescript
interface Employee {
  name:string;
  salary:number;
}

const emp1:Employee = {
  name:"john",
  salary:8000
}

const emp2:Employee = {
  name:"zhangsan",
  salary:8000
}
```

### 添加可选属性

```typescript
interface Employee {
  name:string		 // ；可以不写，只要统一就行
  salary:number
  bonus?:number  // ？表示这个属性可选，emp1 没有这个属性，也是允许的
}

const emp1:Employee = {
  name:"john",
  salary:8000
}

const emp2:Employee = {
  name:"zhangsan",
  salary:8000,
  bonus:200,
}

function updateBonus (e:Employee, p:number){
  if(e.bonus){  //可选字段，需要判断
    e.bonus = e.salary * p
  }
}
```

### 字段类型是函数

```typescript
interface Employee {
  name:string;
  salary:number;
  bonus?:number;
  updateBonus (p:number):void
}

const emp1:Employee = {
  name:"john",
  salary:8000,
  //这样的话，每一个对象都需要手写一个updateBonus方法，会很麻烦，所以这种场景一般使用类来实现。
  updateBonus(p:number){
    if(this.bonus){
      this.bonus = this.salary * p
    }
  }
}

//上面这种使用是错误的，一般接口里定义方法，是使用对象做参数时，用接口来描述对象有哪些字段
//比如微信里用到的一些场景，用来定义参数对象需要有哪些方法
interface GetUserProfileOption {
  /** 声明获取用户个人信息后的用途，不超过30个字符 */
  desc: string
  /** 接口调用结束的回调函数（调用成功、失败都会执行） */
  complete?: GetUserProfileCompleteCallback
  /** 接口调用失败的回调函数 */
  fail?: GetUserProfileFailCallback
  
  lang?: 'en' | 'zh_CN' | 'zh_TW'
  /** 接口调用成功的回调函数 */
  success?: GetUserProfileSuccessCallback
}
 wx.getUserProfile({
      desc: '需要你的个人信息来为你提供服务', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
      success: resolve,   //等同于 res => resolve(res),
      fail:reject,
 })
```

### 只读字段

```typescript
//通常用的比较少，一般定义组件某个属性只读可能会用到
interface Employee {
  readonly name:string;
  salary:number;
}
const emp1:Employee = {
  name:"john",
  salary:8000,
}
emp1.name = 'aa'  //只读字段，不允许这样写
```



## 2、 接口的高级技巧

### ？可选字段串联

```typescript
interface Employee {
  name?:{
    first?:string
    last:string
  };
  salary:number;
  bonus?:number;
}


function hasBadname(e:Employee){
  if(e.name && e.name.first){  //如果多个可选字段，这样写很麻烦
      return e.name.first.startsWith('aaa')
  }else{
    return true
  } 
}
// 简写方式：可选字段串联
function hasBadname(e:Employee){
  return e.name?.first?.startsWith('aaa') //如果不存在则返回 undefined，不再继续执行
}

console.log(hasBadname({
  name:{
    first:'john',
    last:'smith'
  },
  salary:200,
}))

console.log(hasBadname({
  salary:200,
}))
```



### ！非空断言

```typescript
// 使用！告诉编译器一定有值，不要编译报错，没值后果自负
function hasBadname(e:Employee){
  return e.name!.first!.startsWith('aaa')
}
```

### 接口扩展

类似于继承

```typescript
interface Employee extends HasName
 {
  salary:number;
  bonus?:number;
}

interface HasName {
 name?:{
    first?:string
    last:string
  }
}

function hasBadname(e:Employee){
  return e.name?.first?.startsWith('aaa')
}
```



### 类型的并

```typescript
interface WxButton {
  visivle:boolean
  enabled:boolean
  onClick():void
}
interface WxImage {
  visivle:boolean
  src:string
  width:number
  height:number
}
function hideElement(e:WxButton|WxImage){
		e.visivle = false //因为不管是哪种类型，都有visible 这个字段，所以可以这样写，其他字段不可以
}

```

### 类型的断言

```typescript
function processElement(e:WxButton|WxImage){
  //直接使用e.onClick会报错，所以需要先转成any,不让编译器判断类型
  //通过代码判断是否有 onClick 字段来判断 e 的类型
  if((e as any).onClick){
    const btn = e as WxButton //类型的断言
    btn.onClick()
  }else{
     const img = e as WxImage
     console.log(img.src)
  }
}
```

### 类型的判断

```typescript
//另外一种比较简单的写法，但是比较难理解。
function processElement(e:WxButton|WxImage){
  if(isButton(e)){
    e.onClick()
  }else{
     console.log(e.src)
  }
}

// e  is WxButton,这种写法就是类型的判断 
function isButton (e:WxButton|WxImage): e is WxButton {
  return (e as WxButton).onClick!== undefined
}
```

## 3、类

### 类的定义

在ts中使用类的场景并不多

```typescript
class Employee {
  name:string = '';  //类中必须初始化
}

//或者有构造函数
class Employee 
 {
  name:string
  salary:number
  constructor(name:string,salary:number){
    this.name = name
    this.salary = salary
  }
}
const emp1 = new Employee('john',8000)
console.log(emp1)
```

### 设置私有字段

```typescript
class Employee 
 {
  name:string
  salary:number
  private bonus:number = 0
  constructor(name:string,salary:number){
    this.name = name
    this.salary = salary
  }
}
const emp1 = new Employee('john',8000)  // 无法通过 emp1.bonus进行访问
```

### 简写

```typescript
class Employee 
 {
  private bonus?:number  //可选参数可以不初始化
  //添加public、private相当于定义了字段，可以少些一些代码,函数里边也不需要写初始化语句
  constructor(public name:string,private salary:number){}
}
```

### 添加方法

```typescript
class Employee 
 {
  private bonus?:number
  constructor(public name:string,private salary:number){}
  updateBonus(){
    if(this.bonus){
      this.bonus = 2000
    }
  }
}

const emp1 = new Employee('john',8000)
console.log(emp1)
```

### 特殊方法 getter 和 setter

```typescript
class Employee {
  private allocatedBonus?:number
  constructor(public name:string,public salary:number){}
  //并不是只有类中可以这样写，普通的json对象中也可以这样写，只是这样写没必要，所以普通的json对象中不使用。
  //类中这样写，是因为写在类中，getter 和 setter 可以复用。
  set bonus(v:number){
      this.allocatedBonus = v
  }
  get bonus(){
      return this.allocatedBonus || 0
  }
}
const emp1 = new Employee('john',8000)
emp1.bonus =2000
console.log(emp1)
```

### 继承

```typescript
class Manager extends Employee {
  private reporters: Employee[] = []
  addReporters(e:Employee){
    this.reporters.push(e)
  }
}

const manager1 = new Manager('mary',18000)

//子类也可以定义构造函数，但是参数中不能有public,否则相当于在子类中重新添加了 name,salary字段
class Manager extends Employee {
  private reporters: Employee[]

  constructor(name:string,salary:number){
    super(name,salary)
    this.reporters = []
  }
  addReporters(e:Employee){
    this.reporters.push(e)
  }
}

const emp1 = new Employee('john',8000)
emp1.bonus =2000

const manager1 = new Manager('mary',18000)
manager1.addReporters(emp1)
console.log(manager1)
```

## 4、 用类来实现接口

### 隐式实现

```typescript
interface Employee{
  name:string
  salary:number
  bonus?:number
}

class EmplImpl {
  private allocatedBonus?:number
  constructor(public name:string,public salary:number){}
  set bonus(v:number){
      this.allocatedBonus = v
  }
  get bonus(){
      return this.allocatedBonus || 0
  }
}

class Manager extends EmplImpl {
  private reporters: EmplImpl[]

  constructor(name:string,salary:number){
    super(name,salary)
    this.reporters = []
  }
  addReporters(e:EmplImpl){
    this.reporters.push(e)
  }
}

const emp1:Employee = new EmplImpl('john',8000)
```

### 使用 implements 显示实现

```typescript
interface Employee{
  name:string
  salary:number
  bonus?:number
}

class EmplImpl implements Employee{
  private allocatedBonus?:number
  constructor(public name:string,public salary:number){}
  set bonus(v:number){
      this.allocatedBonus = v
  }
  get bonus(){
      return this.allocatedBonus || 0
  }
}
```

### 使用implements的场景

implements 可以使用，也可以不使用，那什么场景下使用呢？

假设项目非常大，写 service.ts 的人和写 login.ts 的人不是同一个人。

1、定义者 = 实现者 != 使用者  这时候用显示实现，**这种方法不推荐**

```typescript
//service.ts
//定义接口
interface Service {
  login(): void
  getTrips(): string
  getLic(): string
  startTrip(): void
  updateLic(lic: string): void
}
//实现接口
class RPCService implements Service {
  login(): void {
    throw new Error("Method not implemented.")
  }
  getTrips(): string {
    throw new Error("Method not implemented.")
  }
  getLic(): string {
    throw new Error("Method not implemented.")
  }
  startTrip(): void {
    throw new Error("Method not implemented.")
  }
  updateLic(lic: string): void {
    throw new Error("Method not implemented.")
  }

}

// login page
// file: login.ts
const page = {
  service:new RPCService() as Service,//此处假设就是要使用Service接口
  onLoginButtonClicked(){
    //使用接口
    this.service.login()
  }
}
```

2、定义者 = 使用者 （推荐）,这时候使用隐式实现。

在登陆页面，使用者只想使用 login() 这一个方法，没必要拿到其他接口。

作为一个使用者，我知道使用哪些方法，我自己定义接口。这样每一个接口都非常小，实现者只要实现了我要使用的方法就可以。

这种方法符合代码解耦

```typescript
//实现接口
class RPCService {
  login(): void {
    throw new Error("Method not implemented.")
  }
  getTrips(): string {
    throw new Error("Method not implemented.")
  }
  getLic(): string {
    throw new Error("Method not implemented.")
  }
  startTrip(): void {
    throw new Error("Method not implemented.")
  }
  updateLic(lic: string): void {
    throw new Error("Method not implemented.")
  }
}


// login page
// file: login.ts
interface LoginService{
  login(): void
}

const page = {
  service:new RPCService() as LoginService,//此处假设就是要使用Service接口
  onLoginButtonClicked(){
    //使用接口
    this.service.login()
  }
}
```

### 5、 泛型
