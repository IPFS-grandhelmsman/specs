---
menuTitle: Interpreter - 解释器
statusIcon: ⚠️
title: VM Interpreter - Message Invocation (Outside VM) - VM解释器-消息调用(VM外部)
entries:
- vm_outside
- vm_inside
# suppressMenu: true
---

(You can see the _old_ VM interpreter [here](docs/systems/filecoin_vm/vm_interpreter_old) )

(你可以查阅_老的_ VM解释器[这里](docs/systems/filecoin_vm/vm_interpreter_old) )

# `vm/interpreter` interface

{{< readfile file="vm_interpreter.id" code="true" lang="go" >}}

# `vm/interpreter` implementation

{{< readfile file="vm_interpreter.go" code="true" lang="go" >}}

# `vm/interpreter/registry`

{{< readfile file="vm_registry.go" code="true" lang="go" >}}
