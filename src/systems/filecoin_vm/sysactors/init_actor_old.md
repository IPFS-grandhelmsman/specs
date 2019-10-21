# Init Actor - 初始化角色

- **Code Cid**: `<codec:raw><mhType:identity><"init">`

The init actor is responsible for creating new actors on the filecoin network. This is a built-in actor and cannot be replicated. In the future, this actor will be responsible for loading new code into the system (for user programmable actors). ID allocation for user instantiated actors starts at 100. This means that `NextID` will initially be set to 100.

初始化角色负责在filecoin网络上创建新的角色，这是一个内置的角色，不能被复制。将来，这个角色将负责向系统中加载新代码(对于用户可编程角色)。用户实例化的参与者的ID分配从100开始。这意味着`NextID`将会被初始化设置为100。

```sh
type InitActorState struct {
    addressMap {Address:ID}<Hamt>
    nextId UInt
}
```

## Methods - 方法

| Name | Method ID |
|--------|-------------|
| `Constructor` | 1 |
| `Exec` | 2 |
| `GetIdForAddress` | 3 |

## `Constructor - 构造器`

**Parameters - 参数**

```sh
type InitConstructor struct {
}

```

**Algorithm - 算法**

## `Exec - 执行`

This method is the core of the `Init Actor`. It handles instantiating new actors and assigning them their IDs.

这个方法是`初始化角色`的核心。它处理实例化新的参与者并为其分配id。

**Parameters - 参数**

```sh
type Exec struct {
    ## Reference to the location at which the code of the actor to create is stored.
    code &Code
    ## Parameters passed to the constructor of the actor.
    params ActorMethod
} representation tuple
```

**Algorithm - 算法**

```go
func Exec(code Cid, params ActorMethod) Address {
  // Get the actor ID for this actor.
  actorID = self.NextID
  self.NextID++

  // Make sure that only the actors defined in the spec can be launched.
  if !IsBuiltinActor(code) {
    Fatal("cannot launch actor instance that is not a builtin actor")
  }

  // Ensure that singeltons can be only launched once.
  // TODO: do we want to enforce this? If so how should actors be marked as such?
  if IsSingletonActor(code) {
    Fatal("cannot launch another actor of this type")
  }

  // This generates a unique address for this actor that is stable across message
  // reordering
  // TODO: where do `creator` and `nonce` come from?
  addr := VM.ComputeActorAddress(creator, nonce)

  // Set up the actor itself
  actor := Actor{
    Code:    code,
    Balance: msg.Value,
    Head:    nil,
    Nonce:   0,
  }

  // The call to the actors constructor will set up the initial state
  // from the given parameters, setting `actor.Head` to a new value when successfull.
  // TODO: can constructors fail?
  actor.Constructor(params)

  VM.GlobalState.Set(actorID, actor)

  // Store the mapping of address to actor ID.
  self.AddressMap[addr] = actorID

  return addr
}

func IsSingletonActor(code Cid) bool {
  return code == StorageMarketActor || code == InitActor
}
```

```go
// TODO: find a better home for this logic
func (VM VM) ComputeActorAddress(creator Address, nonce Integer) Address {
  return NewActorAddress(bytes.Concat(creator.Bytes(), nonce.BigEndianBytes()))
}
```

## `GetIdForAddress - 为地址获取ID`

This method allows for fetching the corresponding ID of a given Address

此方法允许获取给定地址的对应ID

**Parameters - 参数**

```sh
type GetIdForAddress struct {
    addr Address
} representation tuple
```

**Algorithm -参数**

```go
func GetIdForAddress(addr Address) UInt {
  id := self.AddressMap[addr]
  if id == nil {
    Fault("unknown address")
  }
  return id
}
```
