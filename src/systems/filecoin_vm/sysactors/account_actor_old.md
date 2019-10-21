
# Account Actor - 帐号角色

- **Code Cid**: `<codec:raw><mhType:identity><"account">`

The Account actor is the actor used for normal keypair backed accounts on the filecoin network.

帐户角色是用于filecoin网络上备份普通密钥对帐户的角色。

```sh
type AccountActorState struct {
    address Address
}
```

## Methods - 方法

| Name | Method ID |
|--------|-------------|
| `AccountConstructor` | 1 |
| `GetAddress` | 2 |

```
type AccountConstructor struct {
}
```

## `GetAddress - 获取地址`

**Parameters - 参数**

```sh
type GetAddress struct {
} representation tuple
```

**Algorithm - 算法**

```go
func GetAddress() Address {
  return self.address
}
```
